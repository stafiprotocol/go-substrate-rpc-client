package client

import (
	scale "github.com/itering/scale.go"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/submodel"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type EraUpdated struct {
	Symbol submodel.RSymbol
	OldEra uint32
	NewEra uint32
}

type Transfer struct {
	From   types.AccountID
	To     types.AccountID
	Symbol submodel.RSymbol
	Value  types.U128
}

type Minted struct {
	To     types.AccountID
	Symbol submodel.RSymbol
	Value  types.U128
}

type Burned struct {
	From   types.AccountID
	Symbol submodel.RSymbol
	Value  types.U128
}

// LiquidityUnBond(AccountId, RSymbol, Vec<u8>, u128, u128, u128, Vec<u8>),
// RawEvent::LiquidityUnBond(who, symbol, pool, value, left_value, balance, recipient)
type Unbond struct {
	From      types.AccountID
	Symbol    submodel.RSymbol
	Pool      types.Bytes
	Value     types.U128
	LeftValue types.U128
	Balance   types.U128
	Recipient types.Bytes
}

// Swap: (account, symbol, input amount, output amount, fee amount, input is fis, fis balance, rtoken balance)
// Swap(AccountId, RSymbol, u128, u128, u128, bool, u128, u128),
type RdexSwap struct {
	From          types.AccountID
	Symbol        submodel.RSymbol
	InputAmount   types.U128
	OutputAmount  types.U128
	FeeAmount     types.U128
	InputIsFis    bool
	FisBalance    types.U128
	RTokenBalance types.U128
}

// AddLiquidity: (account, symbol, fis amount, rToken amount, new total unit, add lp unit, fis balance, rtoken balance)
// AddLiquidity(AccountId, RSymbol, u128, u128, u128, u128, u128, u128),
type RdexAddLiquidity struct {
	From          types.AccountID
	Symbol        submodel.RSymbol
	FisAmount     types.U128
	RTokenAmount  types.U128
	NewTotalUnit  types.U128
	AddUnit       types.U128
	FisBalance    types.U128
	RTokenBalance types.U128
}

// RemoveLiquidity: (account, symbol, rm unit, swap unit, rm fis amount, rm rToken amount, input is fis, fis balance, rtoken balance)
// RemoveLiquidity(AccountId, RSymbol, u128, u128, u128, u128, bool, u128, u128),
type RdexRemoveLiquidity struct {
	From               types.AccountID
	Symbol             submodel.RSymbol
	RemoveUnit         types.U128
	SwapUnit           types.U128
	RemoveFisAmount    types.U128
	RemoveRTokenAmount types.U128
	InputIsFis         bool
	FisBalance         types.U128
	RTokenBalance      types.U128
}

// RawEvent::LiquidityUnBond(who, controller, value, left_value, balance)
type RFisUnbond struct {
	From      types.AccountID
	Pool      types.Bytes
	Value     types.U128
	LeftValue types.U128
	Balance   types.U128
}

// /// symbol, old_bonding_duration, new_bonding_duration
// BondingDurationUpdated(RSymbol, u32, u32),
type BondingDuration struct {
	Symbol      submodel.RSymbol
	OldDuration types.U32
	NewDuration types.U32
}

// / \[era_index, validator_payout, remainder\]
// EraPayout(EraIndex, Balance, Balance),
type EraPayout struct {
	EraIndex types.U32
	Balance  types.U128
	Balance2 types.U128
}

// RateSet(RSymbol, RateType)
type RateSet struct {
	Symbol submodel.RSymbol
	Rate   types.U64
}

type MultiEventFlow struct {
	EventId         string
	Symbol          submodel.RSymbol
	EventData       interface{}
	Block           uint64
	Index           uint32
	Threshold       uint16
	SubAccounts     []types.Bytes
	Key             *signature.KeyringPair
	Others          []types.AccountID
	OpaqueCalls     []*MultiOpaqueCall
	PaymentInfo     *rpc.PaymentQueryInfo
	NewMulCallHashs map[string]bool
	MulExeCallHashs map[string]bool
}

type EventNewMultisig struct {
	Who, ID     types.AccountID
	CallHash    types.Hash
	CallHashStr string
	TimePoint   *submodel.OptionTimePoint
	Approvals   []types.AccountID
}

type Multisig struct {
	When      types.TimePoint
	Deposit   types.U128
	Depositor types.AccountID
	Approvals []types.AccountID
}

type EventMultisigExecuted struct {
	Who, ID     types.AccountID
	TimePoint   types.TimePoint
	CallHash    types.Hash
	CallHashStr string
	Result      bool
}

type MultiCallParam struct {
	TimePoint *submodel.OptionTimePoint
	Opaque    []byte
	Extrinsic string
	CallHash  string
}

type Receive struct {
	Recipient []byte
	Value     types.UCompact
}

type Era struct {
	Type  string `json:"type"`
	Value uint32 `json:"value"`
}

type ChainEvent struct {
	ModuleId       string             `json:"module_id" `
	EventId        string             `json:"event_id" `
	EventIndex     int                `json:"event_idx"`
	ExtrinsicIndex int                `json:"extrinsic_idx"`
	Params         []scale.EventParam `json:"params"`
}

type MultiOpaqueCall struct {
	Extrinsic string
	Opaque    []byte
	CallHash  string
	TimePoint *submodel.OptionTimePoint
}

type Transaction struct {
	ExtrinsicHash  string
	CallModuleName string
	CallName       string
	Address        interface{}
	Params         []scale.ExtrinsicParam
}

type TransInfoSingle struct {
	Block      uint64
	Index      uint32
	DestSymbol submodel.RSymbol
	Info       TransInfo
}

type TransInfoList struct {
	Block      uint64
	DestSymbol submodel.RSymbol
	List       []TransInfo
}

type TransInfoKey struct {
	Symbol submodel.RSymbol
	Block  uint64
}

type TransInfo struct {
	Account  types.AccountID
	Receiver []byte
	Value    types.U128
	IsDeal   bool `json:"is_deal"`
}

type TransResultWithBlock struct {
	Symbol submodel.RSymbol
	Block  uint64
}

type TransResultWithIndex struct {
	Symbol submodel.RSymbol
	Block  uint64
	Index  uint32
}

type GetLatestDealBLockParam struct {
	Symbol submodel.RSymbol
	Block  chan uint64
}

type GetSignaturesParam struct {
	Symbol     submodel.RSymbol
	Block      uint64
	ProposalId []byte
	Signatures chan []types.Bytes
}

type GetSignaturesKey struct {
	Block      uint64
	ProposalId types.Bytes
}

type SubmitSignatureParams struct {
	Symbol     submodel.RSymbol
	Block      types.U64
	ProposalId types.Bytes
	Signature  types.Bytes
}

type MintRewardAct struct {
	Begin                  types.U32
	End                    types.U32
	Cycle                  types.U32
	RewardRate             types.U128
	TotalReward            types.U128
	LeftAmount             types.U128
	UserLimit              types.U128
	LockedBlocks           types.U32
	TotalRtokenAmount      types.U128
	TotalNativeTokenAmount types.U128
}

type EvtExecuteBondAndSwap struct {
	AccountId     types.AccountID
	Symbol        submodel.RSymbol
	BondId        types.Hash
	Amount        types.U128
	DestRecipient types.Bytes
	DestId        types.U8
}
