package client

import (
	"errors"

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
	_, err = sc.QueryStorage(config.RTokenRateModuleId, config.StorageEraRate, symBz, eraIndex, &rate)
	return
}
