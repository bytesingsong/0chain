package storagesc

import (
	"context"
	"fmt"
	"github.com/0chain/common/core/logging"
	"go.uber.org/zap"
	"net/url"

	"0chain.net/chaincore/smartcontract"

	"github.com/rcrowley/go-metrics"

	chainstate "0chain.net/chaincore/chain/state"
	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
)

const (
	ADDRESS = "6dba10422e368813802877a85039d3985d96760ed844092319743fb3a76712d7"
	name    = "storage"

	KB = 1024      // kilobyte
	MB = 1024 * KB // megabyte
	GB = 1024 * MB // gigabyte
)

type StorageSmartContract struct {
	*sci.SmartContract
}

func NewStorageSmartContract() sci.SmartContractInterface {
	var sscCopy = &StorageSmartContract{
		SmartContract: sci.NewSC(ADDRESS),
	}
	sscCopy.setSC(sscCopy.SmartContract, &smartcontract.BCContext{})
	return sscCopy
}

func (ipsc *StorageSmartContract) GetHandlerStats(ctx context.Context, params url.Values) (interface{}, error) {
	return ipsc.SmartContract.HandlerStats(ctx, params)
}

func (ipsc *StorageSmartContract) GetExecutionStats() map[string]interface{} {
	return ipsc.SmartContractExecutionStats
}

func (ipsc *StorageSmartContract) GetCostTable(balances chainstate.StateContextI) (map[string]int, error) {
	node, err := ipsc.getConfig(balances, true)
	if err != nil {
		return map[string]int{}, err
	}
	if node.Cost == nil {
		return map[string]int{}, err
	}
	return node.Cost, nil
}

func (ssc *StorageSmartContract) setSC(sc *sci.SmartContract, _ sci.BCContextI) {
	ssc.SmartContract = sc
	// sc configurations
	ssc.SmartContractExecutionStats["update_settings"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "update_settings"), nil)
	// reading / writing
	ssc.SmartContractExecutionStats["read_redeem"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "read_redeem"), nil)
	ssc.SmartContractExecutionStats["commit_connection"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "commit_connection"), nil)
	// allocation
	ssc.SmartContractExecutionStats["new_allocation_request"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "new_allocation_request"), nil)
	ssc.SmartContractExecutionStats["update_allocation_request"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "update_allocation_request"), nil)
	ssc.SmartContractExecutionStats["finalize_allocation"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "finalize_allocation"), nil)
	ssc.SmartContractExecutionStats["cancel_allocation"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "cancel_allocation"), nil)
	ssc.SmartContractExecutionStats["free_allocation_request"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "free_allocation_request"), nil)
	// challenge
	ssc.SmartContractExecutionStats["challenge_response"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "challenge_response"), nil)
	ssc.SmartContractExecutionStats["generate_challenge"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "generate_challenge"), nil)
	// validator
	ssc.SmartContractExecutionStats["add_validator"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "add_validator (add/update SC function)"), nil)
	ssc.SmartContractExecutionStats["update_validator_settings"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "update_validator_settings"), nil)
	ssc.SmartContractExecutionStats["validator_health_check"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID,
		"validator_health_check"), nil)
	// validators stat (not function calls)
	ssc.SmartContractExecutionStats[statAddValidator] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "add_validator"), nil)
	ssc.SmartContractExecutionStats[statUpdateValidator] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "update_validator"), nil)
	ssc.SmartContractExecutionStats[statNumberOfValidators] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "number of validators"), nil)
	// blobber
	ssc.SmartContractExecutionStats["add_blobber"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "add_blobber (add/update/remove SC function)"), nil)
	ssc.SmartContractExecutionStats["update_blobber_settings"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "update_blobber_settings"), nil)
	ssc.SmartContractExecutionStats["blobber_block_rewards"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "blobber_block_rewards"), nil)
	ssc.SmartContractExecutionStats["shut-down-blobber"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "shut-down-blobber"), nil)
	ssc.SmartContractExecutionStats["kill-blobber"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "kill-blobber"), nil)
	ssc.SmartContractExecutionStats["shut-down-validator"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "shut-down-validator"), nil)
	ssc.SmartContractExecutionStats["kill-validator"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "kill-validator"), nil)
	ssc.SmartContractExecutionStats["blobber_health_check"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID,
		"blobber_health_check"), nil)

	// blobber statistic (not function calls)
	ssc.SmartContractExecutionStats[statNumberOfBlobbers] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "stat: number of blobbers"), nil)
	ssc.SmartContractExecutionStats[statAddBlobber] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "stat: add bblober"), nil)
	ssc.SmartContractExecutionStats[statUpdateBlobber] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "stat: update blobber"), nil)
	ssc.SmartContractExecutionStats[statRemoveBlobber] = metrics.GetOrRegisterCounter(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "stat: remove blobber"), nil)
	// read pool
	ssc.SmartContractExecutionStats["read_pool_lock"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "read_pool_lock"), nil)
	ssc.SmartContractExecutionStats["read_pool_unlock"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "read_pool_unlock"), nil)
	// write pool
	ssc.SmartContractExecutionStats["write_pool_lock"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "write_pool_lock"), nil)
	// stake pool
	ssc.SmartContractExecutionStats["stake_pool_lock"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "stake_pool_lock"), nil)
	ssc.SmartContractExecutionStats["stake_pool_unlock"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "stake_pool_unlock"), nil)
	ssc.SmartContractExecutionStats["pay_reward"] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", ssc.ID, "pay_reward (add/update/remove SC function)"), nil)

}

func (ssc *StorageSmartContract) GetName() string {
	return name
}

func (ssc *StorageSmartContract) GetAddress() string {
	return ADDRESS
}

// stat not belongs to SC function calls

const (
	statAddValidator       = "stat: add validator"
	statUpdateValidator    = "stat: update validator"
	statNumberOfValidators = "stat: number of validators"
	statNumberOfBlobbers   = "stat: number of blobbers"
	statAddBlobber         = "stat: add blobber"
	statUpdateBlobber      = "stat: update blobber"
	statRemoveBlobber      = "stat: remove blobber"
)

func (ssc *StorageSmartContract) statIncr(name string) {
	var (
		metric interface{}
		count  metrics.Counter
		ok     bool
	)
	if metric, ok = ssc.SmartContractExecutionStats[name]; !ok {
		return
	}
	if count, ok = metric.(metrics.Counter); !ok {
		return
	}
	count.Inc(1)
}

func (ssc *StorageSmartContract) statDecr(name string) {
	var (
		metric interface{}
		count  metrics.Counter
		ok     bool
	)
	if metric, ok = ssc.SmartContractExecutionStats[name]; !ok {
		return
	}
	if count, ok = metric.(metrics.Counter); !ok {
		return
	}
	count.Dec(1)
}

// functions execution

func (sc *StorageSmartContract) Execute(t *transaction.Transaction,
	funcName string, input []byte, balances chainstate.StateContextI) (
	resp string, err error) {

	switch funcName {

	// read/write markers

	case "read_redeem":
		if resp, err = sc.commitBlobberRead(t, input, balances); err != nil {
			return
		}

	case "commit_connection":
		resp, err = sc.commitBlobberConnection(t, input, balances)
		if err != nil {
			return
		}

	// allocations

	case "new_allocation_request":
		resp, err = sc.newAllocationRequest(t, input, balances, nil)
	case "update_allocation_request":
		resp, err = sc.updateAllocationRequest(t, input, balances)
	case "finalize_allocation":
		resp, err = sc.finalizeAllocation(t, input, balances)
	case "cancel_allocation":
		resp, err = sc.cancelAllocationRequest(t, input, balances)

	// free allocations

	case "add_free_storage_assigner":
		resp, err = sc.addFreeStorageAssigner(t, input, balances)
	case "free_allocation_request":
		resp, err = sc.freeAllocationRequest(t, input, balances)

	// blobbers

	case "add_blobber":
		resp, err = sc.addBlobber(t, input, balances)
	case "add_validator":
		resp, err = sc.addValidator(t, input, balances)
	case "blobber_health_check":
		resp, err = sc.blobberHealthCheck(t, input, balances)
	case "validator_health_check":
		resp, err = sc.validatorHealthCheck(t, input, balances)
	case "update_blobber_settings":
		resp, err = sc.updateBlobberSettings(t, input, balances)
	case "update_validator_settings":
		resp, err = sc.updateValidatorSettings(t, input, balances)
	case "blobber_block_rewards":
		err = sc.blobberBlockRewards(t, input, balances)

	case "shutdown_blobber":
		_, err = sc.shutdownBlobber(t, input, balances)
	case "kill_blobber":
		_, err = sc.killBlobber(t, input, balances)
	case "shutdown_validator":
		_, err = sc.shutdownValidator(t, input, balances)
	case "kill_validator":
		_, err = sc.killValidator(t, input, balances)

	// read_pool

	case "read_pool_lock":
		resp, err = sc.readPoolLock(t, input, balances)
	case "read_pool_unlock":
		resp, err = sc.readPoolUnlock(t, input, balances)

	// write pool

	case "write_pool_lock":
		resp, err = sc.writePoolLock(t, input, balances)

	// stake pool

	case "stake_pool_lock":
		resp, err = sc.stakePoolLock(t, input, balances)
	case "stake_pool_unlock":
		resp, err = sc.stakePoolUnlock(t, input, balances)
	case "collect_reward":
		resp, err = sc.collectReward(t, input, balances)
	case "generate_challenge":
		var conf *Config
		if conf, err = sc.getConfig(balances, true); err != nil {
			return "", fmt.Errorf("can't get SC configurations: %v", err.Error())
		}
		if conf.ChallengeEnabled {
			err = sc.generateChallenge(t, balances.GetBlock(), input, conf, balances)
			if err != nil {
				return
			}
		} else {
			return "OpenChallenges disabled in the config", nil
		}
		return "OpenChallenges generated", nil

	case "challenge_response":
		resp, err = sc.verifyChallenge(t, input, balances)

	// configurations

	case "update_settings":
		resp, err = sc.updateSettings(t, input, balances)

	case "commit_settings_changes":
		resp, err = sc.commitSettingChanges(t, input, balances)

	default:
		logging.Logger.Info("Storage function name", zap.String("function", funcName))
		processedResetStats := false
		actErr := chainstate.WithActivation(balances, "ares", func() error {
			logging.Logger.Info("Before ares", zap.String("function", funcName))
			return nil
		}, func() error {
			if funcName == "reset_blobber_stats" {
				_ = chainstate.WithActivation(balances, "athena", func() error {
					return nil
				}, func() error {
					processedResetStats = true
					return nil
				})
				resp, err = sc.resetBlobberStats(t, input, balances)
				return err
			}
			return nil
		})
		if actErr != nil || resp != "" {
			return resp, actErr
		}

		if processedResetStats {
			return
		}

		actErr = chainstate.WithActivation(balances, "artemis", func() error {
			return nil
		}, func() error {
			if funcName == "reset_allocation_stats" {
				resp, err = sc.resetAllocationStats(t, input, balances)
				return err
			}
			return nil
		})
		if actErr != nil || resp != "" {
			return resp, actErr
		}

		err = common.NewErrorf("invalid_storage_function_name",
			"Invalid storage function '%s' called", funcName)
	}

	return
}
