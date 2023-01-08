package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	scale "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	scaleTypes "github.com/itering/scale.go/types"
	scaleBytes "github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/rpc"
	gsrpcConfig "github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/recws"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/websocket_pool"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/submodel"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	commonTypes "github.com/stafiprotocol/go-substrate-rpc-client/types/common"
	"github.com/stafiprotocol/go-substrate-rpc-client/types/stafi"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

type GsrpcClient struct {
	endpoint    string
	addressType string
	rpcs        *gsrpc.RPCS
	key         *signature.KeyringPair
	genesisHash types.Hash
	stop        <-chan int

	wsPool    websocket_pool.Pool
	log       Logger
	chainType string
	typesPath string

	currentSpecVersion int

	stafiMetaDecoderMap map[int]*stafi.MetadataDecoder
	polkaMetaDecoderMap map[int]*scale.MetadataDecoder
	sync.RWMutex

	metaDataVersion int
}

func NewGsrpcClient(chainType, endpoint, typesPath, addressType string, key *signature.KeyringPair, log Logger, stop <-chan int) (*GsrpcClient, error) {
	log.Info("Connecting to substrate chain with sarpc", "endpoint", endpoint)

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
		stop:                stop,
		wsPool:              nil,
		log:                 log,
		typesPath:           typesPath,
		currentSpecVersion:  -1,
		stafiMetaDecoderMap: make(map[int]*stafi.MetadataDecoder),
		polkaMetaDecoderMap: make(map[int]*scale.MetadataDecoder),
	}

	sc.regCustomTypes()

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

func (s *GsrpcClient) getStafiMetaDecoder(blockHash string) (*stafi.MetadataDecoder, error) {
	v := &rpc.JsonRpcResult{}
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

	md := &stafi.MetadataDecoder{}
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
	v := &rpc.JsonRpcResult{}
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

func (sc *GsrpcClient) regCustomTypes() {
	content, err := os.ReadFile(sc.typesPath)
	if err != nil {
		panic(err)
	}

	switch sc.chainType {
	case ChainTypeStafi:
		stafi.RuntimeType{}.Reg()
		stafi.RegCustomTypes(source.LoadTypeRegistry(content))
	case ChainTypePolkadot:
		scaleTypes.RegCustomTypes(source.LoadTypeRegistry(content))
	default:
		panic("chainType not supported")
	}
}

func (sc *GsrpcClient) initial() (*websocket_pool.PoolConn, error) {
	var err error
	if sc.wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 2 * time.Minute}
			SubscribeConn.Dial(sc.endpoint, nil)
			sc.log.Debug("conn factory create new conn", "endpoint", sc.endpoint)
			return SubscribeConn, err
		}
		if sc.wsPool, err = websocket_pool.NewChannelPool(1, 25, factory); err != nil {
			sc.log.Error("wbskt.NewChannelPool", "err", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := sc.wsPool.Get()
	return conn, err
}

func (sc *GsrpcClient) sendWsRequest(v interface{}, action []byte) error {

	retry := 0
	for {
		if retry >= 100 {
			return fmt.Errorf("sendWsRequest reach retry limit")
		}

		if poolConn, err := sc.initial(); err != nil {
			return err
		} else {
			p := poolConn.Conn

			if err = p.WriteMessage(websocket.TextMessage, action); err != nil {
				poolConn.MarkUnusable()
				poolConn.Close()

				sc.log.Warn("websocket send error", "err", err)
				time.Sleep(time.Millisecond * 100)
				retry++
				continue
			}

			if err = p.ReadJSON(v); err != nil {
				poolConn.MarkUnusable()
				poolConn.Close()

				sc.log.Warn("websocket read error", "err", err)
				time.Sleep(time.Millisecond * 100)
				retry++
				continue
			}
			poolConn.Close()
			return nil
		}
	}
}

func (sc *GsrpcClient) GetBlock(blockHash string) (*rpc.Block, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(v, rpc.ChainGetBlock(wsId, blockHash)); err != nil {
		return nil, err
	}
	rpcBlock := v.ToBlock()
	return &rpcBlock.Block, nil
}

func (sc *GsrpcClient) GetExtrinsics(blockHash string) ([]*submodel.Transaction, error) {
	blk, err := sc.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	exts := make([]*submodel.Transaction, 0)
	switch sc.chainType {
	case ChainTypeStafi:
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		e := new(stafi.ExtrinsicDecoder)
		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata}
		for _, raw := range blk.Extrinsics {
			e.Init(stafi.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
			e.Process()
			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				ext := &submodel.Transaction{
					ExtrinsicHash:  e.ExtrinsicHash,
					CallModuleName: e.CallModule.Name,
					CallName:       e.Call.Name,
					Address:        e.Address,
					Params:         e.Params,
				}
				exts = append(exts, ext)
			}
		}
		return exts, nil
	case ChainTypePolkadot:
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		e := new(scale.ExtrinsicDecoder)
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata}
		for _, raw := range blk.Extrinsics {
			e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
			e.Process()
			decodeExtrinsic := e.Value.(map[string]interface{})
			var ce ChainExtrinsic
			eb, _ := json.Marshal(decodeExtrinsic)
			_ = json.Unmarshal(eb, &ce)
			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				params := make([]commonTypes.ExtrinsicParam, 0)
				for _, p := range e.Params {
					params = append(params, commonTypes.ExtrinsicParam{
						Name:  p.Name,
						Type:  p.Type,
						Value: p.Value,
					})
				}

				ext := &submodel.Transaction{
					ExtrinsicHash:  e.ExtrinsicHash,
					CallModuleName: ce.CallModule,
					CallName:       ce.CallModuleFunction,
					Address:        e.Address,
					Params:         params,
				}
				exts = append(exts, ext)
			}
		}
		return exts, nil
	default:
		return nil, errors.New("chainType not supported")
	}
}

func (sc *GsrpcClient) GetBlockHash(blockNum uint64) (string, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return "", fmt.Errorf("websocket get block hash error: %v", err)
	}

	blockHash, err := v.ToString()
	if err != nil {
		return "", fmt.Errorf("ChainGetBlockHash get error %v", err)
	}
	if blockHash == "" {
		return "", fmt.Errorf("ChainGetBlockHash error, blockHash empty")
	}

	return blockHash, nil
}

func (sc *GsrpcClient) GetChainEvents(blockHash string) ([]*submodel.ChainEvent, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	var events []*submodel.ChainEvent
	switch sc.chainType {
	case ChainTypeStafi:
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		e := stafi.EventsDecoder{}
		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata}
		e.Init(stafi.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
		e.Process()
		b, err := json.Marshal(e.Value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &events)
		if err != nil {
			return nil, err
		}

	case ChainTypePolkadot:
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata}

		e := scale.EventsDecoder{}
		e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)

		e.Process()

		b, err := json.Marshal(e.Value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &events)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("chainType not supported")
	}

	return events, nil
}

func (sc *GsrpcClient) GetEvents(blockNum uint64) ([]*submodel.ChainEvent, error) {
	blockHash, err := sc.GetBlockHash(blockNum)
	if err != nil {
		return nil, err
	}

	evts, err := sc.GetChainEvents(blockHash)
	if err != nil {
		return nil, err
	}

	return evts, nil
}

func (sc *GsrpcClient) GetPaymentQueryInfo(encodedExtrinsic string) (paymentInfo *rpc.PaymentQueryInfo, err error) {
	v := &rpc.JsonRpcResult{}
	if err = sc.sendWsRequest(v, rpc.SystemPaymentQueryInfo(wsId, encodedExtrinsic)); err != nil {
		return
	}

	paymentInfo = v.ToPaymentQueryInfo()
	if paymentInfo == nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	return
}