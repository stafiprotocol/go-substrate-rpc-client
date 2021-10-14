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

package rpc

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/client"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/chain"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/state"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/system"
)

type RPC struct {
	Author *author.Author
	Chain  *chain.Chain
	State  *state.State
	System *system.System
	client client.Client
}

func NewRPC(cl client.Client) (*RPC, error) {
	return &RPC{
		Author: author.NewAuthor(cl),
		Chain:  chain.NewChain(cl),
		State:  state.NewState(cl),
		System: system.NewSystem(cl),
		client: cl,
	}, nil
}
