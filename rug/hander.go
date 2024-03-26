package rug

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"math/big"
	"sync"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/artela-network/galxe-integration/uniswapv2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

const (
	Interval   = 10 * time.Second
	BatchCount = 10
)

type Rug struct {
	sync.Mutex

	db         *sql.DB
	client     *goclient.Client
	contract   *uniswapv2.UniswapV2
	privateKey *ecdsa.PrivateKey
	publickKey *ecdsa.PublicKey
	nonce      uint64
}

func NewRug(db *sql.DB) (*Rug, error) {
	url := "http://47.251.61.27:8545" // TODO from config
	address := ""                     // contract address, TODO from config
	keyfile := "./privateKey.txt"     // TODO

	c, err := goclient.NewClient(url)
	if err != nil {
		return nil, err
	}

	privKey, pubKey, err := goclient.ReadKey(keyfile)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(address)
	instance, err := uniswapv2.NewUniswapV2(contractAddress, c)
	if err != nil {
		return nil, err
	}

	nonce, err := goclient.Client.NonceAt(*c, context.Background(), contractAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		return nil, err
	}

	return &Rug{
		db:         db,
		client:     c,
		contract:   instance,
		privateKey: privKey,
		publickKey: pubKey,
		nonce:      nonce,
	}, nil
}

func (s *Rug) GetNonce() uint64 {
	s.Lock()
	defer s.Unlock()
	s.nonce++
	return s.nonce
}

func (s *Rug) Start() error {
	for {
		tasks, err := s.getTasks(BatchCount)
		if err != nil {
			log.Debugf("getTxs return error, will try agin after %2f seconds, err %v\n", Interval.Seconds(), err)
		}
		if err != nil || len(tasks) == 0 {
			time.Sleep(Interval)
			continue
		}

		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			go func(task biz.AddressTask) {
				receipt, err := s.process(task)
				if err != nil {
					s.updateTask(task, common.Hash{}.Hex(), 0) // TODO
				} else {
					s.updateTask(task, receipt.TxHash.Hex(), receipt.Status)
				}

				wg.Done()
			}(task)
		}

		wg.Wait()
		time.Sleep(Interval)
	}
}

func (s *Rug) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAspectPullTask(s.db, count)
}

func (s *Rug) updateTask(task biz.AddressTask, hash string, status uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &hash

	return biz.UpdateTask(s.db, req)
}

func (s *Rug) process(task biz.AddressTask) (*types.Receipt, error) {
	fromAddress := crypto.PubkeyToAddress(*s.publickKey)
	opts := s.client.DefaultTxOpts(s.privateKey, fromAddress) // TODO optimize
	opts.Nonce = big.NewInt(int64(s.GetNonce()))
	// tx, err := s.contract.AddLiquidity(opts)
	return nil, nil
}
