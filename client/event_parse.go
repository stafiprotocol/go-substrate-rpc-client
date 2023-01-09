package client

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/go-substrate-rpc-client/submodel"
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

func parseRsymbol(value interface{}) (submodel.RSymbol, error) {
	sym, ok := value.(string)
	if !ok {
		return submodel.RSymbol(""), ErrValueNotString
	}

	return submodel.RSymbol(sym), nil
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

type EvtExecuteBondAndSwap struct {
	AccountId     types.AccountID
	Symbol        submodel.RSymbol
	BondId        types.Hash
	Amount        types.U128
	DestRecipient types.Bytes
	DestId        types.U8
}

func ParseLiquidityBondAndSwapEvent(evt *submodel.ChainEvent) (*EvtExecuteBondAndSwap, error) {
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
