package client

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	scale "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	scaleTypes "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/rpc"
	gsrpcConfig "github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/recws"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/stafidecoder"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/websocketpool"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

const (
	ChainTypeStafi    = "stafi"
	ChainTypePolkadot = "polkadot"

	AddressTypeAccountId    = "AccountId"
	AddressTypeMultiAddress = "MultiAddress"
)

var (
	ErrorTerminated           = errors.New("terminated")
	ErrorBondEqualToUnbond    = errors.New("ErrorBondEqualToUnbond")
	ErrorDiffSmallerThanLeast = errors.New("ErrorDiffSmallerThanLeast")
	ErrorValueNotExist        = errors.New("value not exist")
)

type GsrpcClient struct {
	endpoint    string
	addressType string
	rpcs        *gsrpc.RPCS
	key         *signature.KeyringPair
	genesisHash types.Hash

	wsPool    websocket_pool.Pool
	log       Logger
	chainType string
	typesPath string

	currentSpecVersion int

	stafiMetaDecoderMap map[int]*stafi_decoder.MetadataDecoder
	polkaMetaDecoderMap map[int]*scale.MetadataDecoder
	sync.RWMutex

	metaDataVersion int
}

func NewGsrpcClient(chainType, endpoint, typesPath, addressType string, key *signature.KeyringPair, log Logger) (*GsrpcClient, error) {
	log.Info("Connecting to substrate chain with gsrpc client", "endpoint", endpoint)

	if addressType != AddressTypeAccountId && addressType != AddressTypeMultiAddress {
		return nil, errors.New("addressType not supported")
	}

	rpcs, err := gsrpc.NewRPCS(endpoint)
	if err != nil {
		return nil, err
	}

	gsrpcConfig.SetSubscribeTimeout(2 * time.Minute)
	latestHash, err := rpcs.Chain.GetFinalizedHead()
	if err != nil {
		return nil, err
	}
	log.Info("NewGsrpcClient", "latestHash", latestHash.Hex())

	genesisHash, err := rpcs.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	sc := &GsrpcClient{
		endpoint:            endpoint,
		chainType:           chainType,
		addressType:         addressType,
		rpcs:                rpcs,
		key:                 key,
		genesisHash:         genesisHash,
		wsPool:              nil,
		log:                 log,
		typesPath:           typesPath,
		currentSpecVersion:  -1,
		stafiMetaDecoderMap: make(map[int]*stafi_decoder.MetadataDecoder),
		polkaMetaDecoderMap: make(map[int]*scale.MetadataDecoder),
	}

	err = sc.regCustomTypes()
	if err != nil {
		return nil, err
	}

	switch chainType {
	case ChainTypeStafi:
		_, err := sc.getStafiMetaDecoder(latestHash.Hex())
		if err != nil {
			return nil, err
		}
	case ChainTypePolkadot:
		_, err := sc.getPolkaMetaDecoder(latestHash.Hex())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("chain type not support: %s", chainType)
	}

	return sc, nil
}

func (s *GsrpcClient) getStafiMetaDecoder(blockHash string) (*stafi_decoder.MetadataDecoder, error) {
	v := &model.JsonRpcResult{}
	// runtime version
	if err := s.sendWsRequest(v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return nil, err
	}

	r := v.ToRuntimeVersion()
	if r == nil {
		return nil, fmt.Errorf("runtime version nil")
	}
	s.RLock()
	if decoder, exist := s.stafiMetaDecoderMap[r.SpecVersion]; exist {
		s.RUnlock()
		return decoder, nil
	}
	s.RUnlock()

	// check metadata need update, maybe  get ahead hash
	if err := s.sendWsRequest(v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
		return nil, err
	}
	metaRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	md := &stafi_decoder.MetadataDecoder{}
	md.Init(utiles.HexToBytes(metaRaw))
	if err := md.Process(); err != nil {
		return nil, err
	}
	s.Lock()
	s.stafiMetaDecoderMap[r.SpecVersion] = md
	s.Unlock()

	if r.SpecVersion > s.currentSpecVersion {
		s.currentSpecVersion = r.SpecVersion
		s.metaDataVersion = md.Metadata.MetadataVersion
	}
	return md, nil

}

func (s *GsrpcClient) getPolkaMetaDecoder(blockHash string) (*scale.MetadataDecoder, error) {
	v := &model.JsonRpcResult{}
	// runtime version
	if err := s.sendWsRequest(v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return nil, err
	}

	r := v.ToRuntimeVersion()
	if r == nil {
		return nil, fmt.Errorf("runtime version nil")
	}
	s.RLock()
	if decoder, exist := s.polkaMetaDecoderMap[r.SpecVersion]; exist {
		s.RUnlock()
		return decoder, nil
	}
	s.RUnlock()

	// check metadata need update, maybe  get ahead hash
	if err := s.sendWsRequest(v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
		return nil, err
	}
	metaRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	md := scale.MetadataDecoder{}
	md.Init(utiles.HexToBytes(metaRaw))
	if err := md.Process(); err != nil {
		return nil, err
	}
	s.Lock()
	s.polkaMetaDecoderMap[r.SpecVersion] = &md
	s.Unlock()

	if r.SpecVersion > s.currentSpecVersion {
		s.currentSpecVersion = r.SpecVersion
		s.metaDataVersion = md.Metadata.MetadataVersion
	}

	return &md, nil

}

func (sc *GsrpcClient) regCustomTypes() error {
	content := []byte(stafi_decoder.DefaultStafiCustumTypes)
	var err error
	if len(sc.typesPath) != 0 {
		content, err = os.ReadFile(sc.typesPath)
		if err != nil {
			return err
		}
	}

	switch sc.chainType {
	case ChainTypeStafi:
		stafi_decoder.RuntimeType{}.Reg()
		stafi_decoder.RegCustomTypes(stafi_decoder.LoadTypeRegistry(content))
	case ChainTypePolkadot:
		scaleTypes.RegCustomTypes(source.LoadTypeRegistry(content))
	default:
		return fmt.Errorf("chainType not supported")
	}
	return nil
}

func (sc *GsrpcClient) initial() (*websocket_pool.PoolConn, error) {
	var err error
	if sc.wsPool == nil {
		factory := func() (*websocket_pool.PoolConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 2 * time.Minute}
			SubscribeConn.Dial(sc.endpoint, nil)
			sc.log.Debug("conn factory create new conn", "endpoint", sc.endpoint)
			return websocket_pool.WrapConn(SubscribeConn), err
		}
		if sc.wsPool, err = websocket_pool.NewWsPool(1, 100, factory); err != nil {
			sc.log.Error("wbskt.NewChannelPool", "err", err)
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := sc.wsPool.Get()
	return conn, err
}
