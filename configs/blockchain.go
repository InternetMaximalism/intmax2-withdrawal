package configs

type Blockchain struct {
	ScrollNetworkChainID string `env:"BLOCKCHAIN_SCROLL_NETWORK_CHAIN_ID"`

	RollupContractAddress string `env:"BLOCKCHAIN_ROLLUP_CONTRACT_ADDRESS,required"`

	WithdrawalContractAddress string `env:"BLOCKCHAIN_WITHDRAWAL_CONTRACT_ADDRESS"`
	WithdrawalPrivateKeyHex   string `env:"BLOCKCHAIN_ETHEREUM_WITHDRAWAL_PRIVATE_KEY_HEX"`

	WithdrawalAggregatorThreshold        uint64 `env:"BLOCKCHAIN_WITHDRAWAL_AGGREGATOR_THRESHOLD" envDefault:"8"`
	WithdrawalAggregatorMinutesThreshold uint64 `env:"BLOCKCHAIN_WITHDRAWAL_AGGREGATOR_MINUTES_THRESHOLD" envDefault:"15"`
}
