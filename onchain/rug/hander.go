package rug

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/artela-network/galxe-integration/uniswapv2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

type Rug struct {
	sync.Mutex

	url        string
	db         *sql.DB
	client     *goclient.Client
	contract   *uniswapv2.UniswapV2
	privateKey *ecdsa.PrivateKey
	publickKey *ecdsa.PublicKey
	nonce      uint64

	queue *llq.Queue
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

	accountAddress := crypto.PubkeyToAddress(*pubKey)
	nonce, err := goclient.Client.NonceAt(*c, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		return nil, err
	}

	return &Rug{
		url:        url,
		db:         db,
		client:     c,
		contract:   instance,
		privateKey: privKey,
		publickKey: pubKey,
		nonce:      nonce,
		queue:      llq.New(),
	}, nil
}

func (s *Rug) getNonce() uint64 {
	s.Lock()
	defer s.Unlock()
	s.nonce++
	return s.nonce
}

func (s *Rug) Start() {
	go s.pullTasks()
	go s.handleTasks()
}

func (s *Rug) pullTasks() {
	for {
		if s.queue.Size() > QueueMaxSize {
			time.Sleep(PullInterval)
			continue
		}

		tasks, err := s.getTasks(PullBatchCount)
		if err != nil {
			log.Error("getTasks failed", err)
			time.Sleep(PullInterval)
			continue
		}

		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Rug) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAspectPullTask(s.db, count)
}

func (s *Rug) handleTasks() {
	for {
		var wg sync.WaitGroup

		for i := 0; i < PushBatchCount; i++ {
			elem, ok := s.queue.Dequeue()
			if !ok {
				break
			}

			wg.Add(1)
			task := elem.(biz.AddressTask)
			go func(task biz.AddressTask) {
				receipt, err := s.process(task) // TODO handle timeout
				if err != nil {
					if strings.Contains(err.Error(), "nonce") { // TODO fix error string
						// nonce is not match, update the nonce
						s.updateNonce()
					} else if strings.Contains(err.Error(), "connected") { // TODO fix error string
						// client is disconnected
						s.connect()
					}
					s.queue.Enqueue(task)
				}
				s.updateTask(task, receipt.TxHash.Hex(), receipt.Status)
				wg.Done()
			}(task)
		}
		wg.Wait()
		time.Sleep(PushInterval)
	}
}

func (s *Rug) updateTask(task biz.AddressTask, hash string, status uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &hash
	if status == 0 {
		req.TaskStatus = &TaskStatusFail
	} else {
		req.TaskStatus = &TaskStatusSuccess
	}

	return biz.UpdateTask(s.db, req)
}

func (s *Rug) process(task biz.AddressTask) (*types.Receipt, error) {
	fromAddress := crypto.PubkeyToAddress(*s.publickKey)
	opts := s.client.DefaultTxOpts(s.privateKey, fromAddress) // TODO optimize
	opts.Nonce = big.NewInt(int64(s.getNonce()))
	// tx, err := s.contract.AddLiquidity(opts)
	return nil, nil
}

func (s *Rug) updateNonce() {
	accountAddress := crypto.PubkeyToAddress(*s.publickKey)
	nonce, err := goclient.Client.NonceAt(*s.client, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		log.Error("get nonce failed")
		// try to reconnect the client
		s.connect()
		time.Sleep(100 * time.Millisecond)
	}
	s.nonce = nonce
}

func (s *Rug) connect() {
	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("connect failed")
		return
	}
	s.client = c
}
