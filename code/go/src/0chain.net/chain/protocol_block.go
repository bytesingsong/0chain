package chain

import (
	"context"
	"math"
	"time"

	"0chain.net/config"
	"0chain.net/datastore"
	"0chain.net/node"
	"0chain.net/state"
	"0chain.net/util"

	"0chain.net/block"
	"0chain.net/common"
	. "0chain.net/logging"
	"go.uber.org/zap"
)

/*VerifyTicket - verify the ticket */
func (c *Chain) VerifyTicket(ctx context.Context, blockHash string, bvt *block.VerificationTicket) error {
	sender := c.Miners.GetNode(bvt.VerifierID)
	if sender == nil {
		return common.InvalidRequest("Verifier unknown or not authorized at this time")
	}

	if ok, _ := sender.Verify(bvt.Signature, blockHash); !ok {
		return common.InvalidRequest("Couldn't verify the signature")
	}
	return nil
}

/*VerifyNotarization - verify that the notarization is correct */
func (c *Chain) VerifyNotarization(ctx context.Context, blockHash string, bvt []*block.VerificationTicket) error {
	if bvt == nil {
		return common.NewError("no_verification_tickets", "No verification tickets for this block")
	}
	ticketsMap := make(map[string]bool, len(bvt))
	for _, vt := range bvt {
		if _, ok := ticketsMap[vt.VerifierID]; ok {
			return common.NewError("duplicate_ticket_signature", "Found duplicate signatures in the notarization of the block")
		}
		ticketsMap[vt.VerifierID] = true
	}
	if !c.reachedNotarization(bvt) {
		return common.NewError("block_not_notarized", "Verification tickets not sufficient to reach notarization")
	}
	for _, vt := range bvt {
		if err := c.VerifyTicket(ctx, blockHash, vt); err != nil {
			return err
		}
	}
	return nil
}

/*IsBlockNotarized - check if the block is notarized */
func (c *Chain) IsBlockNotarized(ctx context.Context, b *block.Block) bool {
	if b.IsBlockNotarized() {
		return true
	}
	notarized := c.reachedNotarization(b.VerificationTickets)
	if notarized {
		b.SetBlockNotarized()
	}
	return notarized
}

func (c *Chain) reachedNotarization(bvt []*block.VerificationTicket) bool {
	if c.ThresholdByCount > 0 {
		numSignatures := len(bvt)
		if numSignatures < c.GetNotarizationThresholdCount() {
			return false
		}
	}
	if c.ThresholdByStake > 0 {
		verifiersStake := 0
		for _, ticket := range bvt {
			verifiersStake += c.getMiningStake(ticket.VerifierID)
		}
		if verifiersStake < c.ThresholdByStake {
			return false
		}
	}
	return true
}

/*UpdateNodeState - based on the incoming valid blocks, update the nodes that notarized the block to be active
 Useful to increase the speed of node status discovery which increases the reliablity of the network
Simple 3 miner scenario :

1) a discovered b & c.
2) b discovered a.
3) b and c are yet to discover each other
4) a generated a block and sent it to b & c, got it notarized and next round started
5) c is the generator who generated the block. He will only send it to a as b is not discovered to be active.
    But if the prior block has b's signature (may or may not, but if it did), c can discover b is active before generating the block and so will send it to b
*/
func (c *Chain) UpdateNodeState(b *block.Block) {
	for _, vt := range b.VerificationTickets {
		signer := c.Miners.GetNode(vt.VerifierID)
		if signer == nil {
			Logger.Error("this should not happen!")
			continue
		}
		if signer.GetStatus() != node.NodeStatusActive {
			signer.SetStatus(node.NodeStatusActive)
		}
	}
}

/*AddVerificationTicket - add a verified ticket to the list of verification tickets of the block */
func (c *Chain) AddVerificationTicket(ctx context.Context, b *block.Block, bvt *block.VerificationTicket) bool {
	added := b.AddVerificationTicket(bvt)
	if added {
		c.IsBlockNotarized(ctx, b)
	}
	return added
}

/*MergeVerificationTickets - merge a set of verification tickets (already validated) for a given block */
func (c *Chain) MergeVerificationTickets(ctx context.Context, b *block.Block, vts []*block.VerificationTicket) {
	vtlen := len(b.VerificationTickets)
	b.MergeVerificationTickets(vts)
	if len(b.VerificationTickets) != vtlen {
		c.IsBlockNotarized(ctx, b)
	}
}

func (c *Chain) finalizeBlock(ctx context.Context, fb *block.Block, bsh BlockStateHandler) {
	bNode := node.GetNode(fb.MinerID)
	ms := bNode.ProtocolStats.(*MinerStats)
	Logger.Info("finalize block", zap.Int64("round", fb.Round), zap.Int64("current_round", c.CurrentRound), zap.Int64("lf_round", c.LatestFinalizedBlock.Round), zap.String("hash", fb.Hash), zap.Int("round_rank", fb.RoundRank), zap.Int8("state", fb.GetBlockState()))
	ms.FinalizationCountByRank[fb.RoundRank]++
	if time.Since(ssFTs) < 20*time.Second {
		SteadyStateFinalizationTimer.UpdateSince(ssFTs)
	}
	StartToFinalizeTimer.UpdateSince(fb.ToTime())
	ssFTs = time.Now()
	c.UpdateChainInfo(fb)
	if fb.GetStateStatus() != block.StateSuccessful {
		err := c.ComputeState(ctx, fb)
		if err != nil {
			if config.DevConfiguration.State {
				Logger.Error("finalize block state not successful", zap.Int64("finalized_round", fb.Round), zap.String("hash", fb.Hash), zap.Int8("state", fb.GetBlockState()), zap.Error(err))
				if state.Debug() {
					Logger.DPanic("finalize block - state not successful")
				}
			}
		}
	}
	if fb.ClientState != nil {
		ts := time.Now()
		err := fb.ClientState.SaveChanges(c.stateDB, false)
		duration := time.Since(ts)
		StateSaveTimer.UpdateSince(ts)
		p95 := StateSaveTimer.Percentile(.95)
		changes := fb.ClientState.GetChangeCollector().GetChanges()
		if len(changes) > 0 {
			StateChangeSizeMetric.Update(int64(len(changes)))
		}
		if StateSaveTimer.Count() > 100 && 2*p95 < float64(duration) {
			Logger.Error("finalize block - save state slow", zap.Int64("round", fb.Round), zap.String("block", fb.Hash), zap.Int("block_size", len(fb.Txns)), zap.Int("changes", len(changes)), zap.String("client_state", util.ToHex(fb.ClientStateHash)), zap.Duration("duration", duration), zap.Duration("p95", time.Duration(math.Round(p95/1000000))*time.Millisecond))
		} else {
			Logger.Debug("finalize block - save state", zap.Int64("round", fb.Round), zap.String("block", fb.Hash), zap.Int("block_size", len(fb.Txns)), zap.Int("changes", len(changes)), zap.String("client_state", util.ToHex(fb.ClientStateHash)), zap.Duration("duration", duration))
		}
		if err != nil {
			Logger.Error("finalize block - save state", zap.Int64("round", fb.Round), zap.String("block", fb.Hash), zap.Int("block_size", len(fb.Txns)), zap.Int("changes", len(changes)), zap.String("client_state", util.ToHex(fb.ClientStateHash)), zap.Duration("duration", duration), zap.Error(err))
		}
		c.rebaseState(fb)
	}
	if config.Development() {
		ts := time.Now()
		for _, txn := range fb.Txns {
			StartToFinalizeTxnTimer.Update(ts.Sub(common.ToTime(txn.CreationDate)))
		}
	}
	go bsh.UpdateFinalizedBlock(ctx, fb)
	c.BlockChain.Value = fb.GetSummary()
	c.BlockChain = c.BlockChain.Next()

	// Deleting dead blocks from a couple of rounds before (helpful for visualizer and potential rollback scenrio)
	pfb := fb
	for idx := 0; idx < 10 && pfb != nil; idx, pfb = idx+1, pfb.PrevBlock {

	}
	if pfb == nil {
		return
	}
	frb := c.GetRoundBlocks(pfb.Round)
	var deadBlocks []*block.Block
	for _, b := range frb {
		if b.Hash != pfb.Hash {
			deadBlocks = append(deadBlocks, b)
		}
	}
	// Prune all the dead blocks
	c.DeleteBlocks(deadBlocks)
}

/*GetNotarizedBlock - get a notarized block for a round */
func (c *Chain) GetNotarizedBlock(blockHash string) *block.Block {
	nbrequestor := MinerNotarizedBlockRequestor
	cround := c.CurrentRound
	params := map[string]string{"block": blockHash}
	ctx := common.GetRootContext()
	var b *block.Block
	handler := func(ctx context.Context, entity datastore.Entity) (interface{}, error) {
		Logger.Info("get notarized block", zap.String("block", blockHash), zap.Int64("cround", cround), zap.Int64("current_round", c.CurrentRound))
		nb, ok := entity.(*block.Block)
		if !ok {
			return nil, datastore.ErrInvalidEntity
		}
		if err := c.VerifyNotarization(ctx, nb.Hash, nb.VerificationTickets); err != nil {
			Logger.Error("get notarized block - validate notarization", zap.String("block", blockHash), zap.Error(err))
			return nil, err
		}
		if err := nb.Validate(ctx); err != nil {
			Logger.Error("get notarized block - validate", zap.String("block", blockHash), zap.Any("block_obj", nb), zap.Error(err))
			return nil, err
		}
		r := c.GetRound(nb.Round)
		if r == nil {
			Logger.Error("get notarized block - no round (TODO)", zap.String("block", blockHash), zap.Int64("round", nb.Round), zap.Int64("cround", cround), zap.Int64("current_round", c.CurrentRound))
			b = c.AddBlock(nb)
		} else {
			c.SetRandomSeed(r, nb.RoundRandomSeed)
			b = c.AddRoundBlock(r, nb)
			b, _ = r.AddNotarizedBlock(b)
		}
		Logger.Info("get notarized block", zap.Int64("round", b.Round), zap.String("block", b.Hash))
		return b, nil
	}
	n2n := c.Miners
	n2n.RequestEntity(ctx, nbrequestor, params, handler)
	return b
}

/*GetPreviousBlock - get the previous block from the network */
func (c *Chain) GetPreviousBlock(ctx context.Context, b *block.Block) *block.Block {
	if b.PrevBlock != nil {
		return b.PrevBlock
	}
	pb, err := c.GetBlock(ctx, b.PrevHash)
	if err == nil {
		b.SetPreviousBlock(pb)
		return pb
	}
	blocks := make([]*block.Block, 0, 10)
	Logger.Info("fetch previous block", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.String("prev_block", b.PrevHash))
	cb := b
	for idx := 0; idx < 10; idx++ {
		Logger.Error("fetching previous block", zap.Int("idx", idx), zap.Int64("cround", cb.Round), zap.String("cblock", cb.Hash), zap.String("cprev_block", cb.PrevHash))
		nb := c.GetNotarizedBlock(cb.PrevHash)
		if nb == nil {
			Logger.Error("get previous block (unable to get prior blocks)", zap.Int64("current_round", c.CurrentRound), zap.Int("idx", idx), zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Int64("cround", cb.Round), zap.String("cblock", cb.Hash), zap.String("cprev_block", cb.PrevHash))
			return nil
		}
		cb = nb
		blocks = append(blocks, cb)
		pb, err = c.GetBlock(ctx, cb.PrevHash)
		if pb != nil {
			cb.SetPreviousBlock(pb)
			break
		}
	}
	if cb.PrevBlock == nil { // This happens after fetching as far as per the previous for loop and still not having the prior block
		Logger.Error("get previous block (missing continuity)", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Int64("oldest_fetched_round", cb.Round), zap.String("oldest_fetched_block", cb.Hash), zap.String("missing_prior_block", cb.PrevHash))
		return nil
	}
	for idx := len(blocks) - 1; idx >= 0; idx-- {
		cb := blocks[idx]
		if cb.PrevBlock == nil {
			pb, err := c.GetBlock(ctx, cb.PrevHash)
			if err != nil {
				Logger.Error("get previous block (missing continuity)", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Int64("cb_round", cb.Round), zap.String("cb_block", cb.Hash), zap.String("missing_prior_block", cb.PrevHash))
				return nil
			}
			cb.SetPreviousBlock(pb)
		}
		c.ComputeState(ctx, cb)
	}
	pb, err = c.GetBlock(ctx, b.PrevHash)
	if err == nil {
		b.SetPreviousBlock(pb)
	}
	return pb
}

//Note: this is expected to work only for small forks
func (c *Chain) commonAncestor(ctx context.Context, b1 *block.Block, b2 *block.Block) *block.Block {
	if b1 == nil || b2 == nil {
		return nil
	}
	if b1 == b2 || b1.Hash == b2.Hash {
		return b1
	}
	if b2.Round < b1.Round {
		b1, b2 = b2, b1
	}
	for b2.Round != b1.Round {
		b2 = c.GetPreviousBlock(ctx, b2)
		if b2 == nil {
			return nil
		}
	}
	for b1 != b2 {
		b1 = c.GetPreviousBlock(ctx, b1)
		if b1 == nil {
			return nil
		}
		b2 = c.GetPreviousBlock(ctx, b2)
		if b2 == nil {
			return nil
		}
	}
	return b1
}
