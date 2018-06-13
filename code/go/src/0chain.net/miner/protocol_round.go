package miner

import (
	"context"
	"fmt"
	"time"

	"0chain.net/block"
	"0chain.net/common"
	"0chain.net/datastore"
	. "0chain.net/logging"
	"0chain.net/memorystore"
	"0chain.net/node"
	"go.uber.org/zap"
)

const BLOCK_TIME = 300 * time.Millisecond
const FINALIZATION_TIME = 300 * time.Millisecond

/*GetBlockToExtend - Get the block to extend from the given round */
func (mc *Chain) GetBlockToExtend(r *Round) *block.Block {
	//TODO: We need to ensure the block exists but also that it has received the notarization
	if r.Block != nil {
		return r.Block
	}
	return nil
}

/*GenerateRoundBlock - given a round number generates a block*/
func (mc *Chain) GenerateRoundBlock(ctx context.Context, r *Round) (*block.Block, error) {
	pround := mc.GetRound(r.Number - 1)
	if pround == nil {
		Logger.Error("generate block (prior round not found)", zap.Any("round", r.Number-1))
		return nil, common.NewError("invalid_round,", "Round not available")
	}
	pb := mc.GetBlockToExtend(pround)
	if pb == nil {
		Logger.Error("generate block (prior block not found)", zap.Any("round", r.Number))
		return nil, common.NewError("block_gen_no_block_to_extend", "Do not have the block to extend this round")
	}
	b := datastore.GetEntityMetadata("block").Instance().(*block.Block)
	b.ChainID = mc.ID
	b.SetPreviousBlock(pb)
	err := mc.GenerateBlock(ctx, b)
	if err != nil {
		Logger.Error("generate block error", zap.Error(err))
		return nil, err
	}
	mc.AddBlock(b)
	mc.SendBlock(ctx, b)
	return b, nil
}

/*CollectBlocksForVerification - keep collecting the blocks till timeout and then start verifying */
func (mc *Chain) CollectBlocksForVerification(ctx context.Context, r *Round) {
	var blockTimeTimer = time.NewTimer(BLOCK_TIME)
	var sendVerification = false
	verifyAndSend := func(ctx context.Context, r *Round, b *block.Block) bool {
		pb := r.Block
		r.Block = b
		bvt, err := mc.VerifyRoundBlock(ctx, r, b)
		if err != nil {
			r.Block = pb
			if err == ErrRoundMismatch {
				Logger.Info("verify round block", zap.Any("round", r.Number), zap.Any("block", b.Hash), zap.Any("current_round", mc.CurrentRound))
			} else {
				Logger.Error("verify round block", zap.Any("round", r.Number), zap.Any("block", b.Hash), zap.Error(err))
			}
			return false
		}
		if b.MinerID != node.Self.GetKey() {
			mc.SendVerificationTicket(ctx, b, bvt)
		}
		mc.ProcessVerifiedTicket(ctx, r, b, &bvt.VerificationTicket)
		return true
	}
	var blocks = make([]*block.Block, 0, 10)
	for true {
		select {
		case <-ctx.Done():
			return
		case <-blockTimeTimer.C:
			sendVerification = true
			// Sort the accumulated blocks by the rank and process them
			blocks = r.GetBlocksByRank(blocks)
			// Keep verifying all the blocks collected so far in the best rank order till the first
			// successul verification
			for _, b := range blocks {
				if verifyAndSend(ctx, r, b) {
					break
				}
			}
		case b := <-r.GetBlocksToVerifyChannel():
			if sendVerification {
				// Is this better than the current best block
				if r.Block == nil || b.RoundRank < r.Block.RoundRank {
					verifyAndSend(ctx, r, b)
				}
			} else { // Accumulate all the blocks into this array till the BlockTime timeout
				blocks = append(blocks, b)
			}
		}
	}
}

/*VerifyRoundBlock - given a block is verified for a round*/
func (mc *Chain) VerifyRoundBlock(ctx context.Context, r *Round, b *block.Block) (*block.BlockVerificationTicket, error) {
	if mc.CurrentRound != r.Number {
		return nil, ErrRoundMismatch
	}
	if b.MinerID == node.Self.GetKey() {
		return mc.SignBlock(ctx, b)
	}
	prevBlock, err := mc.GetBlock(ctx, b.PrevHash)
	if err != nil {
		//TODO: create previous round AND request previous block from miner who sent current block for verification
		Logger.Error("verify round", zap.Any("round", r.Number), zap.Any("block", b.Hash), zap.Any("prev_block", b.PrevHash), zap.Error(err))
		return nil, common.NewError("prev_block_error", "Error getting the previous block")
	}

	if prevBlock == nil {
		//TODO: create previous round AND request previous block from miner who sent current block for verification
		return nil, common.NewError("invalid_block", fmt.Sprintf("Previous block doesn't exist: %v", b.PrevHash))
	}

	/* Note: We are verifying the notrization of the previous block we have with
	   the prev verification tickets of the current block. This is right as all the
	   necessary verification tickets & notarization message may not have arrived to us */
	if err := mc.VerifyNotarization(ctx, prevBlock, b.PrevBlockVerficationTickets); err != nil {
		return nil, err
	}

	bvt, err := mc.VerifyBlock(ctx, b)
	if err != nil {
		return nil, err
	}
	return bvt, nil
}

/*ComputeFinalizedBlock - compute the block that has been finalized. It should be the one in the prior round */
func (mc *Chain) ComputeFinalizedBlock(ctx context.Context, r *Round) *block.Block {
	// TODO: current behavior is we are returning r.Block
	//return r, r.Block
	tips := r.GetNotarizedBlocks()
	for true {
		ntips := make([]*block.Block, 0, 1)
		for _, b := range tips {
			if b.Hash == mc.LatestFinalizedBlock.Hash {
				break
			}
			found := false
			for _, nb := range ntips {
				if b.PrevHash == nb.Hash {
					found = true
					break
				}
			}
			if found {
				continue
			}
			ntips = append(ntips, b.PrevBlock)
		}
		tips = ntips
		if len(tips) == 1 {
			break
		}
	}
	if len(tips) != 1 {
		return nil
	}
	fb := tips[0]
	if fb.Round == r.Number {
		return nil
	}
	return fb
}

/*FinalizeRound - starting from the given round work backwards and identify the round that can be assumed to be finalized as only one chain has survived */
func (mc *Chain) FinalizeRound(ctx context.Context, r *Round) {
	/*TODO: This is incorrect because when we ask r to finalize, it's actually the r-1 round's block that gets finalized */
	if r.IsFinalized() {
		return
	}
	var finzalizeTimer = time.NewTimer(FINALIZATION_TIME)
	select {
	case <-finzalizeTimer.C:
		break
	}
	fb := mc.ComputeFinalizedBlock(ctx, r)
	if fb == nil {
		Logger.Info("finalization - no decisive block to finalize yet", zap.Any("round", r.Number))
		return
	}
	frchain := make([]*block.Block, 0, 1)
	for b := fb; b != nil && b.Hash != mc.LatestFinalizedBlock.Hash; b = b.PrevBlock {
		frchain = append(frchain, b)
	}
	for idx := range frchain {
		fb = frchain[len(frchain)-1-idx]
		fr := mc.GetRound(fb.Round)
		fr.Finalize()
		Logger.Info("finalizing round", zap.Any("round", r.Number), zap.Any("finalized_round", fb.Round), zap.Any("hash", fb.Hash))
		mc.UpdateFinalizedBlock(ctx, fb)
		mc.SendFinalizedBlock(ctx, fb)
	}
	mc.LatestFinalizedBlock = fb
}

/*UpdateFinalizedBlock - update the latest finalized block */
func (mc *Chain) UpdateFinalizedBlock(ctx context.Context, b *block.Block) {
	if b.Hash == mc.LatestFinalizedBlock.Hash {
		return
	}
	txnEntityMetadata := datastore.GetEntityMetadata("txn")
	ctx = memorystore.WithEntityConnection(ctx, txnEntityMetadata)
	defer memorystore.Close(ctx)
	mc.FinalizeBlock(ctx, b)
}
