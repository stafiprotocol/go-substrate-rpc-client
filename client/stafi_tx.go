package client

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

func (gc *GsrpcClient) UpdateRethClaimInfo(txHashs, pubkeys [][]byte, mintValues, nativeValues []types.U128) error {
	ext, err := gc.NewUnsignedExtrinsic(config.MethodUpdateRethClaimInfo, txHashs, pubkeys, mintValues, nativeValues)
	if err != nil {
		return err
	}
	return gc.SignAndSubmitTx(ext)
}
