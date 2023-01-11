package config

const (
	RTokenSeriesModuleId = "RTokenSeries"

	LiquidityBondEventId      = "LiquidityBond"
	ExecuteBondAndSwapEventId = "ExecuteBondAndSwap"
	NominationUpdatedEventId  = "NominationUpdated"
	ValidatorUpdatedEventId   = "ValidatorUpdated"

	RClaimModuleId = "RClaim"

	StorageActLatestCycle     = "ActLatestCycle"
	StorageREthActLatestCycle = "REthActLatestCycle"
	StorageActs               = "Acts"
	StorageREthActs           = "REthActs"

	StorageBondRecords       = "BondRecords"
	StorageBondStates        = "BondStates"
	MethodExecuteBondRecord  = "RTokenSeries.execute_bond_record"
	MethodExecuteBondAndSwap = "RTokenSeries.execute_bond_and_swap"
	StorageNominated         = "Nominated"

	RtokenVoteModuleId         = "RTokenVotes"
	StorageVotes               = "Votes"
	MethodRacknowledgeProposal = "RTokenVotes.acknowledge_proposal"

	RTokenLedgerModuleId   = "RTokenLedger"
	EraPoolUpdatedEventId  = "EraPoolUpdated"
	EraUpdatedEventId      = "EraUpdated"
	BondingDurationEventId = "BondingDurationUpdated"

	RTokenRelayersModuleId                    = "Relayers"
	StorageChainEras                          = "ChainEras"
	StorageCurrentEraSnapShots                = "CurrentEraSnapShots"
	StorageRelayerThreshold                   = "RelayerThreshold"
	StorageEraSnapShots                       = "EraSnapShots"
	StorageLeastBond                          = "LeastBond"
	StoragePendingStake                       = "PendingStake"
	StoragePendingReward                      = "PendingReward"
	MethodSetChainEra                         = "RTokenLedger.set_chain_era"
	MethodBondReport                          = "RTokenLedger.bond_report"
	MethodNewBondReport                       = "RTokenLedger.new_bond_report"
	MethodActiveReport                        = "RTokenLedger.active_report"
	MethodNewActiveReport                     = "RTokenLedger.new_active_report"
	MethodBondAndReportActive                 = "RTokenLedger.bond_and_report_active"
	MethodBondAndReportActiveWithPendingValue = "RTokenLedger.bond_and_report_active_with_pending_value"
	MethodWithdrawReport                      = "RTokenLedger.withdraw_report"
	MethodTransferReport                      = "RTokenLedger.transfer_report"
	BondReportedEventId                       = "BondReported"
	ActiveReportedEventId                     = "ActiveReported"
	WithdrawReportedEventId                   = "WithdrawReported"
	TransferReportedEventId                   = "TransferReported"
	StorageSubAccounts                        = "SubAccounts"
	StorageMultiThresholds                    = "MultiThresholds"
	StorageBondedPools                        = "BondedPools"
	StorageSnapshots                          = "Snapshots"
	StoragePoolUnbonds                        = "PoolUnbonds"
	SignaturesEnoughEventId                   = "SignaturesEnough"
	StorageSignatures                         = "Signatures"
	SubmitSignatures                          = "RTokenSeries.submit_signatures"

	RTokenUnbondEventId = "LiquidityUnBond"

	RTokenBalanceModuleId = "RBalances"
	RTokenTransferEventId = "Transfer"
	RTokenMintedEventId   = "Minted"
	RTokenBurnedEventId   = "Burned"

	RTokenRateModuleId   = "RTokenRate"
	RTokenRateSetEventId = "RateSet"
	StorageEraRate       = "EraRate"

	RFisModuleId      = "RFis"
	RFisUnbondEventId = "LiquidityUnBond"

	EraPayoutEventId = "EraPayout"

	RDexSwapModuleId        = "RDexSwap"
	RDexSwapEventId         = "Swap"
	RDexAddLiquidityEventId = "AddLiquidity"
	RDexRmLiquidityEventId  = "RemoveLiquidity"
)
