package storagesc

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"0chain.net/core/datastore"
	"0chain.net/smartcontract/dbs/benchmark"
	"0chain.net/smartcontract/provider"

	"0chain.net/smartcontract/stakepool/spenum"
	"github.com/0chain/common/core/currency"

	"0chain.net/smartcontract/stakepool"

	"0chain.net/smartcontract/partitions"

	"0chain.net/smartcontract/dbs/event"

	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/core/encryption"
	sc "0chain.net/smartcontract/benchmark"
	"github.com/spf13/viper"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/core/common"
)

const (
	mockFinalizedAllocationIndex = 2
)

func AddMockAllocations(
	clients, publicKeys []string,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) {
	for i := 0; i < viper.GetInt(sc.NumAllocations); i++ {
		cIndex := getMockOwnerFromAllocationIndex(i, len(clients))
		addMockAllocation(
			i,
			clients,
			cIndex,
			publicKeys[cIndex],
			eventDb,
			balances,
		)
	}
}

func benchAllocationExpire(now common.Timestamp) common.Timestamp {
	return common.Timestamp(viper.GetDuration(sc.TimeUnit).Seconds()) + now
}

func addMockAllocation(
	i int,
	clients []string,
	cIndex int,
	publicKey string,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) {
	const mockWriePoolSize = 600000000
	id := getMockAllocationId(i)
	sa := &StorageAllocation{
		ID:              id,
		DataShards:      viper.GetInt(sc.NumBlobbersPerAllocation) / 2,
		ParityShards:    viper.GetInt(sc.NumBlobbersPerAllocation) / 2,
		Size:            viper.GetInt64(sc.StorageMinAllocSize),
		Expiration:      benchAllocationExpire(balances.GetTransaction().CreationDate),
		Owner:           clients[cIndex],
		OwnerPublicKey:  publicKey,
		ReadPriceRange:  PriceRange{0, currency.Coin(viper.GetFloat64(sc.StorageMaxReadPrice) * 1e10)},
		WritePriceRange: PriceRange{0, currency.Coin(viper.GetFloat64(sc.StorageMaxWritePrice) * 1e10)},
		StartTime:       balances.GetTransaction().CreationDate,
		Stats: &StorageAllocationStats{
			UsedSize:                  viper.GetInt64(sc.StorageMinAllocSize) / 2,
			NumWrites:                 1,
			NumReads:                  1,
			TotalChallenges:           10,
			OpenChallenges:            2,
			SuccessChallenges:         6,
			FailedChallenges:          2,
			LastestClosedChallengeTxn: "latest closed challenge transaction:" + id,
		},
		TimeUnit:  viper.GetDuration(sc.TimeUnit),
		Finalized: i == mockFinalizedAllocationIndex,
		WritePool: 2e10,
	}

	startBlobbers := getMockBlobberBlockFromAllocationIndex(i)
	for j := 0; j < viper.GetInt(sc.NumBlobbersPerAllocation); j++ {
		bIndex := startBlobbers + j
		bId := getMockBlobberId(bIndex)
		ba := BlobberAllocation{
			BlobberID:    bId,
			AllocationID: sa.ID,
			Size:         viper.GetInt64(sc.StorageMinAllocSize),
			Stats: &StorageAllocationStats{
				UsedSize:                  (sa.Stats.UsedSize) / 2,
				NumWrites:                 sa.Stats.NumWrites,
				NumReads:                  sa.Stats.NumReads,
				TotalChallenges:           (sa.Stats.TotalChallenges) / 2,
				OpenChallenges:            (sa.Stats.OpenChallenges) / 2,
				SuccessChallenges:         (sa.Stats.SuccessChallenges) / 2,
				FailedChallenges:          (sa.Stats.FailedChallenges) / 2,
				LastestClosedChallengeTxn: sa.Stats.LastestClosedChallengeTxn,
			},
			Terms:                         getMockBlobberTerms(),
			AllocationRoot:                encryption.Hash("allocation root"),
			LastWriteMarker:               &WriteMarker{},
			LatestFinalizedChallCreatedAt: 0,
		}
		sa.BlobberAllocs = append(sa.BlobberAllocs, &ba)

		blobAllocPart, err := partitionsBlobberAllocations(bId, balances)
		if err != nil {
			log.Fatal("add blob alloc partition", err)
		}
		if err := blobAllocPart.Add(balances, &BlobberAllocationNode{ID: sa.ID}); err != nil {
			log.Fatal("add blob alloc node", err)
		}
		if err := blobAllocPart.Save(balances); err != nil {
			log.Fatal("save blob alloc part", err)
		}
	}

	if _, err := balances.InsertTrieNode(sa.GetKey(ADDRESS), sa); err != nil {
		log.Fatal(err)
	}

	if viper.GetBool(sc.EventDbEnabled) {
		allocationTerms := make([]event.AllocationBlobberTerm, 0)
		for _, b := range sa.BlobberAllocs {
			allocationTerms = append(allocationTerms, event.AllocationBlobberTerm{
				BlobberID:        b.BlobberID,
				AllocationIdHash: b.AllocationID,
				ReadPrice:        int64(b.Terms.ReadPrice),
				WritePrice:       int64(b.Terms.WritePrice),
			})
		}

		allocationDb := event.Allocation{
			AllocationID:             sa.ID,
			DataShards:               sa.DataShards,
			ParityShards:             sa.ParityShards,
			Size:                     sa.Size,
			Expiration:               int64(sa.Expiration),
			Owner:                    sa.Owner,
			OwnerPublicKey:           sa.OwnerPublicKey,
			StartTime:                int64(sa.StartTime),
			UsedSize:                 sa.Stats.UsedSize,
			TimeUnit:                 int64(sa.TimeUnit),
			NumWrites:                sa.Stats.NumWrites,
			NumReads:                 sa.Stats.NumReads,
			TotalChallenges:          sa.Stats.TotalChallenges,
			OpenChallenges:           sa.Stats.OpenChallenges,
			FailedChallenges:         sa.Stats.FailedChallenges,
			LatestClosedChallengeTxn: sa.Stats.LastestClosedChallengeTxn,
			Terms:                    allocationTerms,
		}
		if err := eventDb.Store.Get().Create(&allocationDb).Error; err != nil {
			log.Fatal(err)
		}
	}
}

func AddMockChallenges(
	validatorIds []string,
	blobbers []*StorageNode,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) {
	numAllocations := viper.GetInt(sc.NumAllocations)
	allocationChall := make([]AllocationChallenges, numAllocations)

	challengeReadyBlobbersPart, err := partitions.CreateIfNotExists(balances,
		ALL_CHALLENGE_READY_BLOBBERS_KEY, allChallengeReadyBlobbersPartitionSize)
	if err != nil {
		log.Fatal(err)
	}

	var (
		numAllocBlobbers        = viper.GetInt(sc.NumBlobbersPerAllocation)
		numValidators           = numAllocBlobbers / 2
		numChallengesPerBlobber = viper.GetInt(sc.NumChallengesBlobber)
		numAllocs               = viper.GetInt(sc.NumAllocations)
	)

	challenges := make([]*StorageChallenge, 0, numAllocs*numAllocBlobbers*numChallengesPerBlobber)

	for i := 0; i < numAllocs; i++ {
		startBlobbers := getMockBlobberBlockFromAllocationIndex(i)
		blobInd := rand.Intn(startBlobbers + 1)

		// When the challenge_response is called, the challenge shouldn't be expired
		// the challenge expires when the difference between the current round and the creation round of the challenge is more than max_challenge_completion_rounds from the configuration
		// Current round for the benchmark tests is set to NumBlocks when executing the tests.
		// Thus, creation round of the challenge should be less than NumBlocks by a value less than max_challenge_completion_rounds
		round := viper.GetInt(sc.NumBlocks) - (i % viper.GetInt(sc.StorageMaxChallengeCompletionRounds))

		cs := setupMockChallenge(
			numChallengesPerBlobber,
			numValidators,
			getMockAllocationId(i),
			validatorIds,
			blobbers[blobInd],
			&allocationChall[i],
			eventDb,
			balances,
			round,
		)
		challenges = append(challenges, cs...)
	}
	blobAlloc := make(map[string]map[string]*AllocOpenChallenge)

	// adding blobber challenges and blobber challenge partition
	blobbersMap := make(map[string]struct{})
	for _, ch := range challenges {
		if _, ok := blobbersMap[ch.BlobberID]; ok {
			continue
		}

		err := challengeReadyBlobbersPart.Add(balances, &ChallengeReadyBlobber{
			BlobberID: ch.BlobberID,
		})
		if err != nil {
			panic(err)
		}

		blobbersMap[ch.BlobberID] = struct{}{}
	}

	err = challengeReadyBlobbersPart.Save(balances)
	if err != nil {
		panic(err)
	}

	// adding allocation challenges
	for _, ch := range allocationChall {
		_, err := balances.InsertTrieNode(ch.GetKey(ADDRESS), &ch)
		if err != nil {
			panic(err)
		}
		for _, oc := range ch.OpenChallenges {
			if _, ok := blobAlloc[oc.BlobberID]; !ok {
				blobAlloc[oc.BlobberID] = make(map[string]*AllocOpenChallenge)
			}
			blobAlloc[oc.BlobberID][ch.AllocationID] = oc
		}
	}
}

func AddMockReadPools(clients []string, eventDb *event.EventDb, balances cstate.StateContextI) {
	rps := make([]*readPool, len(clients))
	for i := range clients {
		rps[i] = &readPool{
			Balance: 10 * 1e10,
		}
	}
	var edbRps []event.ReadPool
	for i := 0; i < len(rps); i++ {
		if _, err := balances.InsertTrieNode(readPoolKey(ADDRESS, clients[i]), rps[i]); err != nil {
			log.Fatal(err)
		}
		if viper.GetBool(sc.EventDbEnabled) {
			edbRps = append(edbRps, event.ReadPool{
				UserID:  clients[i],
				Balance: rps[i].Balance,
			})
		}
	}
	if viper.GetBool(sc.EventDbEnabled) {
		err := eventDb.InsertReadPool(edbRps)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func AddMockChallengePools(eventDb *event.EventDb, balances cstate.StateContextI) {
	var challengePools []event.ChallengePool
	for i := 0; i < viper.GetInt(sc.NumAllocations); i++ {
		allocationId := getMockAllocationId(i)
		cp := newChallengePool()
		cp.TokenPool.ID = challengePoolKey(ADDRESS, allocationId)
		cp.Balance = 10
		if _, err := balances.InsertTrieNode(challengePoolKey(ADDRESS, allocationId), cp); err != nil {
			log.Fatal(err)
		}

		if viper.GetBool(sc.EventDbEnabled) {
			challengePool := event.ChallengePool{
				ID:           cp.ID,
				AllocationID: allocationId,
				Balance:      int64(cp.Balance),
				Finalized:    false,
			}
			challengePools = append(challengePools, challengePool)
		}
	}
	if len(challengePools) > 0 {
		if err := eventDb.Store.Get().Create(&challengePools).Error; err != nil {
			log.Fatal(err)
		}
	}
}

func setupMockChallenge(
	challengesPerBlobber int,
	totalValidatorsNum int,
	allocationId string,
	validatorIds []string,
	blobber *StorageNode,
	ac *AllocationChallenges,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
	index int,
) []*StorageChallenge {
	ac.AllocationID = allocationId

	if len(validatorIds) < viper.GetInt(sc.StorageValidatorsPerChallenge) {
		log.Fatalf("number of validators %d less than validators per challenge %d",
			len(validatorIds), viper.GetInt(sc.StorageValidatorsPerChallenge))
	}

	challenges := make([]*StorageChallenge, 0, challengesPerBlobber)
	bb := blobber.mustBase()
	challenge := &StorageChallenge{
		ID:              getMockChallengeId(bb.ID, allocationId),
		AllocationID:    allocationId,
		TotalValidators: totalValidatorsNum,
		BlobberID:       bb.ID,
		ValidatorIDs:    validatorIds[:viper.GetInt(sc.StorageValidatorsPerChallenge)],
		RoundCreatedAt:  int64(index) - 1,
	}
	_, err := balances.InsertTrieNode(challenge.GetKey(ADDRESS), challenge)
	if err != nil {
		log.Fatal(err)
	}
	if ac.addChallenge(challenge) {
		challenges = append(challenges, challenge)
	}

	if viper.GetBool(sc.EventDbEnabled) {
		challengeRow := event.Challenge{
			ChallengeID:    challenge.ID,
			CreatedAt:      balances.GetTransaction().CreationDate,
			AllocationID:   challenge.AllocationID,
			BlobberID:      challenge.BlobberID,
			RoundResponded: int64(index),
		}
		if err = eventDb.Store.Get().Create(&challengeRow).Error; err != nil {
			log.Fatal(err)
		}
	}

	return challenges
}

func AddMockBlobbers(
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) []*StorageNode {

	numRewardPartitionBlobbers := viper.GetInt(sc.NumRewardPartitionBlobber)
	numBlobbers := viper.GetInt(sc.NumBlobbers)
	if numRewardPartitionBlobbers > numBlobbers {
		log.Fatal("reward_partition_blobber cannot be greater than total blobbers")
	}

	partition, err := getActivePassedBlobberRewardsPartitions(balances, viper.GetInt64(sc.StorageBlockRewardTriggerPeriod))
	if err != nil {
		log.Fatal("getting active passed blobber rewards partition", err)
	}

	var sscId = StorageSmartContract{
		SmartContract: sci.NewSC(ADDRESS),
	}.ID
	var blobbers StorageNodes
	var rtvBlobbers []*StorageNode
	blobbersDb := make([]event.Blobber, 0, viper.GetInt(sc.NumBlobbers))
	for i := 0; i < viper.GetInt(sc.NumBlobbers); i++ {
		id := getMockBlobberId(i)
		const mockUsedData = 1000
		blobber := &StorageNode{}
		blobber.SetEntity(&storageNodeV2{
			Provider: provider.Provider{
				ID:              id,
				ProviderType:    spenum.Blobber,
				LastHealthCheck: balances.GetTransaction().CreationDate,
			},
			BaseURL:           getMockBlobberUrl(i),
			Terms:             getMockBlobberTerms(),
			Capacity:          viper.GetInt64(sc.StorageMinBlobberCapacity) * 10000,
			Allocated:         mockUsedData,
			PublicKey:         "",
			StakePoolSettings: getMockStakePoolSettings(id),
			NotAvailable:      false,
		})
		blobbers.Nodes.add(blobber)
		rtvBlobbers = append(rtvBlobbers, blobber)
		_, err := balances.InsertTrieNode(blobber.GetKey(), blobber)
		if err != nil {
			log.Fatal("insert blobber into mpt", err)
		}
		_, err = balances.InsertTrieNode(blobber.GetUrlKey(sscId), &datastore.NOIDField{})
		if err != nil {
			log.Fatal("insert blobber url into mpt", err)
		}

		bb := blobber.mustBase()
		if viper.GetBool(sc.EventDbEnabled) {
			blobberDb := event.Blobber{
				BaseURL:    bb.BaseURL,
				ReadPrice:  bb.Terms.ReadPrice,
				WritePrice: bb.Terms.WritePrice,
				Capacity:   bb.Capacity,
				Allocated:  bb.Allocated,
				ReadData:   bb.Allocated * 2,
				Provider: event.Provider{
					ID:              bb.ID,
					DelegateWallet:  bb.StakePoolSettings.DelegateWallet,
					NumDelegates:    bb.StakePoolSettings.MaxNumDelegates,
					ServiceCharge:   bb.StakePoolSettings.ServiceChargeRatio,
					LastHealthCheck: bb.LastHealthCheck,
				},
				ChallengesPassed:    uint64(i),
				ChallengesCompleted: uint64(i + 1),
				RankMetric:          float64(i) / (float64(i) + 1),
				NotAvailable:        bb.NotAvailable,
			}
			blobberDb.TotalStake, err = currency.ParseZCN(viper.GetFloat64(sc.StorageMaxStake) / 2)
			if err != nil {
				log.Fatal("convert currency", err)
			}
			blobbersDb = append(blobbersDb, blobberDb)
		}

		if i < numRewardPartitionBlobbers {
			err = partition.Add(balances,
				&BlobberRewardNode{
					ID:                bb.ID,
					SuccessChallenges: 10,
					WritePrice:        bb.Terms.WritePrice,
					ReadPrice:         bb.Terms.ReadPrice,
					TotalData:         sizeInGB(int64(i * 1000)),
					DataRead:          float64(i) * 0.1,
				})
			if err != nil {
				log.Fatal("add partition", err)
			}
		}
	}
	addMockBlobberSnapshots(blobbersDb, eventDb)
	if viper.GetBool(sc.EventDbEnabled) {
		if err := eventDb.Store.Get().Create(&blobbersDb).Error; err != nil {
			log.Fatal(err)
		}
	}

	err = partition.Save(balances)
	if err != nil {
		log.Fatal("Save partition", err)
	}
	return rtvBlobbers
}

func addMockBlobberSnapshots(blobbers []event.Blobber, edb *event.EventDb) {
	if edb == nil {
		return
	}
	var mockChallengesPassed = viper.GetUint64(sc.EventDbAggregatePeriod)
	var mockChallengesCompleted = viper.GetUint64(sc.EventDbAggregatePeriod) + 1
	const mockInactiveRounds = 17
	var aggregates []event.BlobberAggregate
	for _, blobber := range blobbers {
		for i := sc.GetOldestAggregateRound(); i < viper.GetInt64(sc.NumBlocks); i++ {
			aggregate := event.BlobberAggregate{
				Round:               int64(i),
				BlobberID:           blobber.ID,
				WritePrice:          blobber.WritePrice,
				Capacity:            blobber.Capacity,
				Allocated:           blobber.Allocated,
				SavedData:           blobber.SavedData,
				ReadData:            blobber.ReadData,
				OffersTotal:         blobber.OffersTotal,
				TotalStake:          blobber.TotalStake,
				ChallengesPassed:    mockChallengesPassed * uint64(i),
				ChallengesCompleted: mockChallengesCompleted * uint64(i),
				InactiveRounds:      mockInactiveRounds,
			}
			aggregates = append(aggregates, aggregate)
		}
	}

	res := edb.Store.Get().Create(&aggregates)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
}

func addMockValidatorSnapshots(validators []event.Validator, edb *event.EventDb) {
	if edb == nil {
		return
	}
	var aggregates []event.ValidatorAggregate
	for _, validator := range validators {
		for i := sc.GetOldestAggregateRound(); i < viper.GetInt64(sc.NumBlocks); i++ {
			aggregate := event.ValidatorAggregate{
				Round:         i,
				ValidatorID:   validator.ID,
				TotalStake:    validator.TotalStake,
				TotalRewards:  validator.GetTotalRewards(),
				ServiceCharge: validator.GetServiceCharge(),
			}
			aggregates = append(aggregates, aggregate)
		}
	}

	res := edb.Store.Get().Create(&aggregates)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
}

func AddMockSnapshots(edb *event.EventDb) {
	if edb == nil {
		return
	}
	var snapshots []event.Snapshot
	for i := sc.GetOldestAggregateRound(); i < viper.GetInt64(sc.NumBlocks); i++ {
		snapshot := event.Snapshot{
			Round:                i,
			TotalMint:            i + 10,
			TotalChallengePools:  int64(currency.Coin(i + (1 * 1e10))),
			ActiveAllocatedDelta: i,
			TotalStaked:          int64(currency.Coin(i * (0.001 * 1e10))),
			SuccessfulChallenges: (i - 1) / 2,
			TotalChallenges:      i - 1,
			ZCNSupply:            100000 * int64(i+10),
			AllocatedStorage:     i * 1024,
			MaxCapacityStorage:   i * 10240,
			StakedStorage:        i * 512,
			UsedStorage:          i * 256,
			ClientLocks:          int64(currency.Coin(i * (0.0001 * 1e10))),
		}
		snapshots = append(snapshots, snapshot)
	}

	rounds := make([]int64, 0, len(snapshots))
	for _, snapshot := range snapshots {
		rounds = append(rounds, snapshot.Round)
	}

	log.Printf("addMockSnapshots rounds: %v\n", rounds)

	res := edb.Store.Get().Create(&snapshots)
	if res.Error != nil {
		log.Fatal("mock snapshot failed on create edb row", res.Error)
	}
}

func AddMockValidators(
	ids, publicKeys []string,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) []*ValidationNode {
	var sscId = StorageSmartContract{
		SmartContract: sci.NewSC(ADDRESS),
	}.ID

	valParts, err := partitions.CreateIfNotExists(balances, ALL_VALIDATORS_KEY, allValidatorsPartitionSize)
	if err != nil {
		panic(err)
	}
	if len(ids) != len(publicKeys) {
		log.Fatalf("length validator ids %d does not equal length of public keys %d",
			len(ids), len(publicKeys))
	}
	validatorNodes := make([]*ValidationNode, 0, len(ids))
	validators := make([]event.Validator, 0, len(ids))
	for i, id := range ids {
		url := getMockValidatorUrl(id)
		validator := &ValidationNode{
			Provider: provider.Provider{
				ID:           id,
				ProviderType: spenum.Validator,
			},
			BaseURL:           url,
			PublicKey:         publicKeys[i],
			StakePoolSettings: getMockStakePoolSettings(id),
		}
		_, err := balances.InsertTrieNode(validator.GetKey(sscId), validator)
		if err != nil {
			panic(err)
		}
		validatorNodes = append(validatorNodes, validator)
		vpn := ValidationPartitionNode{
			Id:  id,
			Url: id + ".com",
		}
		if viper.GetBool(sc.EventDbEnabled) {
			validator := event.Validator{
				BaseUrl:   validator.BaseURL,
				PublicKey: publicKeys[i],
				Provider: event.Provider{
					ID:             validator.ID,
					DelegateWallet: validator.StakePoolSettings.DelegateWallet,
					NumDelegates:   validator.StakePoolSettings.MaxNumDelegates,
					ServiceCharge:  validator.StakePoolSettings.ServiceChargeRatio,
				},
			}
			validators = append(validators, validator)
		}

		if err := valParts.Add(balances, &vpn); err != nil {
			panic(err)
		}
	}

	if viper.GetBool(sc.EventDbEnabled) {
		addMockValidatorSnapshots(validators, eventDb)
		if err := eventDb.Store.Get().Create(&validators).Error; err != nil {
			log.Fatal(err)
		}
	}

	err = valParts.Save(balances)
	if err != nil {
		panic(err)
	}
	return validatorNodes
}

func GetMockBlobberStakePools(
	clients []string,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) []*stakePool {
	sps := make([]*stakePool, 0, viper.GetInt(sc.NumBlobbers))
	for i := 0; i < viper.GetInt(sc.NumBlobbers); i++ {
		bId := getMockBlobberId(i)
		sp := &stakePool{
			StakePool: &stakepool.StakePool{
				Pools:    make(map[string]*stakepool.DelegatePool),
				Reward:   0,
				Settings: getMockStakePoolSettings(bId),
				Minter:   cstate.MinterStorage,
			},
			TotalOffers: currency.Coin(100000),
		}
		for j := 0; j < viper.GetInt(sc.NumBlobberDelegates)-1; j++ {
			id := getMockBlobberStakePoolId(i, j, clients)
			clientIndex := (i&len(clients) + j) % len(clients)
			sp.Pools[id] = &stakepool.DelegatePool{
				Reward: 10,
			}
			bal := currency.Coin(viper.GetInt64(sc.StorageMaxStake) * 1e10 / 1000)
			sp.Pools[id].Balance = bal
			sp.Pools[id].DelegateID = clients[clientIndex]

			if viper.GetBool(sc.EventDbEnabled) {
				dp := event.DelegatePool{
					PoolID:       id,
					ProviderType: spenum.Blobber,
					ProviderID:   bId,
					DelegateID:   sp.Pools[id].DelegateID,
					Balance:      sp.Pools[id].Balance,
					Reward:       10,
					TotalReward:  10,
					TotalPenalty: 0,
					Status:       spenum.Active,
					RoundCreated: 1,
					StakedAt:     sp.Pools[id].StakedAt,
				}
				if err := eventDb.Store.Get().Create(&dp).Error; err != nil {
					log.Fatal(err)
				}
			}
		}
		sps = append(sps, sp)
	}
	return sps
}

func GetMockValidatorStakePools(
	validatorIds []string,
	balances cstate.StateContextI,
) {
	if len(validatorIds) < viper.GetInt(sc.NumValidators) {
		log.Fatalf("length of validator ids %d less than the num of validaotrs %d",
			len(validatorIds), viper.GetInt(sc.NumValidators))
	}

	for i := 0; i < viper.GetInt(sc.NumValidators); i++ {
		bId := validatorIds[i]
		sp := &stakePool{
			StakePool: &stakepool.StakePool{
				Pools:    make(map[string]*stakepool.DelegatePool),
				Reward:   0,
				Settings: getMockStakePoolSettings(bId),
			},
		}
		for j := 0; j < viper.GetInt(sc.NumBlobberDelegates); j++ {
			id := getMockValidatorStakePoolId(validatorIds[i], j)
			sp.Pools[id] = &stakepool.DelegatePool{}
			sp.Pools[id].Balance = currency.Coin(viper.GetInt64(sc.StorageMaxStake) * 1e10 / 2)
			err := sp.Save(spenum.Validator, validatorIds[i], balances)
			if err != nil {
				panic(err)
			}
		}
	}
}

func SaveMockStakePools(
	sps []*stakePool,
	balances cstate.StateContextI,
) {
	for i, sp := range sps {
		bId := getMockBlobberId(i)
		err := sp.Save(spenum.Blobber, bId, balances)
		if err != nil {
			panic(err)
		}
	}
}

func AddMockFreeStorageAssigners(
	clients []string,
	keys []string,
	balances cstate.StateContextI,
) {
	var sscId = StorageSmartContract{
		SmartContract: sci.NewSC(ADDRESS),
	}.ID
	for i := 0; i < viper.GetInt(sc.NumFreeStorageAssigners); i++ {
		_, err := balances.InsertTrieNode(
			freeStorageAssignerKey(sscId, clients[i]),
			&freeStorageAssigner{
				ClientId:        clients[i],
				PublicKey:       keys[i],
				IndividualLimit: currency.Coin(viper.GetFloat64(sc.StorageMaxIndividualFreeAllocation) * 1e10),
				TotalLimit:      currency.Coin(viper.GetFloat64(sc.StorageMaxTotalFreeAllocation) * 1e10),
				CurrentRedeemed: 0,
				RedeemedNonces:  []int64{},
			},
		)
		if err != nil {
			panic(err)
		}
	}
}

func AddMockReadMarkers(
	clients, publicKeys []string,
	eventDb *event.EventDb,
	balances cstate.StateContextI,
) {
	numReadMarkers := viper.GetInt(sc.NumReadMarkersAllocation)
	numAllocations := viper.GetInt(sc.NumAllocations)
	for i := 0; i < numAllocations; i++ {
		var rms []event.ReadMarker
		for j := 0; j < numReadMarkers; j++ {
			round := getMockReadMarkerRound(i, j)
			if round >= viper.GetInt64(sc.NumBlocks) {
				break
			}
			client := getMockOwnerFromAllocationIndex(i, len(clients))
			mockReadCounter := int64(2)
			rm := ReadMarker{
				ClientID:        clients[client],
				ClientPublicKey: publicKeys[client],
				BlobberID:       getMockBlobberId(getMockBlobberBlockFromAllocationIndex(i)),
				AllocationID:    getMockAllocationId(i),
				OwnerID:         clients[client],
				ReadCounter:     mockReadCounter,
			}
			commitRead := &ReadConnection{
				ReadMarker: &rm,
			}
			_, err := balances.InsertTrieNode(commitRead.GetKey(ADDRESS), commitRead)
			if err != nil {
				panic(err)
			}

			if viper.GetBool(sc.EventDbEnabled) {
				mockRound := getMockReadMarkerRound(i, j)
				readMarker := event.ReadMarker{
					ClientID:      rm.ClientID,
					BlobberID:     rm.BlobberID,
					AllocationID:  rm.AllocationID,
					TransactionID: benchmark.GetMockTransactionHash(mockRound, j),
					OwnerID:       rm.OwnerID,
					ReadCounter:   rm.ReadCounter,
					ReadSize:      100,
					BlockNumber:   mockRound,
				}
				rms = append(rms, readMarker)
			}
		}
		if err := eventDb.Store.Get().Create(&rms).Error; err != nil {
			log.Fatal(err)
		}
	}
}

func AddMockWriteMarkers(
	clients []string,
	eventDb *event.EventDb,
) {
	numWriteRedeemAllocation := viper.GetInt(sc.NumWriteRedeemAllocation)
	numAllocations := viper.GetInt(sc.NumAllocations)
	for i := 0; i < numAllocations; i++ {
		var wms []event.WriteMarker
		for j := 0; j < numWriteRedeemAllocation; j++ {
			client := getMockOwnerFromAllocationIndex(i, len(clients))
			if viper.GetBool(sc.EventDbEnabled) {
				mockRound := getMockWriteMarkerRound(i, j)
				if mockRound >= viper.GetInt64(sc.NumBlocks) {
					break
				}
				writeMarker := event.WriteMarker{
					ClientID:       clients[client],
					BlobberID:      getMockBlobberId(getMockBlobberBlockFromAllocationIndex(i)),
					AllocationID:   getMockAllocationId(i),
					TransactionID:  benchmark.GetMockTransactionHash(mockRound, j),
					AllocationRoot: "mock allocation root",
					BlockNumber:    mockRound,
					Size:           100,
				}
				wms = append(wms, writeMarker)
			}
		}
		if viper.GetBool(sc.EventDbEnabled) {
			if err := eventDb.Store.Get().Create(&wms).Error; err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getMockWriteMarkerRound(allocationIndex, writeMarkerIndex int) int64 {
	return int64(allocationIndex + writeMarkerIndex*viper.GetInt(sc.NumRoundsBetweenWrites) +
		viper.GetInt(sc.NumRoundsBetweenWrites))
}

func getMockReadMarkerRound(allocationIndex, readMarkerIndex int) int64 {
	return int64(allocationIndex + readMarkerIndex*viper.GetInt(sc.NumRoundsBetweenReads) +
		2*viper.GetInt(sc.NumRoundsBetweenWrites))
}

func getMockBlobberTerms() Terms {
	return Terms{
		ReadPrice:  currency.Coin(0.01 * 1e10),
		WritePrice: currency.Coin(0.01 * 1e10),
	}
}

func getMockStakePoolSettings(blobber string) stakepool.Settings {
	return stakepool.Settings{
		DelegateWallet:     blobber,
		MaxNumDelegates:    viper.GetInt(sc.NumBlobberDelegates),
		ServiceChargeRatio: viper.GetFloat64(sc.StorageMaxCharge),
	}
}

func getMockBlobberStakePoolId(blobber, stake int, clients []string) string {
	index := viper.GetInt(sc.NumBlobberDelegates)*blobber + stake
	clinetIndex := index % len(clients)
	clinetIndex = clinetIndex
	return clients[index%len(clients)]
}

func getMockValidatorStakePoolId(validator string, stake int) string {
	return encryption.Hash(validator + ":pool:" + strconv.Itoa(stake))
}

func getMockBlobberId(index int) string {
	return encryption.Hash("mockBlobber_" + strconv.Itoa(index))
}

func getMockBlobberUrl(index int) string {
	return getMockBlobberId(index) + ".com"
}

func getMockValidatorUrl(id string) string {
	return id + ".com"
}

func getMockAllocationId(allocation int) string {
	return encryption.Hash("mock allocation id" + strconv.Itoa(allocation))
}

func getMockOwnerFromAllocationIndex(allocation, numClinets int) int {
	return allocation % (numClinets - 1 - viper.GetInt(sc.NumAllocationPayerPools))
}

func getMockBlobberBlockFromAllocationIndex(i int) int {
	return i % (viper.GetInt(sc.NumBlobbers) - viper.GetInt(sc.NumBlobbersPerAllocation))
}

func getMockChallengeId(blobberID, allocationId string) string {
	return encryption.Hash("challenge" + allocationId + blobberID)
}

func SetMockConfig(
	balances cstate.StateContextI,
) (conf *Config) {
	conf = newConfig()

	conf.TimeUnit = viper.GetDuration(sc.TimeUnit)
	conf.ChallengeEnabled = true
	conf.MinAllocSize = viper.GetInt64(sc.StorageMinAllocSize)
	conf.MinBlobberCapacity = viper.GetInt64(sc.StorageMinBlobberCapacity)
	conf.ValidatorReward = 0.025

	conf.HealthCheckPeriod = 1 * time.Hour
	conf.BlobberSlash = 0.1
	conf.CancellationCharge = 0.2
	conf.MaxReadPrice = 100e10  // 100 tokens per GB max allowed (by 64 KB)
	conf.MaxWritePrice = 100e10 // 100 tokens per GB max allowed
	conf.MinWritePrice = 0
	conf.NumValidatorsRewarded = viper.GetInt(sc.StorageNumValidatorsRewarded)
	conf.ValidatorsPerChallenge = viper.GetInt(sc.StorageValidatorsPerChallenge)
	conf.MaxDelegates = viper.GetInt(sc.StorageMaxDelegates)
	conf.MaxChallengeCompletionRounds = viper.GetInt64(sc.StorageMaxChallengeCompletionRounds)
	conf.MaxCharge = viper.GetFloat64(sc.StorageMaxCharge)
	conf.MinStake = currency.Coin(viper.GetInt64(sc.StorageMinStake) * 1e10)
	conf.MaxStake = currency.Coin(viper.GetInt64(sc.StorageMaxStake) * 1e10)
	conf.MaxTotalFreeAllocation = currency.Coin(viper.GetInt64(sc.StorageMaxTotalFreeAllocation) * 1e10)
	conf.MaxIndividualFreeAllocation = currency.Coin(viper.GetInt64(sc.StorageMaxIndividualFreeAllocation) * 1e10)
	conf.ReadPool = &readPoolConfig{}
	var err error
	conf.ReadPool.MinLock, err = currency.ParseZCN(viper.GetFloat64(sc.StorageReadPoolMinLock))
	if err != nil {
		panic(err)
	}
	conf.WritePool = &writePoolConfig{
		MinLock: currency.Coin(viper.GetFloat64(sc.StorageWritePoolMinLock) * 1e10),
	}
	conf.OwnerId = viper.GetString(sc.FaucetOwner)
	conf.StakePool = &stakePoolConfig{}
	conf.StakePool.KillSlash = 0.5
	conf.FreeAllocationSettings = freeAllocationSettings{
		DataShards:   viper.GetInt(sc.StorageFasDataShards),
		ParityShards: viper.GetInt(sc.StorageFasParityShards),
		Size:         viper.GetInt64(sc.StorageFasSize),
		ReadPriceRange: PriceRange{
			Min: currency.Coin(viper.GetFloat64(sc.StorageFasReadPriceMin) * 1e10),
			Max: currency.Coin(viper.GetFloat64(sc.StorageFasReadPriceMax) * 1e10),
		},
		WritePriceRange: PriceRange{
			Min: currency.Coin(viper.GetFloat64(sc.StorageFasWritePriceMin) * 1e10),
			Max: currency.Coin(viper.GetFloat64(sc.StorageFasWritePriceMax) * 1e10),
		},
		ReadPoolFraction: viper.GetFloat64(sc.StorageFasReadPoolFraction),
	}
	conf.BlockReward = new(blockReward)
	conf.BlockReward.BlockReward = currency.Coin(viper.GetFloat64(sc.StorageBlockReward) * 1e10)
	conf.BlockReward.BlockRewardChangePeriod = viper.GetInt64(sc.StorageBlockRewardChangePeriod)
	conf.BlockReward.BlockRewardChangeRatio = viper.GetFloat64(sc.StorageBlockRewardChangeRatio)
	conf.BlockReward.QualifyingStake = currency.Coin(viper.GetFloat64(sc.StorageBlockRewardQualifyingStake) * 1e10)
	conf.MaxBlobbersPerAllocation = viper.GetInt(sc.StorageMaxBlobbersPerAllocation)
	conf.BlockReward.TriggerPeriod = viper.GetInt64(sc.StorageBlockRewardTriggerPeriod)
	if err != nil {
		panic(err)
	}

	_, err = balances.InsertTrieNode(scConfigKey(ADDRESS), conf)
	if err != nil {
		panic(err)
	}
	var mockCost = 100
	conf.Cost = map[string]int{
		"cost.update_settings":           mockCost,
		"cost.read_redeem":               mockCost,
		"cost.commit_connection":         mockCost,
		"cost.new_allocation_request":    mockCost,
		"cost.update_allocation_request": mockCost,
		"cost.finalize_allocation":       mockCost,
		"cost.cancel_allocation":         mockCost,
		"cost.add_free_storage_assigner": mockCost,
		"cost.free_allocation_request":   mockCost,
		"cost.blobber_health_check":      mockCost,
		"cost.update_blobber_settings":   mockCost,
		"cost.pay_blobber_block_rewards": mockCost,
		"cost.challenge_response":        mockCost,
		"cost.generate_challenge":        mockCost,
		"cost.add_validator":             mockCost,
		"cost.update_validator_settings": mockCost,
		"cost.add_blobber":               mockCost,
		"cost.read_pool_lock":            mockCost,
		"cost.read_pool_unlock":          mockCost,
		"cost.write_pool_lock":           mockCost,
		"cost.stake_pool_lock":           mockCost,
		"cost.stake_pool_unlock":         mockCost,
		"cost.commit_settings_changes":   mockCost,
		"cost.collect_reward":            mockCost,
		"cost.kill_blobber":              mockCost,
		"cost.kill_validator":            mockCost,
		"cost.shutdown_blobber":          mockCost,
		"cost.shutdown_validator":        mockCost,
	}
	return
}
