package faucet

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	cfg := &config.FaucetConfig{
		OnChain: config.OnChain{
			URL:     "http://47.251.32.165:8545",
			KeyFile: "../../rug.txt",
		},
		TransferAmount: 1,
		RugAmount:      1000000000,
		RugAddress:     "0x8997ec639d49D2F08EC0e6b858f36317680A6eE7",
	}
	cfg.FillDefaults()

	s, err := NewFaucet(nil, cfg)
	require.Equal(t, nil, err)
	defer s.Client().Close()
	// for i := 0; i < 1000; i++ {
	// fmt.Println("sending transfer", i)
	// {
	// 	hash, err := s.client.Transfer(
	// 		s.privateKey,
	// 		common.HexToAddress("0x58C1B539B469fd15A02Da47b52A3B82bc2ed2b1a"),
	// 		cfg.TransferAmount,
	// 		s.getNonce(),
	// 		&cfg.TxConfig,
	// 	)
	// 	require.Equal(t, nil, err)
	// 	_ = hash
	// }

	{
		nonce := s.GetNonce()

		fromAddress := crypto.PubkeyToAddress(*s.Publickey())
		opts := s.Client().DefaultTxOpts(s.Privatekey(), fromAddress, &s.conf.TxConfig)
		opts.Nonce = big.NewInt(int64(nonce)) // we maintance the nonce ourself
		toAddress := common.HexToAddress("0x58C1B539B469fd15A02Da47b52A3B82bc2ed2b1a")
		rugAmount := big.NewInt(1).Mul(big.NewInt(s.conf.RugAmount), big.NewInt(onchain.UINT))
		// log.Debugf("faucet module: transfering rug from %s to %s for %d, amount %d", fromAddress.Hex(), toAddress.Hex(), task.ID, rugAmount)
		time.Sleep(onchain.PushSleep)
		tx, err := s.rug.Transfer(opts, toAddress, rugAmount)
		require.Equal(t, nil, err)
		fmt.Println(tx.Hash().Hex())
	}
	// }
}

func TestAddTask(t *testing.T) {
	go MockAddTasks()

	ch := make(chan struct{}, 0)
	cfg := &config.FaucetConfig{
		OnChain: config.OnChain{
			URL:     "http://47.251.58.164:8545",
			KeyFile: "../../privateKey.txt",
		},
		TransferAmount: 1,
		RugAddress:     "0x1f9c0A770a25e37698E54ffbAc0a4AfBa84d2a02",
	}
	cfg.FillDefaults()

	s, err := NewFaucet(nil, cfg)
	require.Equal(t, nil, err)
	defer s.Client().Close()
	s.Start()
	<-ch
}
