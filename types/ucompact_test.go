package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestUCompact_Encode(t *testing.T) {
	a := NewUCompact(big.NewInt(100))

	re, err := EncodeToBytes(a)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("re", hexutil.Encode(re))
}
