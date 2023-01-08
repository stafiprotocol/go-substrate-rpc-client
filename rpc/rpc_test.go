package rpc_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

func Example_simpleConnect() {
	// The following example shows how to instantiate a Substrate API and use it to connect to a node

	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	chain, err := rpcs.System.Chain()
	if err != nil {
		panic(err)
	}
	nodeName, err := rpcs.System.Name()
	if err != nil {
		panic(err)
	}
	nodeVersion, err := rpcs.System.Version()
	if err != nil {
		panic(err)
	}

	fmt.Printf("You are connected to chain %v using %v v%v\n", chain, nodeName, nodeVersion)

	// Output: You are connected to chain Development using Substrate Node v2.0.0-a200cdb9-x86_64-linux-gnu
}

func Example_listenToNewBlocks() {
	// This example shows how to subscribe to new blocks.
	//
	// It displays the block number every time a new block is seen by the node you are connected to.
	//
	// NOTE: The example runs until 10 blocks are received or until you stop it with CTRL+C

	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	sub, err := rpcs.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	count := 0

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		count++

		if count == 10 {
			sub.Unsubscribe()
			break
		}
	}
}

func Example_listenToBalanceChange() {
	// This example shows how to instantiate a Substrate API and use it to connect to a node and retrieve balance
	// updates
	//
	// NOTE: The example runs until you stop it with CTRL+C

	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := rpcs.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Known account we want to use (available on dev chain, with funds)
	alice, err := types.HexDecodeString("0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d")
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "Balances", "FreeBalance", alice, nil)
	if err != nil {
		panic(err)
	}

	// Retrieve the initial balance
	var previous types.U128
	ok, err := rpcs.State.GetStorageLatest(key, &previous)
	if err != nil || !ok {
		panic(err)
	}

	fmt.Printf("%#x has a balance of %v\n", alice, previous)
	fmt.Printf("You may leave this example running and transfer any value to %#x\n", alice)

	// Here we subscribe to any balance changes
	sub, err := rpcs.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		// inner loop for the changes within one of those notifications
		for _, chng := range (<-sub.Chan()).Changes {
			var current types.U128
			if err = types.DecodeFromBytes(chng.StorageData, &current); err != nil {
				panic(err)
			}

			// Calculate the delta
			var change = types.U128{Int: big.NewInt(0).Sub(current.Int, previous.Int)}

			// Only display positive value changes (Since we are pulling `previous` above already,
			// the initial balance change will also be zero)
			if change.Cmp(big.NewInt(0)) != 0 {
				previous = current
				fmt.Printf("New balance change of: %v\n", change)
				return
			}
		}
	}
}

func Example_unsubscribeFromListeningToUpdates() {
	// This example shows how to subscribe to and later unsubscribe from listening to block updates.
	//
	// In this example we're calling the built-in unsubscribe() function after a timeOut of 20s to cleanup and
	// unsubscribe from listening to updates.

	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	sub, err := rpcs.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	timeout := time.After(20 * time.Second)

	for {
		select {
		case head := <-sub.Chan():
			fmt.Printf("Chain is at block: #%v\n", head.Number)
		case <-timeout:
			sub.Unsubscribe()
			fmt.Println("Unsubscribed")
			return
		}
	}
}

func Example_makeASimpleTransfer() {
	// This sample shows how to create a transaction to make a transfer from one an account to another.

	// Instantiate the API
	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := rpcs.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(12345))
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := rpcs.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := rpcs.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err := rpcs.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         types.NewUCompactFromUInt(0),
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	// Send the extrinsic
	hash, err := rpcs.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transfer sent with hash %#x\n", hash)
}

func Example_displaySystemEvents() {
	// Query the system events and extract information from them. This example runs until exited via Ctrl-C

	// Create our API with a default connection to the local node
	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	meta, err := rpcs.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}

	sub, err := rpcs.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		set := <-sub.Chan()
		// inner loop for the changes within one of those notifications
		for _, chng := range set.Changes {
			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
				// skip, we are only interested in events with content
				continue
			}

			// Decode the event records
			events := types.EventRecords{}
			err = types.EventRecordsRaw(chng.StorageData).DecodeEventRecords(meta, &events)
			if err != nil {
				panic(err)
			}

			// Show what we are busy with
			for _, e := range events.Balances_Endowed {
				fmt.Printf("\tBalances:Endowed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Balances_DustLost {
				fmt.Printf("\tBalances:DustLost:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Balances_Transfer {
				fmt.Printf("\tBalances:Transfer:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v, %v\n", e.From, e.To, e.Value)
			}
			for _, e := range events.Balances_BalanceSet {
				fmt.Printf("\tBalances:BalanceSet:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v, %v\n", e.Who, e.Free, e.Reserved)
			}
			for _, e := range events.Balances_Deposit {
				fmt.Printf("\tBalances:Deposit:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Grandpa_NewAuthorities {
				fmt.Printf("\tGrandpa:NewAuthorities:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.NewAuthorities)
			}
			for _, e := range events.Grandpa_Paused {
				fmt.Printf("\tGrandpa:Paused:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.Grandpa_Resumed {
				fmt.Printf("\tGrandpa:Resumed:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.ImOnline_HeartbeatReceived {
				fmt.Printf("\tImOnline:HeartbeatReceived:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x\n", e.AuthorityID)
			}
			for _, e := range events.ImOnline_AllGood {
				fmt.Printf("\tImOnline:AllGood:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.ImOnline_SomeOffline {
				fmt.Printf("\tImOnline:SomeOffline:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.IdentificationTuples)
			}
			for _, e := range events.Indices_IndexAssigned {
				fmt.Printf("\tIndices:IndexAssigned:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.AccountIndex)
			}
			for _, e := range events.Indices_IndexFreed {
				fmt.Printf("\tIndices:IndexFreed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.AccountIndex)
			}
			for _, e := range events.Offences_Offence {
				fmt.Printf("\tOffences:Offence:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v%v\n", e.Kind, e.OpaqueTimeSlot)
			}
			for _, e := range events.Session_NewSession {
				fmt.Printf("\tSession:NewSession:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.SessionIndex)
			}
			for _, e := range events.Staking_Reward {
				fmt.Printf("\tStaking:Reward:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.Amount)
			}
			for _, e := range events.Staking_Slash {
				fmt.Printf("\tStaking:Slash:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.Balance)
			}
			for _, e := range events.Staking_OldSlashingReportDiscarded {
				fmt.Printf("\tStaking:OldSlashingReportDiscarded:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.SessionIndex)
			}
			for _, e := range events.System_ExtrinsicSuccess {
				fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.System_ExtrinsicFailed {
				fmt.Printf("\tSystem:ErtrinsicFailed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.DispatchError)
			}
			for _, e := range events.System_CodeUpdated {
				fmt.Printf("\tSystem:CodeUpdated:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.System_NewAccount {
				fmt.Printf("\tSystem:NewAccount:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x\n", e.Who)
			}
			for _, e := range events.System_KilledAccount {
				fmt.Printf("\tSystem:KilledAccount:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#X\n", e.Who)
			}
		}
	}
}

func Example_transactionWithEvents() {
	// Display the events that occur during a transfer by sending a value to bob

	// Instantiate the API
	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := rpcs.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	amount := types.NewUCompactFromUInt(12345)

	c, err := types.NewCall(meta, "Balances.transfer", bob, amount)
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	if err != nil {
		panic(err)
	}

	genesisHash, err := rpcs.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := rpcs.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	// Get the nonce for Alice
	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err := rpcs.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         types.NewUCompactFromUInt(0),
	}

	fmt.Printf("Sending %v from %#x to %#x with nonce %v", amount, signature.TestKeyringPairAlice.PublicKey, bob.AsAccountID, nonce)

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	// Do the transfer and track the actual status
	sub, err := rpcs.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)
			return
		}
	}
}

func TestV13(t *testing.T) {
	url := "wss://kusama-rpc.polkadot.io"
	rpcs, err := rpc.NewRPCS(url)
	if err != nil {
		panic(err)
	}

	genesisHash, err := rpcs.Chain.GetBlockHash(0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(genesisHash.Hex())
}

type StakingLedger struct {
	Stash          types.AccountID
	Total          types.UCompact
	Active         types.UCompact
	Unlocking      []UnlockChunk
	ClaimedRewards []uint32
}

type UnlockChunk struct {
	Value types.UCompact
	Era   types.UCompact
}

func QueryStakingLeder(endpoint string, ac types.AccountID) (*StakingLedger, bool, error) {
	rpcs, err := rpc.NewRPCS(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := rpcs.State.GetMetadataLatest()
	if err != nil {
		return nil, false, err
	}

	key, err := types.CreateStorageKey(meta, "Staking", "Ledger", ac[:], nil)
	if err != nil {
		return nil, false, err
	}

	ledger := new(StakingLedger)
	ok, err := rpcs.State.GetStorageLatest(key, ledger)
	if err != nil {
		return nil, false, err
	}

	return ledger, ok, nil
}
