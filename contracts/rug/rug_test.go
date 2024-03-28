package rug

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	c, err := goclient.NewClient("http://47.251.61.27:8545")
	require.Equal(t, nil, err)
	defer c.Close()

	privKey, pubKey, err := goclient.ReadKey("../../privateKey.txt")
	require.Equal(t, nil, err)

	var address common.Address
	{
		cfg := &config.TxConfig{}
		cfg.FillDefaults()

		fromAddress := crypto.PubkeyToAddress(*pubKey)
		opts := c.DefaultTxOpts(privKey, fromAddress, cfg)
		nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
		require.Equal(t, nil, err)
		opts.Nonce = big.NewInt(int64(nonce))
		// input := "1.0"

		addr, _, _, err := DeployRug(opts, c, "rug", "RUG")
		require.Equal(t, nil, err)
		address = addr
	}

	fmt.Println(address.Hex())
	// load contract
	instance, err := NewRug(address, c)
	require.Equal(t, nil, err)
	_ = instance

}
