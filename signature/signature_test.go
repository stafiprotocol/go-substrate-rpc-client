// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
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

package signature_test

import (
	"crypto/rand"
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var testSecretPhrase = "little orbit comfort eyebrow talk pink flame ridge bring milk equip blood"
var testSecretSeed = "0x167d9a020688544ea246b056799d6a771e97c9da057e4d0b87024537f99177bc"
var testPubKey = "0xdc64bef918ddda3126a39a11113767741ddfdf91399f055e1d963f2ae1ec2535"
var testAddressSS58 = "5H3gKVQU7DfNFfNGkgTrD7p715jjg7QXtat8X3UxiSyw7APW"
var testKusamaAddressSS58 = "HZHyokLjagJ1KBiXPGu75B79g1yUnDiLxisuhkvCFCRrWBk"

func TestKeyRingPairFromSecretPhrase(t *testing.T) {
	p, err := KeyringPairFromSecret(testSecretPhrase, "")
	assert.NoError(t, err)

	assert.Equal(t, KeyringPair{
		URI: testSecretPhrase,
		Address: testAddressSS58,
		PublicKey: types.MustHexDecodeString(testPubKey),
	}, p)
}

func TestKeyringPairFromSecretSeed(t *testing.T) {
	p, err := KeyringPairFromSecret(testSecretSeed, "")
	assert.NoError(t, err)

	assert.Equal(t, KeyringPair{
		URI: testSecretSeed,
		Address: testAddressSS58,
		PublicKey: types.MustHexDecodeString(testPubKey),
	}, p)
}

func TestKeyringPairFromSecretSeedAndNetwork(t *testing.T) {
	p, err := KeyringPairFromSecret(testSecretSeed, "kusama")
	assert.NoError(t, err)

	assert.Equal(t, KeyringPair{
		URI: testSecretSeed,
		Address: testKusamaAddressSS58,
		PublicKey: types.MustHexDecodeString(testPubKey),
	}, p)
}

func TestSignAndVerify(t *testing.T) {
	data := []byte("hello!")

	sig, err := Sign(data, TestKeyringPairAlice.URI)
	assert.NoError(t, err)

	ok, err := Verify(data, sig, types.HexEncodeToString(TestKeyringPairAlice.PublicKey))
	assert.NoError(t, err)

	assert.True(t, ok)
}

func TestSignAndVerifyLong(t *testing.T) {
	data := make([]byte, 258)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	sig, err := Sign(data, TestKeyringPairAlice.URI)
	assert.NoError(t, err)

	ok, err := Verify(data, sig, types.HexEncodeToString(TestKeyringPairAlice.PublicKey))
	assert.NoError(t, err)

	assert.True(t, ok)
}
