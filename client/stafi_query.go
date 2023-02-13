package client

import (
	"errors"
	"fmt"

	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

func (c *GsrpcClient) CurrentChainEra(sym RSymbol) (uint32, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return 0, err
	}

	var era uint32
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageChainEras, symBz, nil, &era)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrorValueNotExist
	}

	return era, nil
}

func (c *GsrpcClient) ActiveChangeRateLimit(sym RSymbol) (uint32, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return 0, err
	}

	var PerBill types.U32
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageActiveChangeRateLimit, symBz, nil, &PerBill)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrorValueNotExist
	}

	return uint32(PerBill), nil
}

func (c *GsrpcClient) RTokenTotalIssuance(sym RSymbol) (types.U128, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return types.U128{}, err
	}

	var issuance types.U128
	exists, err := c.QueryStorage(config.RTokenBalanceModuleId, config.StorageTotalIssuance, symBz, nil, &issuance)
	if err != nil {
		return types.U128{}, err
	}

	if !exists {
		return types.U128{}, ErrorValueNotExist
	}

	return issuance, nil
}

func (c *GsrpcClient) CurrentEraSnapshots(symbol RSymbol) ([]types.Hash, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}

	ids := make([]types.Hash, 0)
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageCurrentEraSnapShots, symBz, nil, &ids)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("storage not exit")
	}
	return ids, nil
}

func (c *GsrpcClient) ActLatestCycle(sym RSymbol) (uint32, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return 0, err
	}

	var cycle uint32
	exists, err := c.QueryStorage(config.RClaimModuleId, config.StorageActLatestCycle, symBz, nil, &cycle)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrorValueNotExist
	}

	return cycle, nil
}

func (c *GsrpcClient) REthActLatestCycle() (uint32, error) {

	var cycle uint32
	exists, err := c.QueryStorage(config.RClaimModuleId, config.StorageREthActLatestCycle, nil, nil, &cycle)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrorValueNotExist
	}

	return cycle, nil
}

func (c *GsrpcClient) Act(sym RSymbol, cycle uint32) (*MintRewardAct, error) {
	key := struct {
		Symbol RSymbol
		Cycle  uint32
	}{
		sym,
		cycle,
	}
	keyBz, err := types.EncodeToBytes(key)
	if err != nil {
		return nil, err
	}

	act := new(MintRewardAct)

	exists, err := c.QueryStorage(config.RClaimModuleId, config.StorageActs, keyBz, nil, act)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrorValueNotExist
	}

	return act, nil
}
func (c *GsrpcClient) RethAct(cycle uint32) (*MintRewardAct, error) {
	cycleBz, err := types.EncodeToBytes(cycle)
	if err != nil {
		return nil, err
	}

	act := new(MintRewardAct)

	exists, err := c.QueryStorage(config.RClaimModuleId, config.StorageREthActs, cycleBz, nil, act)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrorValueNotExist
	}

	return act, nil
}

func (sc *GsrpcClient) GetEraRate(symbol RSymbol, era uint32) (rate uint64, err error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return 0, err
	}
	eraIndex, err := types.EncodeToBytes(types.NewU32(era))
	if err != nil {
		return 0, err
	}
	exists, err := sc.QueryStorage(config.RTokenRateModuleId, config.StorageEraRate, symBz, eraIndex, &rate)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrorValueNotExist
	}
	return rate, nil
}

func (sc *GsrpcClient) GetReceiver() (*types.AccountID, error) {
	ac := new(types.AccountID)
	exists, err := sc.QueryStorage(config.RTokenLedgerModuleId, config.StorageReceiver, nil, nil, ac)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrorValueNotExist
	}
	return ac, nil
}

func (sc *GsrpcClient) GetRFisReceiver() (*types.AccountID, error) {
	ac := new(types.AccountID)
	exists, err := sc.QueryStorage(config.RFisModuleId, config.StorageReceiver, nil, nil, ac)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrorValueNotExist
	}
	return ac, nil
}

func (gc *GsrpcClient) GetREthCurrentCycle() (uint32, error) {
	var cycle uint32
	exists, err := gc.QueryStorage(config.RClaimModuleId, config.StorageREthActCurrentCycle, nil, nil, &cycle)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrorValueNotExist
	}

	return cycle, nil
}

func (gc *GsrpcClient) MintTxHashExist(txHash types.Bytes) (bool, error) {
	txHashBytes, err := types.EncodeToBytes(txHash)
	if err != nil {
		return false, err
	}
	var txExists bool
	exists, err := gc.QueryStorage(config.RClaimModuleId, config.StorageMintTxHashExist, txHashBytes, nil, &txExists)
	if err != nil {
		return false, err
	}
	if !exists {
		return exists, nil
	}

	return txExists, nil
}

// when req reth info
// update timestamp in database every time.
// need send tx to stafi when:
// 1 current cycle begin < now block < current cycle end
// or
// 2  current+x cycle begin < now < current+x cycle end
func (gc *GsrpcClient) CurrentRethNeedSeed() (bool, error) {

	currentCycleExist := false
	currentCycle, err := gc.GetREthCurrentCycle()
	if err != nil {
		if err != ErrorValueNotExist {
			return false, err
		}
	} else {
		currentCycleExist = true
	}

	blockNumber, err := gc.GetLatestBlockNumber()
	if err != nil {
		return false, err
	}
	latestCycle, err := gc.REthActLatestCycle()
	if err != nil {
		if err == ErrorValueNotExist {
			return false, nil
		}
		return false, err
	}

	if latestCycle == 0 {
		return false, nil
	}

	if currentCycleExist && currentCycle > 0 {
		currentAct, err := gc.RethAct(currentCycle)
		if err != nil {
			return false, err
		}
		//case 1
		if uint64(currentAct.Begin) <= blockNumber && blockNumber <= uint64(currentAct.End) {
			return true, nil
		}
	}

	beginCycle := 1
	if currentCycleExist {
		beginCycle = int(currentCycle) + 1
	}

	for i := beginCycle; i <= int(latestCycle); i++ {
		act, err := gc.RethAct(uint32(i))
		if err != nil {
			if err == ErrorValueNotExist {

				return false, fmt.Errorf("cycle: %d err: %s", i, err)
			}
			return false, err
		}

		if act.Begin > types.U32(blockNumber) {
			break
		}
		//case 2
		if uint64(act.Begin) <= blockNumber && blockNumber <= uint64(act.End) {
			return true, nil
		}
	}
	return false, nil
}
