package uniswapv2

import (
	"testing"

	"github.com/artela-network/galxe-integration/goclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	c, err := goclient.NewClient("http://47.251.61.27:8545")
	require.Equal(t, nil, err)
	defer c.Close()

	privKey, pubKey, err := goclient.ReadKey("../privateKey.txt")
	require.Equal(t, nil, err)

	// load contract
	address := common.HexToAddress("0xf7D3D6e7D0FC0f773f10208fbF33Eaa0B6E608d6")
	instance, err := NewUniswapV2(address, c)
	require.Equal(t, nil, err)

	_ = privKey
	_ = pubKey
	_ = instance

	// send a tx
	/*fromAddress := crypto.PubkeyToAddress(*pubKey)
	opts := c.DefaultTxOpts(privKey, fromAddress)
	tx, err := instance.

	storeTx, err := instance.AddLiquidity(opts, )
	require.Equal(t, nil, err)
	require.Equal(t, true, storeTx != nil)
	require.Equal(t, true, storeTx.Hash().Hex() != common.Hash{}.Hex())

	time.Sleep(2 * time.Second)

	{
		// try to query this tx
		tx, isPending, err := c.QueryTxByHash(context.Background(), storeTx.Hash())
		require.Equal(t, nil, err)
		require.Equal(t, false, isPending)

		json, err := json.Marshal(tx)
		require.Equal(t, nil, err)
		fmt.Println(string(json)) //
	}

	{
		// try to get the receipt
		receipt, err := c.TransactionReceipt(context.Background(), storeTx.Hash())
		require.Equal(t, nil, err)
		json, err := json.Marshal(receipt)
		require.Equal(t, nil, err)
		fmt.Println(string(json)) // {"root":"0x","status":"0x1","cumulativeGasUsed":"0x249f0","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","logs":[],"transactionHash":"0xb52e3c6750173bc19390fb79e25aa96194294394291a69a3283f9890fc76f280","contractAddress":"0x0000000000000000000000000000000000000000","gasUsed":"0x493e0","effectiveGasPrice":null,"blockHash":"0x431264894b4228738a5771f38184006820fe1770044ded69e142b0c4c094fca0","blockNumber":"0x1f9413","transactionIndex":"0x0"}
		require.Equal(t, ethtypes.ReceiptStatusSuccessful, receipt.Status)
	}*/
}
