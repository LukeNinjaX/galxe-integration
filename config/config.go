package config

import (
	"encoding/json"
)

type Config struct {
	Notifiers []json.RawMessage `json:"notifiers"`
	Indexers  []*IndexerConfig  `json:"indexers"`
	APIServer *APIConfig        `json:"api_server"`
	Fetcher   *FetcherConfig    `json:"fetcher"`
	DB        *DBConfig         `json:"db"`
	GoPlus    *GoPlusConfig     `json:"biz_goplus"`
	Faucet    *FaucetConfig     `json:"faucet"`
	Rug       *RugConfig        `json:"rug"`
}

type FaucetConfig struct {
	OnChain
	TxConfig
	TransferAmount int64  `json:"transfer_amount"`
	RugAddress     string `json:"rug_address"`
	RugAmount      int64  `json:"rug_amount"`
}

func (c *FaucetConfig) FillDefaults() {
	if c.TransferAmount <= 0 {
		c.TransferAmount = 1
	}
	if c.RugAddress == "" {
		c.RugAddress = "0x8997ec639d49D2F08EC0e6b858f36317680A6eE7"
	}
	if c.RugAmount <= 0 {
		c.RugAmount = 100000
	}
	c.OnChain.FillDefaults()
	c.TxConfig.FillDefaults()
}

type RugConfig struct {
	OnChain
	TxConfig
	ContractAddress string   `json:"contract_address"`
	Path            []string `json:"path"`
}

func (c *RugConfig) FillDefaults() {
	c.OnChain.FillDefaults()
	c.TxConfig.FillDefaults()

	if c.ContractAddress == "" {
		c.ContractAddress = "0xa646F6607af459917EFc14957bADC0Eb87f6dA7c"
	}

	if len(c.Path) == 0 {
		c.Path = make([]string, 2)
		c.Path[0] = "0xaDCd43c78A914c6B14171aB1380bCFcfa25cd3AD"
		c.Path[1] = "0x8997ec639d49D2F08EC0e6b858f36317680A6eE7"
	}
}

type TxConfig struct {
	GasLimit int64 `json:"gas_limit"`
	GasPrice int64 `json:"gas_price"`
	ChainID  int64 `json:"chain_id"`
}

func (c *TxConfig) FillDefaults() {
	if c.GasLimit == 0 {
		c.GasLimit = 300000
	}

	if c.GasPrice == 0 {
		c.GasPrice = 7
	}

	if c.ChainID == 0 {
		c.ChainID = 11822
	}
}

type OnChain struct {
	URL                string `json:"url"`
	KeyFile            string `json:"keyfile"`
	PullInterval       int    `json:"pull_interval"`
	PullBatchCount     int    `json:"pull_batch_count"`
	PushInterval       int    `json:"push_interval"`
	PushBatchCount     int    `json:"push_batch_count"`
	QueueMaxSize       int    `json:"queue_max_size"`
	BlockTime          int    `json:"block_time"`
	GetReceiptInterval int    `json:"get_receipt_interval"`
}

func (c *OnChain) FillDefaults() {
	if c.PullInterval <= 0 {
		c.PullInterval = 1000
	}

	if c.PullBatchCount <= 0 {
		c.PullBatchCount = 20
	}

	if c.PushInterval <= 0 {
		c.PushInterval = 1000
	}

	if c.PushBatchCount <= 0 {
		c.PushBatchCount = 50
	}

	if c.QueueMaxSize <= 0 {
		c.QueueMaxSize = 200
	}

	if c.BlockTime <= 0 {
		c.BlockTime = 600
	}

	if c.GetReceiptInterval <= 0 {
		c.GetReceiptInterval = 100
	}
}

type GoPlusConfig struct {
	ChannelCode string `json:"channelCode"`
	ManageId    string `json:"manageId"`
	ManageKey   string `json:"manageKey"`
	SecwarexUrl string `json:"secwarexUrl"`
}
type DBConfig struct {
	URL           string `json:"url"`
	MaxConnection uint64 `json:"max_connection"`
}

type IndexerConfig struct {
	Type     string `json:"type"`
	Thread   uint64 `json:"thread"`
	Contract string `json:"contract"`
}

type FetcherConfig struct {
	EthereumRPCUrl    string `json:"ethereum_rpc_url"`
	PullIntervalMs    uint64 `json:"pull_interval_ms"`
	RetryIntervalMs   uint64 `json:"retry_interval_ms"`
	BeginBlock        uint64 `json:"begin_block"`
	BlockCacheSize    uint64 `json:"block_cache_size"`
	PollThread        uint64 `json:"poll_thread"`
	BlockMaxRetry     uint64 `json:"block_max_retry"`
	MaxProcessingTime string `json:"max_processing_time"`
}

func (c *FetcherConfig) FillDefaults() *FetcherConfig {
	if c.EthereumRPCUrl == "" {
		c.EthereumRPCUrl = "http://localhost:8545"
	}
	if c.BlockCacheSize == 0 {
		c.BlockCacheSize = 100
	}
	if c.PullIntervalMs == 0 {
		c.PullIntervalMs = 300
	}
	if c.RetryIntervalMs == 0 {
		c.RetryIntervalMs = 200
	}
	if c.BeginBlock == 0 {
		c.BeginBlock = 1
	}
	if c.PollThread == 0 {
		c.PollThread = 10
	}
	if c.BlockMaxRetry == 0 {
		c.BlockMaxRetry = 3
	}
	if c.MaxProcessingTime == "" {
		c.MaxProcessingTime = "5m"
	}
	return c
}

type APIConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type TypeConf struct {
	Type string `json:"type"`
}
