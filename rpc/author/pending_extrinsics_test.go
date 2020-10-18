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

package author

import (
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthor_PendingExtrinsics(t *testing.T) {
	res, err := author.PendingExtrinsics()
	assert.NoError(t, err)
	assert.Equal(t, []types.Extrinsic{types.Extrinsic{Version: 0x84, Signature: types.ExtrinsicSignatureV4{Signer: types.Address{IsAccountID: true, AsAccountID: types.AccountID{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, IsAccountIndex: false, AsAccountIndex: 0x0}, Signature: types.MultiSignature{IsEd25519: true, AsEd25519: types.Signature{0xa0, 0x23, 0xbb, 0xe8, 0x83, 0x40, 0x5b, 0x5f, 0xac, 0x2a, 0xa1, 0x14, 0x9, 0x3f, 0xcf, 0x3d, 0x8, 0x2, 0xd2, 0xf3, 0xd3, 0x71, 0x5e, 0x9, 0x12, 0x9b, 0x0, 0xa4, 0xbf, 0x74, 0x10, 0x48, 0xca, 0xf5, 0x3d, 0x8c, 0x7d, 0x97, 0xe8, 0x72, 0xca, 0xa7, 0x3, 0xe7, 0xd0, 0x4f, 0x17, 0x4a, 0x4e, 0x2e, 0xd4, 0xac, 0xad, 0xee, 0x41, 0x73, 0xa8, 0xb6, 0xba, 0xb7, 0xe4, 0x5c, 0xa, 0x6}, IsSr25519: false, AsSr25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, IsEcdsa: false, AsEcdsa: types.Bytes(nil)}, Era: types.ExtrinsicEra{IsImmortalEra: true, IsMortalEra: false, AsMortalEra: types.MortalEra{First: 0x0, Second: 0x0}}, Nonce: types.NewUCompactFromUInt(0x3), Tip: types.NewUCompactFromUInt(0x0)}, Method: types.Call{CallIndex: types.CallIndex{SectionIndex: 0x6, MethodIndex: 0x0}, Args: types.Args{0xff, 0x8e, 0xaf, 0x4, 0x15, 0x16, 0x87, 0x73, 0x63, 0x26, 0xc9, 0xfe, 0xa1, 0x7e, 0x25, 0xfc, 0x52, 0x87, 0x61, 0x36, 0x93, 0xc9, 0x12, 0x90, 0x9c, 0xb2, 0x26, 0xaa, 0x47, 0x94, 0xf2, 0x6a, 0x48, 0xe5, 0x6c}}}}, res) //nolint:lll,dupl
}
