package rug

import (
	"testing"

	"github.com/artela-network/galxe-integration/config"
	"github.com/stretchr/testify/require"
)

func TestNewRug(t *testing.T) {
	rug, err := NewRug(nil, &config.RugConfig{URL: "http://47.251.61.27:8545", KeyFile: "../privateKey.txt", ContractAddress: "0x"})
	require.Equal(t, nil, err)

	_ = rug
}
