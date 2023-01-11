package client

import (
	"errors"
	"fmt"

	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

func (sc *GsrpcClient) NewUnsignedExtrinsic(callMethod string, args ...interface{}) (interface{}, error) {
	sc.log.Debug("Submitting substrate call...", "callMethod", callMethod, "addressType", sc.addressType, "sender", sc.key.Address)

	ci, err := sc.FindCallIndex(callMethod)
	if err != nil {
		return nil, err
	}
	call, err := types.NewCallWithCallIndex(ci, callMethod, args...)
	if err != nil {
		return nil, err
	}

	if sc.addressType == AddressTypeAccountId {
		unsignedExt := types.NewExtrinsic(call)
		return &unsignedExt, nil
	} else if sc.addressType == AddressTypeMultiAddress {
		unsignedExt := types.NewExtrinsicMulti(call)
		return &unsignedExt, nil
	} else {
		return nil, errors.New("addressType not supported")
	}
}

func (sc *GsrpcClient) SignAndSubmitTx(ext interface{}) error {
	err := sc.signExtrinsic(ext)
	if err != nil {
		return err
	}
	sc.log.Trace("signExtrinsic ok")

	api, err := sc.FlashApi()
	if err != nil {
		return err
	}
	sc.log.Trace("flashApi ok")
	// Do the transfer and track the actual status
	sub, err := api.Author.SubmitAndWatch(ext)
	if err != nil {
		return err
	}
	sc.log.Trace("Extrinsic submission succeeded")
	defer sub.Unsubscribe()

	return sc.watchSubmission(sub)
}

func (sc *GsrpcClient) watchSubmission(sub *author.ExtrinsicStatusSubscription) error {
	for {
		select {
		case <-sc.stop:
			return ErrorTerminated
		case status := <-sub.Chan():
			switch {
			case status.IsInBlock:
				sc.log.Info("Extrinsic included in block", "block", status.AsInBlock.Hex())
				return nil
			case status.IsRetracted:
				return fmt.Errorf("extrinsic retracted: %s", status.AsRetracted.Hex())
			case status.IsDropped:
				return fmt.Errorf("extrinsic dropped from network")
			case status.IsInvalid:
				return fmt.Errorf("extrinsic invalid")
			}
		case err := <-sub.Err():
			sc.log.Trace("Extrinsic subscription error", "err", err)
			return err
		}
	}
}

func (sc *GsrpcClient) signExtrinsic(xt interface{}) error {
	rv, err := sc.GetLatestRuntimeVersion()
	if err != nil {
		return err
	}

	nonce, err := sc.GetLatestNonce()
	if err != nil {
		return err
	}

	o := types.SignatureOptions{
		BlockHash:          sc.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        sc.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	if ext, ok := xt.(*types.Extrinsic); ok {
		sc.log.Debug("signExtrinsic", "addressType", sc.addressType)
		err = ext.Sign(*sc.key, o)
		if err != nil {
			return err
		}
	} else if ext, ok := xt.(*types.ExtrinsicMulti); ok {
		sc.log.Debug("signExtrinsic", "addressType", sc.addressType)
		err = ext.Sign(*sc.key, o)
		if err != nil {
			return fmt.Errorf("sign err: %s", err)
		}
	} else {
		return errors.New("extrinsic cast error")
	}

	return nil
}

func (sc *GsrpcClient) SingleTransferTo(accountId []byte, value types.UCompact) error {
	var addr interface{}
	switch sc.addressType {
	case AddressTypeAccountId:
		addr = types.NewAddressFromAccountID(accountId)
	case AddressTypeMultiAddress:
		addr = types.NewMultiAddressFromAccountID(accountId)
	default:
		return fmt.Errorf("unsupported address type: %s", sc.addressType)
	}
	ext, err := sc.NewUnsignedExtrinsic(config.MethodTransferKeepAlive, addr, value)
	if err != nil {
		return err
	}
	return sc.SignAndSubmitTx(ext)
}
