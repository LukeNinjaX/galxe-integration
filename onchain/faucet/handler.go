package faucet

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/goclient"

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

	cfg *config.FaucetConfig

	queue *llq.Queue
}

func NewFaucet(db *sql.DB, cfg *config.FaucetConfig) (*Faucet, error) {
	url := cfg.URL
	keyfile := cfg.KeyFile

	cfg.FillDefaults()

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
		cfg:        cfg,
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
	log.Debug("starting grab faucet task service...")
	for {
		if s.queue.Size() > s.cfg.QueueMaxSize {
			time.Sleep(time.Duration(s.cfg.PullInterval) * time.Millisecond)
			continue
		}

		tasks, err := s.getTasks(s.cfg.PullBatchCount)
		if err != nil {
			log.Error("getTasks failed", err)
			time.Sleep(time.Duration(s.cfg.PullInterval) * time.Millisecond)
			continue
		}

		if len(tasks) == 0 {
			time.Sleep(time.Duration(s.cfg.PullInterval) * time.Millisecond)
			continue
		}

		log.Debugf("get %d facuet stasks\n", len(tasks))
		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Faucet) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetFaucetTask(s.db, count)
}

func (s *Faucet) handleTasks() {
	log.Debug("starting handling faucet task service...")
	for {
		var wg sync.WaitGroup

		for i := 0; i < s.cfg.PushBatchCount; i++ {
			elem, ok := s.queue.Dequeue()
			if !ok {
				break
			}

			task := elem.(biz.AddressTask)
			// s.process(task)
			fmt.Println("processing task...", task.ID)
			hash, err := s.client.Transfer(s.privateKey, common.HexToAddress(*task.AccountAddress), s.cfg.TransferAmount, s.getNonce(), &s.cfg.TxConfig)
			if err != nil {
				log.Error("transfer err", err)
				if strings.Contains(err.Error(), "invalid nonce") || strings.Contains(err.Error(), "tx already in mempool") {
					// nonce is not match, update the nonce
					s.updateNonce()
				} else if strings.Contains(err.Error(), "connected") { // TODO fix error string
					// client is disconnected
					s.connect()
				}
				s.queue.Enqueue(task) // TODO add retry limition
			}

			wg.Add(1)
			go func(task biz.AddressTask, hash common.Hash) {
				s.processReceipt(task, hash)
				wg.Done()
			}(task, hash)
		}
		wg.Wait()
		time.Sleep(time.Duration(s.cfg.PushInterval) * time.Millisecond)
	}
}

func (s *Faucet) updateTask(task biz.AddressTask, hash string, status uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &hash
	taskStatus := *task.TaskStatus
	if status == 0 {
		taskStatus = string(types.TaskStatusFail)
	} else {
		taskStatus = string(types.TaskStatusSuccess)
	}
	req.TaskStatus = &taskStatus

	return biz.UpdateTask(s.db, req)
}

func (s *Faucet) processReceipt(task biz.AddressTask, hash common.Hash) {
	time.Sleep(time.Duration(s.cfg.BlockTime) * time.Millisecond)
	// TODO handle timeout
	for i := 0; i < 10; i++ {
		receipt, err := s.client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			log.Debug("get receipt failed", hash.Hex(), err)
			time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
			continue
		}
		s.updateTask(task, receipt.TxHash.Hex(), receipt.Status)
		return
	}
	log.Error("failed to get receipt after reaching the upper limit of retry times")
	s.updateTask(task, hash.Hex(), 0)
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
