// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2020 Stafi Protocol
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package teste2e

import (
	"fmt"
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func TestEnd2end(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode.")
	}

	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	assert.NoError(t, err)

	fmt.Println()
	fmt.Printf("Connected to node: %v\n", rpcs.Client.URL())
	fmt.Println()

	runtimeVersion, err := rpcs.State.GetRuntimeVersionLatest()
	assert.NoError(t, err)
	fmt.Printf("authoringVersion: %v\n", runtimeVersion.AuthoringVersion)
	fmt.Printf("specVersion: %v\n", runtimeVersion.SpecVersion)
	fmt.Printf("implVersion: %v\n", runtimeVersion.ImplVersion)
	fmt.Println()

	hash, err := rpcs.Chain.GetBlockHashLatest()
	assert.NoError(t, err)
	fmt.Printf("Latest block: %v\n", hash.Hex())
	fmt.Printf("\tView in Polkadot/Substrate Apps: https://polkadot.js.org/apps/#/explorer/query/%v?"+
		"rpc=wss://serinus-5.kusama.network\n", hash.Hex())
	fmt.Printf("\tView in polkascan.io: https://polkascan.io/pre/kusama-cc2/block/%v\n", hash.Hex())
	fmt.Println()

	header, err := rpcs.Chain.GetHeader(hash)
	assert.NoError(t, err)
	fmt.Printf("Block number: %v\n", header.Number)
	fmt.Printf("Parent hash: %v\n", header.ParentHash.Hex())
	fmt.Printf("State root: %v\n", header.StateRoot.Hex())
	fmt.Printf("Extrinsics root: %v\n", header.ExtrinsicsRoot.Hex())
	fmt.Println()

	block, err := rpcs.Chain.GetBlock(hash)
	assert.NoError(t, err)
	fmt.Printf("Total extrinsics: %v\n", len(block.Block.Extrinsics))
	fmt.Println()

	finHead, err := rpcs.Chain.GetFinalizedHead()
	assert.NoError(t, err)
	fmt.Printf("Last finalized block in the canon chain: %v\n", finHead.Hex())
	fmt.Println()

	meta, err := rpcs.State.GetMetadataLatest()
	assert.NoError(t, err)

	key, err := types.CreateStorageKey(meta, "Session", "Validators", nil, nil)
	assert.NoError(t, err)

	var validators []types.AccountID
	ok, err := rpcs.State.GetStorageLatest(key, &validators)
	assert.NoError(t, err)
	assert.True(t, ok)
	fmt.Printf("Current validators:\n")
	for i, v := range validators {
		fmt.Printf("\tValidator %v: %#x\n", i, v)
	}
	fmt.Println()
}
