package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	scale "github.com/itering/scale.go"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/utils"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

var (
	ErrValueNotStringSlice = errors.New("value not string slice")
	ErrValueNotString      = errors.New("value not string")
	ErrValueNotFloat       = errors.New("value not float64")
	ErrValueNotMap         = errors.New("value not map")
	ErrValueNotU32         = errors.New("value not u32")
)

func parseBytes(value interface{}) ([]byte, error) {
	val, ok := value.(string)
	if !ok {
		return nil, ErrValueNotString
	}

	bz, err := hexutil.Decode(utiles.AddHex(val))
	if err != nil {
		if err.Error() == hexutil.ErrSyntax.Error() {
			return []byte(val), nil
		}
		return nil, err
	}

	return bz, nil
}

func parseAccountId(value interface{}) (types.AccountID, error) {
	val, ok := value.(string)
	if !ok {
		return types.NewAccountID([]byte{}), ErrValueNotString
	}
	ac, err := hexutil.Decode(utiles.AddHex(val))
	if err != nil {
		return types.NewAccountID([]byte{}), err
	}

	return types.NewAccountID(ac), nil
}

func parseRsymbol(value interface{}) (RSymbol, error) {
	sym, ok := value.(string)
	if !ok {
		return RSymbol(""), ErrValueNotString
	}

	return RSymbol(sym), nil
}

func parseHash(value interface{}) (types.Hash, error) {
	val, ok := value.(string)
	if !ok {
		return types.NewHash([]byte{}), ErrValueNotString
	}

	hash, err := types.NewHashFromHexString(utiles.AddHex(val))
	if err != nil {
		return types.NewHash([]byte{}), err
	}

	return hash, err
}

func parseU128(value interface{}) (types.U128, error) {
	val, ok := value.(string)
	if !ok {
		return types.U128{}, ErrValueNotString
	}
	deci, err := decimal.NewFromString(val)
	if err != nil {
		return types.U128{}, err
	}
	return types.NewU128(*deci.BigInt()), nil
}

func parseU8(value interface{}) (types.U8, error) {
	val, ok := value.(float64)
	if !ok {
		return types.U8(0), ErrValueNotFloat
	}
	return types.NewU8(uint8(decimal.NewFromFloat(val).IntPart())), nil
}

func ParseVecBytes(value interface{}) ([]types.Bytes, error) {
	vals, ok := value.([]interface{})
	if !ok {
		return nil, ErrValueNotStringSlice
	}
	result := make([]types.Bytes, 0)
	for _, val := range vals {
		bz, err := parseBytes(val)
		if err != nil {
			return nil, err
		}

		result = append(result, bz)
	}

	return result, nil
}

func parseBigint(value interface{}) (*big.Int, error) {
	val, ok := value.(string)
	if !ok {
		return nil, ErrValueNotString
	}

	i, ok := utils.StringToBigint(val)
	if !ok {
		return nil, fmt.Errorf("string to bigint error: %s", val)
	}

	return i, nil
}

func parseBoolean(value interface{}) (bool, error) {
	val, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("value not bool type")
	}

	return val, nil
}

func parseU64(value interface{}) (types.U64, error) {
	val, ok := value.(float64)
	if !ok {
		return 0, ErrValueNotString
	}

	ret := types.NewU64(uint64(val))

	return ret, nil
}

func parseU32(value interface{}) (types.U32, error) {
	valueType := reflect.TypeOf(value)
	var ret types.U32
	switch valueType.String() {
	case "float64":
		val, ok := value.(float64)
		if !ok {
			return 0, ErrValueNotString
		}

		ret = types.NewU32(uint32(val))
	default:
		return 0, ErrValueNotFloat
	}

	return ret, nil
}

func ParseLiquidityBondAndSwapEvent(evt *ChainEvent) (*EvtExecuteBondAndSwap, error) {
	if len(evt.Params) != 6 {
		return nil, fmt.Errorf("LiquidityBondEventData params number not right: %d, expected: 6", len(evt.Params))
	}
	accountId, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("LiquidityBondEventData params[0] -> AccountId error: %s", err)
	}
	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("LiquidityBondEventData params[1] -> RSymbol error: %s", err)
	}
	bondId, err := parseHash(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("LiquidityBondEventData params[2] -> BondId error: %s", err)
	}
	amount, err := parseU128(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("LiquidityBondEventData params[3] -> BondId error: %s", err)
	}
	destRecipient, err := parseBytes(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("LiquidityBondEventData params[4] -> BondId error: %s", err)
	}
	destId, err := parseU8(evt.Params[5].Value)
	if err != nil {
		return nil, fmt.Errorf("LiquidityBondEventData params[5] -> BondId error: %s", err)
	}

	return &EvtExecuteBondAndSwap{
		AccountId:     accountId,
		Symbol:        symbol,
		BondId:        bondId,
		Amount:        amount,
		DestRecipient: types.NewBytes(destRecipient),
		DestId:        destId,
	}, nil
}

func EventRateSetData(evt *ChainEvent) (*RateSet, error) {
	switch len(evt.Params) {
	case 1:
		rate, err := parseU64(evt.Params[0].Value)
		if err != nil {
			return nil, fmt.Errorf("EventRateSetData params[1] -> rate error: %s", err)
		}

		return &RateSet{
			Symbol: RFIS,
			Rate:   rate,
		}, nil
	case 2:
		symbol, err := parseRsymbol(evt.Params[0].Value)
		if err != nil {
			return nil, fmt.Errorf("EventRateSetData params[0] -> RSymbol error: %s", err)
		}

		rate, err := parseU64(evt.Params[1].Value)
		if err != nil {
			return nil, fmt.Errorf("EventRateSetData params[1] -> rate error: %s", err)
		}

		return &RateSet{
			Symbol: symbol,
			Rate:   rate,
		}, nil
	}

	return nil, fmt.Errorf("EventRateSetData params number not right: %d, expected:1 or 2", len(evt.Params))
}

func EventEraPayoutData(evt *ChainEvent) (*EraPayout, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventTransferData params number not right: %d, expected: 4", len(evt.Params))
	}
	eraIndex, err := parseU32(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[0] -> from error: %s", err)
	}
	balance, err := parseBigint(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[3] -> value error: %s", err)
	}
	balance2, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[3] -> value error: %s", err)
	}

	return &EraPayout{
		EraIndex: eraIndex,
		Balance:  types.NewU128(*balance),
		Balance2: types.NewU128(*balance2),
	}, nil
}

func EventUnbondData(evt *ChainEvent) (*Unbond, error) {
	if len(evt.Params) != 7 {
		return nil, fmt.Errorf("EventUnbondData params number not right: %d, expected: 7", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[0] -> from error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[1] -> RSymbol error: %s", err)
	}

	pool, err := parseBytes(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[2] -> value error: %s", err)
	}

	value, err := parseBigint(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}

	leftValue, err := parseBigint(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}
	balance, err := parseBigint(evt.Params[5].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}
	recipient, err := parseBytes(evt.Params[6].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[6] -> value error: %s", err)
	}

	return &Unbond{
		From:      from,
		Symbol:    symbol,
		Pool:      pool,
		Value:     types.NewU128(*value),
		LeftValue: types.NewU128(*leftValue),
		Balance:   types.NewU128(*balance),
		Recipient: recipient,
	}, nil
}

func EventRdexSwapData(evt *ChainEvent) (*RdexSwap, error) {
	if len(evt.Params) != 8 {
		return nil, fmt.Errorf("EventUnbondData params number not right: %d, expected: 7", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[0] -> from error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[1] -> RSymbol error: %s", err)
	}

	inputAmount, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}
	outputAmount, err := parseBigint(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}
	feeAmount, err := parseBigint(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}

	inputIsFis, err := parseBoolean(evt.Params[5].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}

	fisBalance, err := parseBigint(evt.Params[6].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}
	rTokenBalance, err := parseBigint(evt.Params[7].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}

	return &RdexSwap{
		From:          from,
		Symbol:        symbol,
		InputAmount:   types.NewU128(*inputAmount),
		OutputAmount:  types.NewU128(*outputAmount),
		FeeAmount:     types.NewU128(*feeAmount),
		InputIsFis:    inputIsFis,
		FisBalance:    types.NewU128(*fisBalance),
		RTokenBalance: types.NewU128(*rTokenBalance),
	}, nil
}

func EventRdexAddLiquidityData(evt *ChainEvent) (*RdexAddLiquidity, error) {
	if len(evt.Params) != 8 {
		return nil, fmt.Errorf("EventUnbondData params number not right: %d, expected: 7", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[0] -> from error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[1] -> RSymbol error: %s", err)
	}

	fisAmount, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}
	rTokenAmount, err := parseBigint(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}
	newTotalUnit, err := parseBigint(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}

	addUnit, err := parseBigint(evt.Params[5].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}

	fisBalance, err := parseBigint(evt.Params[6].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}
	rTokenBalance, err := parseBigint(evt.Params[7].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}

	return &RdexAddLiquidity{
		From:          from,
		Symbol:        symbol,
		FisAmount:     types.NewU128(*fisAmount),
		RTokenAmount:  types.NewU128(*rTokenAmount),
		NewTotalUnit:  types.NewU128(*newTotalUnit),
		AddUnit:       types.NewU128(*addUnit),
		FisBalance:    types.NewU128(*fisBalance),
		RTokenBalance: types.NewU128(*rTokenBalance),
	}, nil
}

func EventRdexRemoveLiquidityData(evt *ChainEvent) (*RdexRemoveLiquidity, error) {
	if len(evt.Params) != 9 {
		return nil, fmt.Errorf("EventUnbondData params number not right: %d, expected: 7", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[0] -> from error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[1] -> RSymbol error: %s", err)
	}

	removeUnit, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}

	swapUnit, err := parseBigint(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}
	removefisAmount, err := parseBigint(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}
	removeRTokenAmount, err := parseBigint(evt.Params[5].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[3] -> value error: %s", err)
	}

	inputIsFis, err := parseBoolean(evt.Params[6].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[4] -> value error: %s", err)
	}

	fisBalance, err := parseBigint(evt.Params[7].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}
	rTokenBalance, err := parseBigint(evt.Params[8].Value)
	if err != nil {
		return nil, fmt.Errorf("EventUnbondData params[5] -> value error: %s", err)
	}

	return &RdexRemoveLiquidity{
		From:               from,
		Symbol:             symbol,
		RemoveUnit:         types.NewU128(*removeUnit),
		SwapUnit:           types.NewU128(*swapUnit),
		RemoveFisAmount:    types.NewU128(*removefisAmount),
		RemoveRTokenAmount: types.NewU128(*removeRTokenAmount),
		InputIsFis:         inputIsFis,
		FisBalance:         types.NewU128(*fisBalance),
		RTokenBalance:      types.NewU128(*rTokenBalance),
	}, nil
}

func EventRFisUnbondData(evt *ChainEvent) (*RFisUnbond, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventRFisUnbondData params number not right: %d, expected: 5", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventRFisUnbondData params[0] -> from error: %s", err)
	}

	pool, err := parseBytes(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventRFisUnbondData params[1] -> value error: %s", err)
	}

	value, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventRFisUnbondData params[2] -> value error: %s", err)
	}

	leftValue, err := parseBigint(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventRFisUnbondData params[3] -> value error: %s", err)
	}
	balance, err := parseBigint(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventRFisUnbondData params[4] -> value error: %s", err)
	}

	return &RFisUnbond{
		From:      from,
		Pool:      pool,
		Value:     types.NewU128(*value),
		LeftValue: types.NewU128(*leftValue),
		Balance:   types.NewU128(*balance),
	}, nil
}

func EventTransferData(evt *ChainEvent) (*Transfer, error) {
	if len(evt.Params) != 4 {
		return nil, fmt.Errorf("EventTransferData params number not right: %d, expected: 4", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[0] -> from error: %s", err)
	}

	to, err := parseAccountId(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[1] -> to error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[2] -> RSymbol error: %s", err)
	}
	value, err := parseBigint(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferData params[3] -> value error: %s", err)
	}

	return &Transfer{
		From:   from,
		To:     to,
		Symbol: symbol,
		Value:  types.NewU128(*value),
	}, nil
}

func EventWithdrawUnbondData(evt *ChainEvent) (*LiquidityWithdrawUnbond, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventWithdrawUnbondData params number not right: %d, expected: 4", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventWithdrawUnbondData params[0] -> from error: %s", err)
	}

	to, err := parseAccountId(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventWithdrawUnbondData params[1] -> to error: %s", err)
	}

	value, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventWithdrawUnbondData params[3] -> value error: %s", err)
	}

	return &LiquidityWithdrawUnbond{
		From:  from,
		To:    to,
		Value: types.NewU128(*value),
	}, nil
}

func EventMintedData(evt *ChainEvent) (*Minted, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventMintedData params number not right: %d, expected: 4", len(evt.Params))
	}

	to, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMintedData params[1] -> to error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMintedData params[2] -> RSymbol error: %s", err)
	}
	value, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMintedData params[3] -> value error: %s", err)
	}

	return &Minted{
		To:     to,
		Symbol: symbol,
		Value:  types.NewU128(*value),
	}, nil
}

func EventBurnedData(evt *ChainEvent) (*Burned, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventBurnedData params number not right: %d, expected: 4", len(evt.Params))
	}
	from, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBurnedData params[0] -> from error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBurnedData params[1] -> RSymbol error: %s", err)
	}
	value, err := parseBigint(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBurnedData params[2] -> value error: %s", err)
	}

	return &Burned{
		From:   from,
		Symbol: symbol,
		Value:  types.NewU128(*value),
	}, nil
}

func EventEraUpdatedData(evt *ChainEvent) (*EraUpdated, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EraPoolUpdatedData params number not right: %d, expected: 4", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[0] -> RSymbol error: %s", err)
	}

	oldEra, err := parseEra(evt.Params[1])
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[1] -> era error: %s", err)
	}
	newEra, err := parseEra(evt.Params[2])
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[2] -> era error: %s", err)
	}

	return &EraUpdated{
		Symbol: symbol,
		OldEra: oldEra.Value,
		NewEra: newEra.Value,
	}, nil
}

func EventBondingDurationData(evt *ChainEvent) (*BondingDuration, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventBondingDurationData params number not right: %d, expected: 3", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondingDurationData params[0] -> RSymbol error: %s", err)
	}

	oldDuration, err := parseU32(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondingDurationData params[1] -> old duration error: %s", err)
	}
	newDuration, err := parseU32(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondingDurationData params[2] -> new duration error: %s", err)
	}

	return &BondingDuration{
		Symbol:      symbol,
		OldDuration: oldDuration,
		NewDuration: newDuration,
	}, nil
}

func EventNewMultisigData(evt *ChainEvent) (*EventNewMultisig, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventNewMultisigData params number not right: %d, expected: 3", len(evt.Params))
	}
	who, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNewMultisig params[0] -> who error: %s", err)
	}

	id, err := parseAccountId(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNewMultisig params[1] -> id error: %s", err)
	}

	hash, err := parseHash(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNewMultisig params[2] -> hash error: %s", err)
	}

	return &EventNewMultisig{
		Who:      who,
		ID:       id,
		CallHash: hash,
	}, nil
}

func EventMultisigExecutedData(evt *ChainEvent) (*EventMultisigExecuted, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventMultisigExecuted params number not right: %d, expected: 5", len(evt.Params))
	}

	approving, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[0] -> approving error: %s", err)
	}

	tp, err := parseTimePoint(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[1] -> timepoint error: %s", err)
	}

	id, err := parseAccountId(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[2] -> id error: %s", err)
	}

	hash, err := parseHash(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[3] -> hash error: %s", err)
	}

	ok, err := parseDispatchResult(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[4] -> dispatchresult error: %s", err)
	}

	return &EventMultisigExecuted{
		Who:       approving,
		TimePoint: tp,
		ID:        id,
		CallHash:  hash,
		Result:    ok,
	}, nil
}

func parseEra(param scale.EventParam) (*Era, error) {
	bz, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	era := new(Era)
	err = json.Unmarshal(bz, era)
	if err != nil {
		return nil, err
	}

	return era, nil
}

func parseTimePoint(value interface{}) (types.TimePoint, error) {
	bz, err := json.Marshal(value)
	if err != nil {
		return types.TimePoint{}, err
	}

	var tp types.TimePoint
	err = json.Unmarshal(bz, &tp)
	if err != nil {
		return types.TimePoint{}, err
	}

	return tp, nil
}

func parseDispatchResult(value interface{}) (bool, error) {
	result, ok := value.(map[string]interface{})
	if !ok {
		return false, ErrValueNotMap
	}
	_, ok = result["Ok"]
	return ok, nil
}
