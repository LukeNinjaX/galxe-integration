package faucet

import (
	"fmt"
	"testing"

	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	cfg := &config.FaucetConfig{
		OnChain: config.OnChain{
			URL:     "http://47.251.58.164:8545",
			KeyFile: "../../privateKey.txt",
		},
		TransferAmount: 1,
	}
	cfg.FillDefaults()

	s, err := NewFaucet(nil, cfg)
	require.Equal(t, nil, err)
	defer s.client.Close()
	for i := 0; i < 1000; i++ {
		fmt.Println("sending transfer", i)
		hash, err := s.client.Transfer(
			s.privateKey,
			common.HexToAddress("0x22b7926DA60F97c1aAD776084174218EbBF05E28"),
			cfg.TransferAmount,
			s.getNonce(),
			&cfg.TxConfig,
		)
		require.Equal(t, nil, err)
		_ = hash
	}
}
