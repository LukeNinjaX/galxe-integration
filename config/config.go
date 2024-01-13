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
