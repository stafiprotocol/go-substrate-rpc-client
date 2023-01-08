package utils

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

func BlakeTwo256(dest []byte) [32]byte {
	return blake2b.Sum256(dest)
}

const Base = 10

func StringToBigint(src string) (*big.Int, bool) {
	return big.NewInt(0).SetString(src, Base)
}
