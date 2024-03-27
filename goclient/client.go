package goclient

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	*ethclient.Client
}

// NewClient with url like http://127.0.0.1:8545
func NewClient(url string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Errorf("failed to connect to node %s, err %v", url, err)
		return nil, err
	}

	return &Client{client}, nil
}

func (c *Client) DefaultTxOpts(privateKey *ecdsa.PrivateKey, fromAddress common.Address, conf *config.TxConfig) *bind.TransactOpts {
	// nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil
	// }

	var err error
	gasPrice := big.NewInt(conf.GasPrice)
	if gasPrice.Int64() == 0 {
		gasPrice, err = c.SuggestGasPrice(context.Background())
		if err != nil {
			log.Error(err)
		}
	}

	chainId := big.NewInt(conf.ChainID)
	if chainId.Int64() == 0 {
		chainId, err = c.ChainID(context.Background())
		if err != nil {
			log.Error(err)
		}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)            // in wei
	auth.GasLimit = uint64(conf.GasLimit) // in units
	auth.GasPrice = gasPrice
	return auth
}

func (c *Client) QueryTxByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return c.TransactionByHash(ctx, hash)
}

func (c *Client) SendRawMessage(ctx context.Context, msg []byte) (common.Hash, error) {
	tx := new(types.Transaction)
	err := rlp.DecodeBytes(msg, &tx)
	if err != nil {
		return common.Hash{}, err
	}

	err = c.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *Client) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return c.Client.TransactionReceipt(ctx, hash)
}

func (c *Client) Transfer(privateKey *ecdsa.PrivateKey, to common.Address, amount int64, nonce uint64, conf *config.TxConfig) (common.Hash, error) {
	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	err := errors.New("error casting public key to ECDSA")
	// 	return common.Hash{}, err
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var err error

	value := big.NewInt(1000000000000000000 * amount) // in wei (1 eth)
	gasLimit := uint64(conf.GasLimit)                 // in units

	gasPrice := big.NewInt(conf.GasPrice)
	if gasPrice.Int64() == 0 {
		gasPrice, err = c.SuggestGasPrice(context.Background())
		if err != nil {
			return common.Hash{}, err
		}
	}

	var data []byte
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	chainID := big.NewInt(conf.ChainID)
	if chainID.Int64() == 0 {
		chainID, err = c.NetworkID(context.Background())
		if err != nil {
			return common.Hash{}, err
		}
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, err
	}
	return signedTx.Hash(), nil
}
