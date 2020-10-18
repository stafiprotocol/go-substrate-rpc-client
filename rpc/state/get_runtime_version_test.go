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

package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_GetRuntimeVersionLatest(t *testing.T) {
	rv, err := state.GetRuntimeVersionLatest()
	assert.NoError(t, err)
	assert.Equal(t, &mockSrv.runtimeVersion, rv)
}

func TestState_GetRuntimeVersion(t *testing.T) {
	rv, err := state.GetRuntimeVersion(mockSrv.blockHashLatest)
	assert.NoError(t, err)
	assert.Equal(t, &mockSrv.runtimeVersion, rv)
}
