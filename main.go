package gsrpc

import (
	"github.com/kaelnew/go-substrate-rpc-client/client"
	"github.com/kaelnew/go-substrate-rpc-client/rpc"
)

type SubstrateAPI struct {
	RPC    *rpc.RPC
	Client client.Client
}

func NewSubstrateAPI(url string) (*SubstrateAPI, error) {
	cl, err := client.Connect(url)
	if err != nil {
		return nil, err
	}

	newRPC, err := rpc.NewRPC(cl)
	if err != nil {
		return nil, err
	}

	return &SubstrateAPI{
		RPC:    newRPC,
		Client: cl,
	}, nil
}
