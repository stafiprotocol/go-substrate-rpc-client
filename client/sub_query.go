package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/websocket"
	scalecodec "github.com/itering/scale.go"
	scaleTypes "github.com/itering/scale.go/types"
	scaleBytes "github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	commonTypes "github.com/stafiprotocol/go-substrate-rpc-client/types/common"
	"github.com/stafiprotocol/go-substrate-rpc-client/types/stafi"
)

func (sc *GsrpcClient) FlashApi() (*gsrpc.RPCS, error) {
	_, err := sc.rpcs.Chain.GetBlockHashLatest()
	if err != nil {
		var rpcs *gsrpc.RPCS
		for i := 0; i < 3; i++ {
			rpcs, err = gsrpc.NewRPCS(sc.endpoint)
			if err == nil {
				break
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		if rpcs != nil {
			sc.rpcs = rpcs
		}
	}
	return sc.rpcs, nil
}

func (sc *GsrpcClient) Address() string {
	return sc.key.Address
}

func (sc *GsrpcClient) GetLatestBlockNumber() (uint64, error) {
	h, err := sc.GetHeaderLatest()
	if err != nil {
		return 0, err
	}

	return uint64(h.Number), nil
}

func (sc *GsrpcClient) GetFinalizedBlockNumber() (uint64, error) {
	hash, err := sc.GetFinalizedHead()
	if err != nil {
		return 0, err
	}

	header, err := sc.GetHeader(hash)
	if err != nil {
		return 0, err
	}

	return uint64(header.Number), nil
}

func (sc *GsrpcClient) GetHeaderLatest() (*types.Header, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return nil, err
	}
	return api.Chain.GetHeaderLatest()
}

func (sc *GsrpcClient) GetFinalizedHead() (types.Hash, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return types.NewHash([]byte{}), err
	}
	return api.Chain.GetFinalizedHead()
}

func (sc *GsrpcClient) GetHeader(blockHash types.Hash) (*types.Header, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return nil, err
	}
	return api.Chain.GetHeader(blockHash)
}

func (sc *GsrpcClient) GetBlockNumber(blockHash types.Hash) (uint64, error) {
	head, err := sc.GetHeader(blockHash)
	if err != nil {
		return 0, err
	}

	return uint64(head.Number), nil
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (sc *GsrpcClient) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	entry, err := sc.FindStorageEntryMetadata(prefix, method)
	if err != nil {
		return false, err
	}

	var key types.StorageKey
	keySeted := false
	if entry.IsNMap() {
		hashers, err := entry.Hashers()
		if err != nil {
			return false, err
		}

		if len(hashers) == 1 {
			key, err = types.CreateStorageKeyWithEntryMeta(uint8(sc.metaDataVersion), entry, prefix, method, arg1)
			if err != nil {
				return false, err
			}
			keySeted = true
		}
	}

	if !keySeted {
		key, err = types.CreateStorageKeyWithEntryMeta(uint8(sc.metaDataVersion), entry, prefix, method, arg1, arg2)
		if err != nil {
			return false, err
		}
	}

	api, err := sc.FlashApi()
	if err != nil {
		return false, err
	}

	ok, err := api.State.GetStorageLatest(key, result)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (sc *GsrpcClient) GetLatestRuntimeVersion() (*types.RuntimeVersion, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return nil, err
	}
	rv, err := api.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (sc *GsrpcClient) GetLatestNonce() (types.U32, error) {
	ac, err := sc.GetAccountInfo()
	if err != nil {
		return 0, err
	}

	return ac.Nonce, nil
}

func (sc *GsrpcClient) GetAccountInfo() (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := sc.QueryStorage("System", "Account", sc.key.PublicKey, nil, &ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("account not exist")
	}

	return ac, nil
}

func (sc *GsrpcClient) PublicKey() []byte {
	return sc.key.PublicKey
}

func (sc *GsrpcClient) StakingLedger(ac types.AccountID) (*StakingLedger, error) {
	s := new(StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, ac[:], nil, s)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get active for account: %s", hexutil.Encode(ac[:]))
	}

	return s, nil
}

func (sc *GsrpcClient) FreeBalance(who []byte) (types.U128, error) {
	if sc.addressType == AddressTypeMultiAddress {
		info, err := sc.NewVersionAccountInfo(who)
		if err != nil {
			return types.U128{}, err
		}
		return info.Data.Free, nil
	}

	info, err := sc.AccountInfo(who)
	if err != nil {
		return types.U128{}, err
	}

	return info.Data.Free, nil
}

func (sc *GsrpcClient) AccountInfo(who []byte) (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := sc.QueryStorage(config.SystemModuleId, config.StorageAccount, who, nil, ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get accountInfo for account: %s", hexutil.Encode(who))
	}

	return ac, nil
}

func (sc *GsrpcClient) NewVersionAccountInfo(who []byte) (*AccountInfo, error) {
	ac := new(AccountInfo)
	exist, err := sc.QueryStorage(config.SystemModuleId, config.StorageAccount, who, nil, ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get accountInfo for account: %s", hexutil.Encode(who))
	}
	return ac, nil
}

func (sc *GsrpcClient) ExistentialDeposit() (types.U128, error) {
	_, err := sc.FlashApi()
	if err != nil {
		return types.U128{}, err
	}
	var e types.U128
	err = sc.GetConst(config.BalancesModuleId, config.ConstExistentialDeposit, &e)
	if err != nil {
		return types.U128{}, err
	}
	return e, nil
}

func (sc *GsrpcClient) GetConst(prefix, name string, res interface{}) error {

	switch sc.chainType {
	case ChainTypeStafi:
		return sc.rpcs.State.GetConst(prefix, name, &res)
	case ChainTypePolkadot:
		blockHash, err := sc.GetFinalizedHead()
		if err != nil {
			return err
		}
		md, err := sc.getPolkaMetaDecoder(blockHash.Hex())
		if err != nil {
			return err
		}

		for _, mod := range md.Metadata.Metadata.Modules {
			if string(mod.Prefix) == prefix {
				for _, cons := range mod.Constants {
					if cons.Name == name {

						return types.DecodeFromHexString(cons.ConstantsValue, res)
					}
				}
			}
		}
		return fmt.Errorf("could not find constant %s.%s", prefix, name)
	default:
		return errors.New("GetConst chainType not supported")
	}
}

func (sc *GsrpcClient) FindStorageEntryMetadata(module string, fn string) (types.StorageEntryMetadata, error) {
	switch sc.chainType {
	case ChainTypeStafi:
		meta, err := sc.rpcs.State.GetMetadataLatest()
		if err != nil {
			return nil, err
		}

		return meta.FindStorageEntryMetadata(module, fn)
	case ChainTypePolkadot:
		blockHash, err := sc.GetFinalizedHead()
		if err != nil {
			return nil, err
		}
		md, err := sc.getPolkaMetaDecoder(blockHash.Hex())
		if err != nil {
			return nil, err
		}

		for _, mod := range md.Metadata.Metadata.Modules {
			if string(mod.Prefix) != module {
				continue
			}
			for _, s := range mod.Storage {
				if string(s.Name) != fn {
					continue
				}

				sfm := types.StorageFunctionMetadataV13{
					Name: types.Text(s.Name),
				}

				if s.Type.PlainType != nil {
					sfm.Type = types.StorageFunctionTypeV13{
						IsType: true,
						AsType: types.Type(*s.Type.PlainType),
					}
				}

				if s.Type.DoubleMapType != nil {
					dmt := types.DoubleMapTypeV10{
						Key1:       types.Type(s.Type.DoubleMapType.Key),
						Key2:       types.Type(s.Type.DoubleMapType.Key2),
						Value:      types.Type(s.Type.DoubleMapType.Value),
						Hasher:     TransformHasher(s.Type.DoubleMapType.Hasher),
						Key2Hasher: TransformHasher(s.Type.DoubleMapType.Key2Hasher),
					}

					sfm.Type = types.StorageFunctionTypeV13{
						IsDoubleMap: true,
						AsDoubleMap: dmt,
					}
				}

				if s.Type.MapType != nil {
					mt := types.MapTypeV10{
						Key:    types.Type(s.Type.MapType.Key),
						Value:  types.Type(s.Type.MapType.Value),
						Linked: s.Type.MapType.IsLinked,
						Hasher: TransformHasher(s.Type.MapType.Hasher),
					}

					sfm.Type = types.StorageFunctionTypeV13{
						IsMap: true,
						AsMap: mt,
					}
				}

				if s.Type.NMapType != nil {
					keys := make([]types.Type, 0)
					for _, key := range s.Type.NMapType.KeyVec {
						keys = append(keys, types.Type(key))
					}

					hashers := make([]types.StorageHasherV10, 0)
					for _, hasher := range s.Type.NMapType.Hashers {
						hashers = append(hashers, TransformHasher(hasher))
					}

					nmt := types.NMapTypeV13{
						Keys:    keys,
						Hashers: hashers,
						Value:   types.Type(s.Type.NMapType.Value),
					}

					sfm.Type = types.StorageFunctionTypeV13{
						IsNMap: true,
						AsNMap: nmt,
					}
				}

				return sfm, nil
			}
			return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
		}
		return nil, fmt.Errorf("module %v not found in metadata", module)
	default:
		return nil, errors.New("chainType not supported")
	}
}

func (sc *GsrpcClient) FindCallIndex(call string) (types.CallIndex, error) {
	switch sc.chainType {
	case ChainTypeStafi:
		meta, err := sc.rpcs.State.GetMetadataLatest()
		if err != nil {
			return types.CallIndex{}, err
		}

		return meta.FindCallIndex(call)
	case ChainTypePolkadot:
		blockHash, err := sc.GetFinalizedHead()
		if err != nil {
			return types.CallIndex{}, err
		}

		md, err := sc.getPolkaMetaDecoder(blockHash.Hex())
		if err != nil {
			return types.CallIndex{}, err
		}
		s := strings.Split(call, ".")

		for _, mod := range md.Metadata.Metadata.Modules {
			if string(mod.Name) != s[0] {
				continue
			}
			for ci, f := range mod.Calls {
				if string(f.Name) == s[1] {
					return types.CallIndex{SectionIndex: uint8(mod.Index), MethodIndex: uint8(ci)}, nil
				}
			}
			return types.CallIndex{}, fmt.Errorf("method %v not found within module %v for call %v", s[1], mod.Name, call)
		}
		return types.CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)

	default:
		return types.CallIndex{}, errors.New("FindCallIndex chainType not supported")
	}
}

func TransformHasher(Hasher string) types.StorageHasherV10 {
	if Hasher == "Blake2_128" {
		return types.StorageHasherV10{IsBlake2_128: true}
	}

	if Hasher == "Blake2_256" {
		return types.StorageHasherV10{IsBlake2_256: true}
	}

	if Hasher == "Blake2_128Concat" {
		return types.StorageHasherV10{IsBlake2_128Concat: true}
	}

	if Hasher == "Twox128" {
		return types.StorageHasherV10{IsTwox128: true}
	}

	if Hasher == "Twox256" {
		return types.StorageHasherV10{IsTwox256: true}
	}

	if Hasher == "Twox64Concat" {
		return types.StorageHasherV10{IsTwox64Concat: true}
	}

	return types.StorageHasherV10{IsIdentity: true}
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

			if err = poolConn.Conn.WriteMessage(websocket.TextMessage, action); err != nil {
				poolConn.MarkUnusable()
				sc.wsPool.Put(poolConn)

				sc.log.Debug("websocket send error", "err", err)
				time.Sleep(time.Millisecond * 100)
				retry++
				continue
			}

			if err = poolConn.Conn.ReadJSON(v); err != nil {
				poolConn.MarkUnusable()
				sc.wsPool.Put(poolConn)

				sc.log.Debug("websocket read error", "err", err)
				time.Sleep(time.Millisecond * 100)
				retry++
				continue
			}

			sc.wsPool.Put(poolConn)
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

func (sc *GsrpcClient) GetExtrinsics(blockHash string) ([]*Transaction, error) {
	blk, err := sc.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	exts := make([]*Transaction, 0)
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
				ext := &Transaction{
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
		e := new(scalecodec.ExtrinsicDecoder)
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

				ext := &Transaction{
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

func (sc *GsrpcClient) GetChainEvents(blockHash string) ([]*ChainEvent, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	var events []*ChainEvent
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

		e := scalecodec.EventsDecoder{}
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

func (sc *GsrpcClient) GetEvents(blockNum uint64) ([]*ChainEvent, error) {
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

func (sc *GsrpcClient) GetBlockTimestampAndExtrinsics(height uint64) (uint64, map[int]*Transaction, error) {
	blockHash, err := sc.GetBlockHash(height)
	if err != nil {
		return 0, nil, err
	}

	blk, err := sc.GetBlock(blockHash)
	if err != nil {
		return 0, nil, err
	}
	if len(blk.Extrinsics) == 0 {
		return 0, nil, fmt.Errorf("no set time extrinsic in block: %d", height)
	}

	exts := make(map[int]*Transaction)
	switch sc.chainType {
	case ChainTypeStafi:
		e := new(stafi.ExtrinsicDecoder)
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return 0, nil, err
		}

		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata, Spec: md.Spec}

		raw := blk.Extrinsics[0]
		e.Init(stafi.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
		e.Process()
		if len(e.Params) == 0 {
			return 0, nil, fmt.Errorf("no params")
		}
		stamp, ok := e.Params[0].Value.(int)
		if !ok {
			return 0, nil, fmt.Errorf("interface not ok: %s", e.Params[0].Value)
		}

		for index, raw := range blk.Extrinsics {
			e.Init(stafi.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
			e.Process()
			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				ext := &Transaction{
					ExtrinsicHash:  utiles.AddHex(e.ExtrinsicHash),
					CallModuleName: e.CallModule.Name,
					CallName:       e.Call.Name,
					Address:        e.Address,
					Params:         e.Params,
				}
				exts[index] = ext
			}
		}
		return uint64(stamp), exts, nil
	case ChainTypePolkadot:
		e := new(scalecodec.ExtrinsicDecoder)
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return 0, nil, err
		}
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata, Spec: md.Spec}

		raw := blk.Extrinsics[0]
		e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
		e.Process()

		if len(e.Params) == 0 {
			return 0, nil, fmt.Errorf("no params")
		}
		stamp, ok := e.Params[0].Value.(int)
		if !ok {
			return 0, nil, fmt.Errorf("interface not ok: %s", e.Params[0].Value)
		}

		for index, raw := range blk.Extrinsics {
			e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
			e.Process()

			call, exist := e.Metadata.CallIndex[e.CallIndex]
			if !exist {
				return 0, nil, fmt.Errorf("callIndex: %s not exist metaData", e.CallIndex)
			}

			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				params := make([]commonTypes.ExtrinsicParam, 0)
				for _, p := range e.Params {
					params = append(params, commonTypes.ExtrinsicParam{
						Name:  p.Name,
						Type:  p.Type,
						Value: p.Value,
					})
				}
				ext := &Transaction{
					ExtrinsicHash:  utiles.AddHex(e.ExtrinsicHash),
					CallModuleName: call.Module.Name,
					CallName:       call.Call.Name,
					Address:        e.Address,
					Params:         params,
				}
				exts[index] = ext
			}
		}

		return uint64(stamp), exts, nil
	default:
		return 0, nil, errors.New("chainType not supported")
	}
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

func (c *GsrpcClient) CurrentEra() (uint32, error) {
	var index uint32
	exist, err := c.QueryStorage(config.StakingModuleId, config.StorageActiveEra, nil, nil, &index)
	if err != nil {
		return 0, err
	}

	if !exist {
		return 0, fmt.Errorf("unable to get activeEraInfo")
	}

	return index, nil
}
