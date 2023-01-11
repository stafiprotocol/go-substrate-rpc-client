package client

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	commonTypes "github.com/stafiprotocol/go-substrate-rpc-client/types/common"
)

type StakingLedger struct {
	Stash          types.AccountID
	Total          types.UCompact
	Active         types.UCompact
	Unlocking      []UnlockChunk
	ClaimedRewards []uint32
}

type UnlockChunk struct {
	Value types.UCompact
	Era   types.UCompact
}

// multiaddress account info
type AccountInfo struct {
	Nonce     uint32
	Consumers uint32
	Providers uint32
	Data      struct {
		Free       types.U128
		Reserved   types.U128
		MiscFrozen types.U128
		FreeFrozen types.U128
	}
}

type Transaction struct {
	ExtrinsicHash  string
	CallModuleName string
	CallName       string
	Address        interface{}
	Params         []commonTypes.ExtrinsicParam
}
