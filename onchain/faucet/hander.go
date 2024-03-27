package faucet

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/artela-network/galxe-integration/api"
	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

type Faucet struct {
	sync.Mutex

	url        string
	db         *sql.DB
	client     *goclient.Client
	privateKey *ecdsa.PrivateKey
	publickKey *ecdsa.PublicKey
	nonce      uint64

	queue *llq.Queue
}

func NewFaucet(db *sql.DB) (*Faucet, error) {
	url := "http://47.251.61.27:8545" // TODO from config
	keyfile := "./privateKey.txt"     // TODO

	c, err := goclient.NewClient(url)
	if err != nil {
		return nil, err
	}

	privKey, pubKey, err := goclient.ReadKey(keyfile)
	if err != nil {
		return nil, err
	}

	accountAddress := crypto.PubkeyToAddress(*pubKey)
	nonce, err := goclient.Client.NonceAt(*c, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		return nil, err
	}

	return &Faucet{
		url:        url,
		db:         db,
		client:     c,
		privateKey: privKey,
		publickKey: pubKey,
		nonce:      nonce,
		queue:      llq.New(),
	}, nil
}

func (s *Faucet) getNonce() uint64 {
	s.Lock()
	defer s.Unlock()
	ret := s.nonce
	s.nonce++
	return ret
}

func (s *Faucet) Start() {
	go s.pullTasks()
	go s.handleTasks()
}

func (s *Faucet) pullTasks() {
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

func (s *Faucet) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetFaucetTask(s.db, count)
}

func (s *Faucet) handleTasks() {
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
				receipt, err := s.process(task)
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

func (s *Faucet) updateTask(task biz.AddressTask, hash string, status uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &hash
	taskStatus := *task.TaskStatus
	if status == 0 {
		taskStatus = string(api.TaskStatusFail)
	} else {
		taskStatus = string(api.TaskStatusSuccess)
	}
	req.TaskStatus = &taskStatus

	return biz.UpdateTask(s.db, req)
}

func (s *Faucet) process(task biz.AddressTask) (*types.Receipt, error) {
	fmt.Println("processing task...", task.ID)
	hash, err := s.client.Transfer(s.privateKey, common.HexToAddress(*task.AccountAddress), 1, s.getNonce())
	if err != nil {
		return nil, err
	}
	time.Sleep(BlockTime)
	// TODO handle timeout
	for i := 0; i < 10; i++ {
		receipt, err := s.client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			log.Debug("get receipt failed", hash.Hex(), err)
			time.Sleep(GetReceiptInterval)
			continue
		}
		return receipt, nil
	}
	return nil, errors.New("failed to get receipt, retry timeout")
}

func (s *Faucet) updateNonce() {
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

func (s *Faucet) connect() {
	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("connect failed")
		return
	}
	s.client = c
}
