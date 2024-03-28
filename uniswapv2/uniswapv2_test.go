package uniswapv2

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	c, err := goclient.NewClient("http://47.251.14.108:8545")
	require.Equal(t, nil, err)
	defer c.Close()

	privKey, pubKey, err := goclient.ReadKey("../rug.txt")
	require.Equal(t, nil, err)

	// load contract
	address := common.HexToAddress("0xa646F6607af459917EFc14957bADC0Eb87f6dA7c")
	instance, err := NewUniswapV2(address, c)
	require.Equal(t, nil, err)

	cfg := &config.TxConfig{}
	cfg.FillDefaults()

	fromAddress := crypto.PubkeyToAddress(*pubKey)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	require.Equal(t, nil, err)

	// send a tx
	opts := c.DefaultTxOpts(privKey, fromAddress, cfg)
	opts.Nonce = big.NewInt(int64(nonce)) // we maintance the nonce ourself
	path := make([]common.Address, 2)
	path[0] = common.HexToAddress("0xaDCd43c78A914c6B14171aB1380bCFcfa25cd3AD")
	path[1] = common.HexToAddress("0x8997ec639d49D2F08EC0e6b858f36317680A6eE7")
	toAddress := fromAddress // rug tokens to the sender
	tx, err := instance.SwapETHForExactTokens(opts, big.NewInt(100), path, toAddress, big.NewInt(int64(time.Now().Second())+10000))

	require.Equal(t, nil, err)
	require.Equal(t, true, tx != nil)
	require.Equal(t, true, tx.Hash().Hex() != common.Hash{}.Hex())

	time.Sleep(5 * time.Second)

	{
		// try to query this tx
		tx, isPending, err := c.QueryTxByHash(context.Background(), tx.Hash())
		require.Equal(t, nil, err)
		require.Equal(t, false, isPending)

		json, err := json.Marshal(tx)
		require.Equal(t, nil, err)
		fmt.Println(string(json)) //
	}

	{
		// try to get the receipt
		receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
		require.Equal(t, nil, err)
		json, err := json.Marshal(receipt)
		require.Equal(t, nil, err)
		fmt.Println(string(json)) // {"root":"0x","status":"0x0","cumulativeGasUsed":"0x34008","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","logs":[],"transactionHash":"0x9bb3aeb3c73d7358baae38cfcd7165cc2ecf88d7ece1f9ba064f025c7d8b3dbd","contractAddress":"0x0000000000000000000000000000000000000000","gasUsed":"0x493e0","effectiveGasPrice":null,"blockHash":"0xc3ffa06a3076257aa85eda2bf2f0ef763c3e6490b91f9a6087ee6aeff84310c2","blockNumber":"0x491657","transactionIndex":"0x3"}
		require.Equal(t, ethtypes.ReceiptStatusFailed, receipt.Status)
	}
}
