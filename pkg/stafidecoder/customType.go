package stafi_decoder

import (
	"regexp"
	"strings"

	"github.com/itering/scale.go/utiles"
)

func newStruct(names, typeString []string) *TypeMapping {
	if len(names) != len(typeString) {
		panic("init newStruct names and typeString length not equal")
	}
	if len(names) == 0 {
		return nil
	}
	return &TypeMapping{Names: names, Types: typeString}
}

func RegCustomTypes(registry map[string]TypeStruct) {
	for key, typeStruct := range registry {
		key = strings.ToLower(key)
		switch typeStruct.Type {
		case "string":
			typeString := typeStruct.TypeString
			instant := TypeRegistry[strings.ToLower(typeString)]
			if instant != nil {
				regCustomKey(key, instant)
				continue
			}

			// Explained
			if explainedType, ok := registry[typeString]; ok {
				if explainedType.Type == "string" {
					instant := TypeRegistry[strings.ToLower(explainedType.TypeString)]
					if instant != nil {
						regCustomKey(key, instant)
						continue
					}
				} else {
					RegCustomTypes(map[string]TypeStruct{key: registry[typeString]})
				}
			}

			// sub type vec|option
			if typeString[len(typeString)-1:] == ">" {
				reg := regexp.MustCompile("^([^<]*)<(.+)>$")
				typeParts := reg.FindStringSubmatch(typeString)
				if len(typeParts) > 2 {
					if strings.ToLower(typeParts[1]) == "vec" {
						v := Vec{}
						v.SubType = typeParts[2]
						regCustomKey(key, &v)
						continue
					} else if strings.ToLower(typeParts[1]) == "option" {
						v := Option{}
						v.SubType = typeParts[2]
						regCustomKey(key, &v)
						continue
					} else if strings.ToLower(typeParts[1]) == "compact" {
						v := Compact{}
						v.SubType = typeParts[2]
						regCustomKey(key, &v)
						continue
					}
				}

			}

			// Tuple
			if typeString != "()" && string(typeString[0]) == "(" && typeString[len(typeString)-1:] == ")" {
				s := Struct{}
				s.TypeString = typeString
				s.buildStruct()
				regCustomKey(key, &s)
				continue
			}

			// Array
			if typeString != "[]" && string(typeString[0]) == "[" && string(typeString[len(typeString)-1:]) == "]" {
				if typePart := strings.Split(string(typeString[1:len(typeString)-1]), ";"); len(typePart) == 2 {
					fixed := FixedLengthArray{
						FixedLength: utiles.StringToInt(strings.TrimSpace(typePart[1])),
						SubType:     strings.TrimSpace(typePart[0]),
					}
					regCustomKey(key, &fixed)
					continue
				}
			}
		case "struct":
			var names, typeStrings []string
			for _, v := range typeStruct.TypeMapping {
				names = append(names, v[0])
				typeStrings = append(typeStrings, v[1])
			}
			s := Struct{}
			s.TypeMapping = newStruct(names, typeStrings)

			regCustomKey(key, &s)
		case "enum":
			var names, typeStrings []string
			for _, v := range typeStruct.TypeMapping {
				names = append(names, v[0])
				typeStrings = append(typeStrings, v[1])
			}
			e := Enum{ValueList: typeStruct.ValueList}
			e.TypeMapping = newStruct(names, typeStrings)
			regCustomKey(key, &e)
		case "set":
			regCustomKey(key, &Set{ValueList: typeStruct.ValueList, BitLength: typeStruct.BitLength})
		}
	}
}

func regCustomKey(key string, rt interface{}) {
	slice := strings.Split(key, "#")
	if len(slice) == 2 { // for Special
		special := Special{Registry: rt, Version: []int{0, 99999999}}
		if version := strings.Split(slice[1], "-"); len(version) == 2 {
			special.Version[0] = utiles.StringToInt(version[0])
			if version[1] != "?" {
				special.Version[1] = utiles.StringToInt(version[1])
			}
		}
		if specialRegistry == nil {
			specialRegistry = make(map[string][]Special)
		}
		if instant, ok := specialRegistry[slice[0]]; ok {
			instant = append(instant, special)
			specialRegistry[slice[0]] = instant
		} else {
			specialRegistry[slice[0]] = []Special{special}
		}

	} else {
		TypeRegistry[key] = rt
	}

}

var DefaultStafiCustumTypes = `{
	"Weight": "u64",
	"DispatchResult": {
	  "type": "enum",
	  "type_mapping": [
		[
		  "Ok",
		  "Null"
		],
		[
		  "Error",
		  "DispatchError"
		]
	  ]
	},
	"Address": "GenericAddress",
	"LookupSource": "GenericAddress",
	"Keys": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "grandpa",
		  "AccountId"
		],
		[
		  "babe",
		  "AccountId"
		],
		[
		  "im_online",
		  "AccountId"
		],
		[
		  "authority_discovery",
		  "AccountId"
		]
	  ]
	},
	"ChainId": "u8",
	"RateType": "u64",
	"ResourceId": "[u8; 32]",
	"DepositNonce": "u64",
	"U256": "H256",
	"XSymbol": {
	  "type": "enum",
	  "value_list": [
		"WRA"
	  ]
	},
	"AccountXData": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "free",
		  "u128"
		]
	  ]
	},
	"RSymbol": {
	  "type": "enum",
	  "value_list": [
		"RFIS",
		"RDOT",
		"RKSM",
		"RATOM",
		"RSOL",
		"RMATIC",
		"RBNB",
		"RETH"
	  ]
	},
	"AccountRData": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "free",
		  "u128"
		]
	  ]
	},
	"ProposalStatus": {
	  "type": "enum",
	  "value_list": [
		"Active",
		"Passed",
		"Expired",
		"Executed"
	  ]
	},
	"ProposalVotes": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "voted",
		  "Vec<AccountId>"
		],
		[
		  "status",
		  "ProposalStatus"
		],
		[
		  "expiry",
		  "BlockNumber"
		]
	  ]
	},
	"BondRecord": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "bonder",
		  "AccountId"
		],
		[
		  "symbol",
		  "RSymbol"
		],
		[
		  "pubkey",
		  "Vec<u8>"
		],
		[
		  "pool",
		  "Vec<u8>"
		],
		[
		  "blockhash",
		  "Vec<u8>"
		],
		[
		  "txhash",
		  "Vec<u8>"
		],
		[
		  "amount",
		  "u128"
		]
	  ]
	},
	"BondReason": {
	  "type": "enum",
	  "value_list": [
		"Pass",
		"BlockhashUnmatch",
		"TxhashUnmatch",
		"PubkeyUnmatch",
		"PoolUnmatch",
		"AmountUnmatch"
	  ]
	},
	"BondState": {
	  "type": "enum",
	  "value_list": [
		"Dealing",
		"Fail",
		"Success"
	  ]
	},
	"RproposalStatus": {
	  "type": "enum",
	  "value_list": [
		"Initiated",
		"Approved",
		"Rejected",
		"Expired"
	  ]
	},
	"RproposalVotes": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "votes_for",
		  "Vec<AccountId>"
		],
		[
		  "votes_against",
		  "Vec<AccountId>"
		],
		[
		  "status",
		  "RproposalStatus"
		],
		[
		  "expiry",
		  "BlockNumber"
		]
	  ]
	},
	"LinkChunk": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "bond",
		  "u128"
		],
		[
		  "unbond",
		  "u128"
		],
		[
		  "active",
		  "u128"
		]
	  ]
	},
	"Unbonding": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "who",
		  "AccountId"
		],
		[
		  "value",
		  "u128"
		],
		[
		  "recipient",
		  "Vec<u8>"
		]
	  ]
	},
	"OriginalTxType": {
	  "type": "enum",
	  "value_list": [
		"Transfer",
		"Bond",
		"Unbond",
		"WithdrawUnbond",
		"ClaimRewards"
	  ]
	},
	"PoolBondState": {
	  "type": "enum",
	  "value_list": [
		"EraUpdated",
		"BondReported",
		"ActiveReported",
		"WithdrawSkipped",
		"WithdrawReported",
		"TransferReported"
	  ]
	},
	"BondSnapshot": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "symbol",
		  "RSymbol"
		],
		[
		  "era",
		  "u32"
		],
		[
		  "pool",
		  "Vec<u8>"
		],
		[
		  "bond",
		  "u128"
		],
		[
		  "unbond",
		  "u128"
		],
		[
		  "active",
		  "u128"
		],
		[
		  "last_voter",
		  "AccountId"
		],
		[
		  "bond_state",
		  "PoolBondState"
		]
	  ]
	},
	"UserUnlockChunk": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "pool",
		  "Vec<u8>"
		],
		[
		  "unlock_era",
		  "u32"
		],
		[
		  "value",
		  "u128"
		],
		[
		  "recipient",
		  "Vec<u8>"
		]
	  ]
	},
	"SigVerifyResult": {
	  "type": "enum",
	  "value_list": [
		"InvalidPubkey",
		"Fail",
		"Pass"
	  ]
	},
	"BondAction": {
	  "type": "enum",
	  "value_list": [
		"BondOnly",
		"UnbondOnly",
		"BothBondUnbond",
		"EitherBondUnbond",
		"InterDeduct"
	  ]
	},
	"SwapTransactionInfo": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "account",
		  "AccountId"
		],
		[
		  "receiver",
		  "Vec<u8>"
		],
		[
		  "value",
		  "u128"
		],
		[
		  "is_deal",
		  "bool"
		]
	  ]
	},
	"SwapRate": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "lock_number",
		  "u64"
		],
		[
		  "rate",
		  "u128"
		]
	  ]
	},
	"ClaimInfo": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "mint_amount",
		  "u128"
		],
		[
		  "native_token_amount",
		  "u128"
		],
		[
		  "total_reward",
		  "Balance"
		],
		[
		  "total_claimed",
		  "Balance"
		],
		[
		  "latest_claimed_block",
		  "BlockNumber"
		],
		[
		  "mint_block",
		  "BlockNumber"
		]
	  ]
	},
	"MintRewardAct": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "begin",
		  "BlockNumber"
		],
		[
		  "end",
		  "BlockNumber"
		],
		[
		  "cycle",
		  "u32"
		],
		[
		  "reward_rate",
		  "u128"
		],
		[
		  "total_reward",
		  "Balance"
		],
		[
		  "left_amount",
		  "Balance"
		],
		[
		  "user_limit",
		  "Balance"
		],
		[
		  "locked_blocks",
		  "u32"
		],
		[
		  "total_rtoken_amount",
		  "u128"
		],
		[
		  "total_native_token_amount",
		  "u128"
		]
	  ]
	},
	"BondSwap": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "bonder",
		  "AccountId"
		],
		[
		  "swap_fee",
		  "Balance"
		],
		[
		  "swap_receiver",
		  "AccountId"
		],
		[
		  "bridger",
		  "AccountId"
		],
		[
		  "recipient",
		  "Vec<u8>"
		],
		[
		  "dest_id",
		  "ChainId"
		],
		[
		  "expire",
		  "BlockNumber"
		],
		[
		  "bond_state",
		  "BondState"
		],
		[
		  "refunded",
		  "bool"
		]
	  ]
	},
	"SwapPool": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "symbol",
		  "RSymbol"
		],
		[
		  "fis_balance",
		  "u128"
		],
		[
		  "rtoken_balance",
		  "u128"
		],
		[
		  "total_unit",
		  "u128"
		]
	  ]
	},
	"StakePool": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "symbol",
		  "RSymbol"
		],
		[
		  "emergency_switch",
		  "bool"
		],
		[
		  "total_stake_lp",
		  "u128"
		],
		[
		  "start_block",
		  "u32"
		],
		[
		  "reward_per_block",
		  "u128"
		],
		[
		  "total_reward",
		  "u128"
		],
		[
		  "left_reward",
		  "u128"
		],
		[
		  "lp_locked_blocks",
		  "u32"
		],
		[
		  "last_reward_block",
		  "u32"
		],
		[
		  "reward_per_share",
		  "u128"
		],
		[
		  "guard_impermanent_loss",
		  "bool"
		]
	  ]
	},
	"StakeUser": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "account",
		  "AccountId"
		],
		[
		  "lp_amount",
		  "u128"
		],
		[
		  "reward_debt",
		  "u128"
		],
		[
		  "reserved_lp_reward",
		  "u128"
		],
		[
		  "total_fis_value",
		  "u128"
		],
		[
		  "total_rtoken_value",
		  "u128"
		],
		[
		  "deposit_height",
		  "u32"
		],
		[
		  "grade_index",
		  "u32"
		],
		[
		  "claimed_reward",
		  "u128"
		]
	  ]
	},
	"AccountLpData": {
	  "type": "struct",
	  "type_mapping": [
		[
		  "free",
		  "u128"
		]
	  ]
	}
  }`
