package client_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/client"
	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/utils"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	LessPolka    = "1334v66HrtqQndbugYxX9m56V6222m97LbavB4KAMmqgjsas"
	From1        = "31d96Cq9idWQqPq3Ch5BFY84zrThVE3r98M7vG4xYaSWHwsX"
	From2        = "1TgYb5x8xjsZRyL5bwvxUoAWBn36psr4viSMHbRXA8bkB2h"
	Wen          = "1swvN162p1siDjm63UhhWoa59bpPZTSNKGVcbCwHUYkfRRW"
	Jun          = "33RQ73d9XfPTaE2SV7dzdhQQ17YaeMQ4yzhzAQhhFVenxMuJ"
	relay1       = "33xzQzUk75cAxt7i3hHb2XWwJNFqzcSULfoCRsAkiCG4jh5d"
	KeystorePath = "/Users/tpkeeper/gowork/stafi/rtoken-relay/keys"
)

var (
	tlog = client.NewLog()
)

const (
	stafiTypesFile  = "/Users/tpkeeper/gowork/stafi/rtoken-relay/network/stafi.json"
	polkaTypesFile  = "/Users/tpkeeper/gowork/stafi/rtoken-relay/network/polkadot.json"
	kusamaTypesFile = "/Users/tpkeeper/gowork/stafi/rtoken-relay/network/kusama.json"
)

func TestSarpcClient_GetChainEvents(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "wss://mainnet-rpc.stafi.io", "", client.AddressTypeAccountId, nil, tlog)
	// sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "wss://scan-rpc.stafi.io", "", client.AddressTypeAccountId, nil, tlog)
	//sc, err := client.NewGsrpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	// sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "ws://127.0.0.1:9944", "", client.AddressTypeAccountId, AliceKey, tlog)

	// sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "wss://stafi-seiya.stafi.io", "", client.AddressTypeAccountId, AliceKey, tlog)
	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://rpc.polkadot.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog,  )
	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-rpc.stafi.io", kusamaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog,  )
	if err != nil {
		t.Fatal(err)
	}

	receiver, err := sc.GetRFisReceiver()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receiver)
	// rate,err:=sc.GetEraRate(client.RKSM,4773)
	// if err!=nil{
	// 	t.Fatal(err)
	// }
	// t.Log(rate)

	// era,err:=sc.CurrentChainEra(client.RDOT)
	// if err!=nil{
	// 	t.Fatal(err)
	// }
	// t.Log(era)
	// rateLimit, err := sc.ActiveChangeRateLimit(client.RDOT)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(rateLimit)
	// 0xd6c1e8d44fe1b4b54773efa907ad6b404756bc761d625d4577264d663144ddd27025e075d5e2f6cde3cc051a31f0766000
	// 0xd6c1e8d44fe1b4b54773efa907ad6b404ddad12338d5de7866ed2d6abb9f61dc4a9e6f9b8d43f6ad008f8c291929dee201
	// for i:=14188841;i<1418884114188841;i++{

	// 	events, err := sc.GetEvents(uint64(i))
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	for _, e := range events {
	// 		t.Log(e.EventId, e.ModuleId)
	// 		if e.EventId == config.RFisWithdrawUnbondEventId {
	// 			w, err := client.EventWithdrawUnbondData(e)
	// 			if err != nil {
	// 				t.Fatal(err)
	// 			}
	// 			t.Log(w)
	// 		}
	// 	}
	// }

	// wg := sync.WaitGroup{}
	// for i := 1588890; i < 1588990; i++ {
	// 	wg.Add(1)
	// 	go func(height uint64) {
	// 		defer func() {
	// 			if err := recover(); err != nil {
	// 				panic(height)
	// 			}
	// 		}()
	// 		_, _ = sc.GetEvents(height)
	// 		fmt.Println(height)

	// 		wg.Done()
	// 	}(uint64(i))
	// }
	// wg.Wait()
}

func TestSarpcClient_GetChainEventNominationUpdated(t *testing.T) {

	// sc, err := client.NewGsrpcClient(ChainTypeStafi, "wss://stafi-seiya.stafi.io", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog,  )
	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "wss://mainnet-rpc.stafi.io", stafiTypesFile, client.AddressTypeAccountId, AliceKey, tlog)
	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot,"wss://polkadot-test-rpc.stafi.io", polkaTypesFile, AddressTypeAccountId, AliceKey, tlog,  )
	if err != nil {
		t.Fatal(err)
	}

	symbz, err := types.EncodeToBytes(client.RKSM)
	if err != nil {
		t.Fatal(err)
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := sc.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("bonded pools not extis")
	}

	t.Log(bondedPools)

	evts, err := sc.GetEvents(11561482)
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range evts {
		if evt.EventId != config.NominationUpdatedEventId {
			continue
		}
	}
}

func TestSarpcClient_GetExtrinsics1(t *testing.T) {
	//sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := client.NewGsrpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)

	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://rpc.polkadot.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	for i := 7411010; i >= 7311010; i-- {
		if i%10 == 0 {
			t.Log("i", i)
		}

		bh, err := sc.GetBlockHash(uint64(i))
		if err != nil {
			t.Fatal(err)
		}
		exts, err := sc.GetExtrinsics(bh)
		if err != nil {
			t.Fatal(err)
		}

		for _, ext := range exts {
			t.Log("exthash", ext.ExtrinsicHash)
			t.Log("moduleName", ext.CallModuleName)
			t.Log("methodName", ext.CallName)
			t.Log("address", ext.Address)
			t.Log(ext.Params)
			//for _, p := range ext.Params {
			//	if p.Name == config.ParamDest && p.Type == config.ParamDestType {
			//		//dest, ok := p.Value.(string)
			//		//fmt.Println("ok", ok)
			//		//fmt.Println(dest)
			//
			//		// polkadot-test
			//		dest, ok := p.Value.(map[string]interface{})
			//		t.Log("ok", ok)
			//		v, ok := dest["Id"]
			//		t.Log("ok1", ok)
			//		val, ok := v.(string)
			//		t.Log("ok2", ok)
			//		t.Log(val)
			//	}
			//
			//	t.Log("name", p.Name, "value", p.Value, "type", p.Type)
			//}
		}
	}
}

func TestSarpcClient_GetExtrinsics2(t *testing.T) {

	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	exts, err := sc.GetExtrinsics("0x6157da60a188b3f31d250afe5acb2da786417fec00973f1c7f863504fbca4642")
	if err != nil {
		t.Fatal(err)
	}

	for _, ext := range exts {
		t.Log("exthash", ext.ExtrinsicHash)
		t.Log("moduleName", ext.CallModuleName)
		t.Log("methodName", ext.CallName)
		t.Log("address", ext.Address)
		t.Log(ext.Params)
		for _, p := range ext.Params {
			if p.Name == config.ParamDest && p.Type == config.ParamDestType {
				//dest, ok := p.Value.(string)
				//fmt.Println("ok", ok)
				//fmt.Println(dest)

				// polkadot-test
				dest, ok := p.Value.(map[string]interface{})
				t.Log("ok", ok)
				v, ok := dest["Id"]
				t.Log("ok1", ok)
				val, ok := v.(string)
				t.Log("ok2", ok)
				t.Log(val)
			}

			t.Log("name", p.Name, "value", p.Value, "type", p.Type)
		}
	}
}

func TestBatchTransfer(t *testing.T) {

	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, client.AddressTypeAccountId, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")

	addrs := []types.Address{less, jun, wen, bao}

	amount, _ := utils.StringToBigint("3000" + "000000000000")
	value := types.NewUCompact(amount)
	calls := make([]types.Call, 0)

	ci, err := sc.FindCallIndex(config.MethodTransferKeepAlive)
	if err != nil {
		t.Fatal(err)
	}

	for _, addr := range addrs {
		call, err := types.NewCallWithCallIndex(
			ci,
			config.MethodTransferKeepAlive,
			addr,
			value,
		)
		if err != nil {
			t.Fatal(err)
		}
		calls = append(calls, call)
	}

	ext, err := sc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(ext)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConst(t *testing.T) {

	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	e, err := sc.ExistentialDeposit()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(e)
}

func TestPolkaQueryStorage(t *testing.T) {

	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	var index uint32
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageActiveEra, nil, nil, &index)
	if err != nil {
		panic(err)
	}

	if !exist {
		panic("not exist")
	}

	t.Log(index)
}

func TestStafiLocalQueryActiveEra(t *testing.T) {

	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, client.AddressTypeAccountId, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	var index uint32
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageActiveEra, nil, nil, &index)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	t.Log("activeEra", index)
}

func TestActive(t *testing.T) {

	//sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog,  )
	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	a := "0xac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c"
	mac, err := types.NewMultiAddressFromHexAccountID(a)
	if err != nil {
		t.Fatal(err)
	}
	ledger := new(client.StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, mac.AsID[:], nil, ledger)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	t.Log("ledger", ledger)

	t.Log(types.NewU128(big.Int(ledger.Active)))
}

func TestActive1(t *testing.T) {

	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://polkadot-test-rpc.stafi.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}

	a := "0x782a467d4ff23b660ca5f1ecf47f8537d4c35049541b6ebbf5381c00c4c158f7"
	b, _ := hexutil.Decode(a) // work
	//mac, err := types.NewAddressFromHexAccountID(a) // work
	//mac, err := types.NewMultiAddressFromHexAccountID(a) // work
	ledger := new(client.StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, b, nil, ledger)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	t.Log(types.NewU128(big.Int(ledger.Active)))
}

func TestHash(t *testing.T) {
	h, _ := types.NewHashFromHexString("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	a, _ := types.EncodeToBytes(h)

	fmt.Println(hexutil.Encode(a))
}

func TestPool(t *testing.T) {
	p := "0x782a467d4ff23b660ca5f1ecf47f8537d4c35049541b6ebbf5381c00c4c158f7"
	pool, _ := hexutil.Decode(p)
	pbz, _ := types.EncodeToBytes(pool)
	fmt.Println(pool)
	fmt.Println(pbz)

	//
	//gc, err := client.NewGsrpcClient("wss://stafi-seiya.stafi.io", AddressTypeAccountId, AliceKey, tlog,  )
	//assert.NoError(t, err)
	//
	//
	////poolBz, err := types.EncodeToBytes(pool)
	//symBz, err := types.EncodeToBytes(core.RKSM)
	//assert.NoError(t, err)
	//
	//var threshold uint16
	//exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, symBz, pbz, &threshold)
	//assert.NoError(t, err)
	//fmt.Println(exist)
	//fmt.Println()

}

func Test_KSM_GsrpcClient_Multisig(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", kusamaTypesFile, client.AddressTypeMultiAddress, krp, tlog)
	if err != nil {
		t.Fatal(err)
	}
	_ = sc

	//pool, err := hexutil.Decode("ac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c")
	//assert.NoError(t, err)

	// threshold := uint16(2)
	// //wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	// // jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	// relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	// others := []types.AccountID{
	// 	relay2.AsID,
	// }

	//for _, oth := range others {
	//	fmt.Println(hexutil.Encode(oth[:]))
	//}

	// bond, _ := utils.StringToBigint("1000000000000")
	// unbond := big.NewInt(0)

	// call, err := sc.BondOrUnbondCall(bond, unbond)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// h := utils.BlakeTwo256(call.Opaque)
	// t.Log("Extrinsic", call.Extrinsic)
	// t.Log("Opaque", hexutil.Encode(call.Opaque))
	// t.Log("callHash", hexutil.Encode(h[:]))

	// info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("info", info.Class, info.PartialFee, info.Weight)

	//optp := types.TimePoint{Height: 1964877, Index: 1}
	//tp := submodel.NewOptionTimePoint(optp)

	// tp := client.NewOptionTimePointEmpty()
	// ext, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, call.Opaque, false, info.Weight)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// err = sc.SignAndSubmitTx(ext)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
func Test_KSM_GsrpcClient_transfer(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", kusamaTypesFile, AddressTypeAccountId, krp, tlog,  )
	sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "ws://127.0.0.1:9944", kusamaTypesFile, client.AddressTypeMultiAddress, krp, tlog)
	if err != nil {
		t.Fatal(err)
	}

	//pool, err := hexutil.Decode("ac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c")
	//assert.NoError(t, err)

	// threshold := uint16(2)
	//wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	// jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	// others := []types.AccountID{
	// 	relay2.AsAccountID,
	// }

	//for _, oth := range others {
	//	fmt.Println(hexutil.Encode(oth[:]))
	//}

	// bond, _ := utils.StringToBigint("1000000000000")
	// unbond := big.NewInt(0)

	// call,err:=sc.TransferCall(relay2.AsAccountID[:],types.NewUCompact(big.NewInt(1000000)))
	// if err!=nil{
	// 	t.Fatal(err)
	// }
	ext, err := sc.NewUnsignedExtrinsic(config.MethodTransfer, relay2, types.NewUCompact(big.NewInt(1e10)))
	if err != nil {
		t.Fatal(err)
	}

	// call, err := sc.BondOrUnbondCall(bond, unbond)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// h := utils.BlakeTwo256(call.Opaque)
	// t.Log("Extrinsic", call.Extrinsic)
	// t.Log("Opaque", hexutil.Encode(call.Opaque))
	// t.Log("callHash", hexutil.Encode(h[:]))

	// info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("info", info.Class, info.PartialFee, info.Weight)

	//optp := types.TimePoint{Height: 1964877, Index: 1}
	//tp := submodel.NewOptionTimePoint(optp)

	// tp := submodel.NewOptionTimePointEmpty()
	// ext, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, call.Opaque, false, info.Weight)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	err = sc.SignAndSubmitTx(ext)
	if err != nil {
		t.Fatal(err)
	}
}
func TestSarpcClient_MintTxhashExist(t *testing.T) {
	//sc, err := client.NewGsrpcClient("wss://mainnet-rpc.stafi.io", stafiTypesFile, tlog)
	//sc, err := client.NewGsrpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, client.AddressTypeAccountId, AliceKey, tlog)

	// sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "wss://stafi-seiya.stafi.io", "", client.AddressTypeAccountId, AliceKey, tlog)
	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog,  )
	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot, "wss://kusama-rpc.stafi.io", kusamaTypesFile, client.AddressTypeMultiAddress, AliceKey, tlog,  )
	if err != nil {
		t.Fatal(err)
	}
	exist, err := sc.MintTxHashExist(types.NewBytes(hexutil.MustDecode("0x12345678")))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}

func TestSarpcClient_RTokenTotalIssuance(t *testing.T) {

	// sc, err := client.NewGsrpcClient(ChainTypeStafi, "wss://stafi-seiya.stafi.io", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog,  )
	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, "wss://mainnet-rpc.stafi.io", stafiTypesFile, client.AddressTypeAccountId, AliceKey, tlog)
	// sc, err := client.NewGsrpcClient(client.ChainTypePolkadot,"wss://polkadot-test-rpc.stafi.io", polkaTypesFile, AddressTypeAccountId, AliceKey, tlog,  )
	if err != nil {
		t.Fatal(err)

	}
	issuance, err := sc.RTokenTotalIssuance(client.RFIS)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(issuance)
}

func TestQueryStakePools(t *testing.T) {
	//rpc:="wss://mainnet-rpc.stafi.io"
	//rpc := "wss://stafi-seiya.stafi.io"
	rpc := "ws://127.0.0.1:9944"
	sc, err := client.NewGsrpcClient(client.ChainTypeStafi, rpc, stafiTypesFile, client.AddressTypeAccountId, AliceKey, tlog)
	if err != nil {
		t.Fatal(err)
	}
	pools, err := sc.StakePool(client.RFIS, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pools)
}
