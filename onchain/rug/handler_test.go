package rug

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/artela-network/galxe-integration/config"
	"github.com/stretchr/testify/require"
)

func TestNewRug(t *testing.T) {
	cfg := &config.RugConfig{
		OnChain: config.OnChain{
			URL:     "http://47.251.58.164:8545",
			KeyFile: "../../privateKey.txt",
		},
		ContractAddress: "0x",
	}
	s, err := NewRug(nil, cfg)
	require.Equal(t, nil, err)
	defer s.client.Close()

	for i := 0; i < 1000; i++ {
		fmt.Println("sending swap", i)
	}
}

func Test1(t *testing.T) {
	num := big.NewInt(1).Mul(big.NewInt(100000), big.NewInt((1000000000000000000)))
	fmt.Println(num.String())
}
