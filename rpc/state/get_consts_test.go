package state

import (
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/client"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func TestState_GetConstWithMetadataV12(t *testing.T) {
	meta := types.NewMetadataV12()
	types.DecodeFromBytes(types.MustHexDecodeString(types.ExamplaryMetadataV12PolkadotString), meta)
	var cst types.U128
	err := state.GetConstWithMetadata(meta, "Balances", "ExistentialDeposit", &cst)
	assert.NoError(t, err)
	assert.Equal(t, "100000000000000", cst.Int.String())
}

func TestState_GetConstWithMetadataV11(t *testing.T) {
	meta := types.NewMetadataV11()
	types.DecodeFromBytes(types.MustHexDecodeString(types.ExamplaryMetadataV11SubstrateString), meta)
	//fmt.Printf("%+v\n", meta)
	var cst types.U64
	err := state.GetConstWithMetadata(meta, "Babe", "EpochDuration", &cst)
	assert.NoError(t, err)
	assert.Equal(t, 200, int(cst))
}

func TestState_GetConstWithMetadataV10(t *testing.T) {
	meta := types.NewMetadataV11()
	types.DecodeFromBytes(types.MustHexDecodeString(types.ExamplaryMetadataV10PolkadotString), meta)
	var cst types.U64
	err := state.GetConstWithMetadata(meta, "Babe", "EpochDuration", &cst)
	assert.NoError(t, err)
	assert.Equal(t, 600, int(cst))
}

func TestState_GetConst(t *testing.T) {
	cl, err := client.Connect("ws://127.0.0.1:9944")
	assert.NoError(t, err)

	var cst types.U64
	s := NewState(cl)
	err = s.GetConst("Babe", "EpochDuration", &cst)
	assert.NoError(t, err)
	assert.Equal(t, 600, int(cst))
}
