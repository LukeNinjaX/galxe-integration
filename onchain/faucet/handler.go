package faucet

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/contracts/rug"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/artela-network/galxe-integration/onchain"
	coretypes "github.com/ethereum/go-ethereum/core/types"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

const UINT = 1000000000000000000

type Faucet struct {
	sync.Mutex

	url        string
	db         *sql.DB
	client     *goclient.Client
	privateKey *ecdsa.PrivateKey
	publickKey *ecdsa.PublicKey
	nonce      uint64
	rug        *rug.Rug

	cfg *config.FaucetConfig

	queue    *llq.Queue
	uptating atomic.Bool
}

func NewFaucet(db *sql.DB, cfg *config.FaucetConfig) (*Faucet, error) {
	url := cfg.URL
	keyfile := cfg.KeyFile

	cfg.FillDefaults()

	c, err := goclient.NewClient(url)
	if err != nil {
		return nil, err
	}

	instance, err := rug.NewRug(common.HexToAddress(cfg.RugAddress), c)
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
		rug:        instance,
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
	log.Debug("faucet module: starting grab faucet task service...")
	for {
		if s.queue.Size() > onchain.QueueSize {
			time.Sleep(onchain.PullSleep)
			continue
		}

		tasks, err := s.getTasks(onchain.PullBatchCount)
		if err != nil {
			log.Error("faucet module: getTasks failed, ", err)
			time.Sleep(onchain.PullSleep)
			continue
		}

		if len(tasks) == 0 {
			time.Sleep(onchain.PullSleep)
			continue
		}

		log.Debugf("faucet module: get %d facuet stasks, ", len(tasks))
		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Faucet) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetFaucetTask(s.db, count)
}

func (s *Faucet) handleTasks() {
	log.Debug("faucet module: starting handling faucet task service...")

	ch := make(chan struct{}, s.cfg.Concurrency)
	for {
		elem, ok := s.queue.Dequeue()
		if !ok || elem == nil {
			time.Sleep(onchain.DeQuequeWait)
			continue
		}

		task, ok := elem.(biz.AddressTask)
		if !ok {
			log.Error("faucet module: element from queue is not a task")
			continue
		}

		var handleError = func(task biz.AddressTask, err error) {
			// log.Error("faucet module: transfer err", err)
			if strings.Contains(err.Error(), "invalid nonce") || strings.Contains(err.Error(), "tx already in mempool") {
				// nonce is not match, try to reconnect and update the nonce
				s.updateNetwork()
			} else if strings.Contains(err.Error(), "connection refused") {
				// client is disconnected
				s.updateNetwork()
			}
			s.queue.Enqueue(task)
		}

		if task.AccountAddress == nil {
			log.Error("faucet module: task AccountAddress is nil", task.ID)
			continue
		}

		fromAddress := crypto.PubkeyToAddress(*s.publickKey)

		// s.process(task)
		artAmount := s.cfg.TransferAmount
		log.Debugf("faucet module: transfering art from %s to %s for %d, amount %d", fromAddress.Hex(), *task.AccountAddress, task.ID, artAmount)
		time.Sleep(onchain.PushSleep)
		hashTransfer, err := s.client.Transfer(s.privateKey, common.HexToAddress(*task.AccountAddress), artAmount, s.getNonce(), &s.cfg.TxConfig)
		if err != nil {
			handleError(task, err)
			continue
		}

		nonce := s.getNonce()
		opts := s.client.DefaultTxOpts(s.privateKey, fromAddress, &s.cfg.TxConfig)
		opts.Nonce = big.NewInt(int64(nonce)) // we maintance the nonce ourself
		toAddress := common.HexToAddress(*task.AccountAddress)
		rugAmount := big.NewInt(1).Mul(big.NewInt(s.cfg.RugAmount), big.NewInt(UINT))
		log.Debugf("faucet module: transfering rug from %s to %s for %d, amount %d", fromAddress.Hex(), toAddress.Hex(), task.ID, rugAmount)
		time.Sleep(onchain.PushSleep)
		txRug, err := s.rug.Transfer(opts, toAddress, rugAmount)
		if err != nil {
			handleError(task, err)
			continue
		}

		err = s.updateTask(task, memo(hashTransfer, txRug.Hash()), nil)
		if err != nil {
			log.Error("faucet module: update task failed", task.ID, err)
		}

		ch <- struct{}{}
		go func(task biz.AddressTask, hashTransfer, hashRug common.Hash) {
			s.processReceipt(task, hashTransfer, hashRug)
			<-ch
		}(task, hashTransfer, txRug.Hash())
	}
}

func (s *Faucet) updateTask(task biz.AddressTask, memo string, status *uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &memo
	if status != nil {
		taskStatus := string(types.TaskStatusFail)
		if *status == 1 {
			taskStatus = string(types.TaskStatusSuccess)
		}
		req.TaskStatus = &taskStatus
		log.Debugf("faucet module: updating task, %d, hash %s, status %s\n", req.ID, *req.Txs, *req.TaskStatus)
	} else {
		log.Debugf("faucet module: updating task, %d, hash %s\n", req.ID, *req.Txs)
	}

	return biz.UpdateTask(s.db, req)
}

func (s *Faucet) processReceipt(task biz.AddressTask, hashTransfer, hashRug common.Hash) {
	var networkErr = func(err error) bool {
		// log.Error("faucet module: get receipt err", err)
		if strings.Contains(err.Error(), "connection refused") {
			// client is disconnected
			s.updateNetwork()
			return true
		}
		return false
	}

	log.Debugf("faucet module: getting receipt for task, %d, trasfer hash %s, rug hash %s", task.ID, hashTransfer.Hex(), hashRug.Hex())
	time.Sleep(time.Duration(s.cfg.BlockTime) * time.Millisecond)
	// TODO handle timeout
	var transferReceipt, rugReceipt *coretypes.Receipt
	var err error
	for i := 0; i < 50; i++ {
		if transferReceipt == nil {
			transferReceipt, err = s.client.TransactionReceipt(context.Background(), hashTransfer)
			if err != nil {
				time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
				if networkErr(err) {
					i--
				}
				continue
			}
		}

		if rugReceipt == nil {
			rugReceipt, err = s.client.TransactionReceipt(context.Background(), hashRug)
			if err != nil {
				time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
				if networkErr(err) {
					i--
				}
				continue
			}
		}

		if transferReceipt != nil && rugReceipt != nil {
			status := transferReceipt.Status
			if rugReceipt.Status != 1 {
				status = rugReceipt.Status
			}

			s.updateTask(task, memo(hashTransfer, hashRug), &status)
			return
		}
	}
	log.Errorf("faucet module: failed to get receipt after reaching the upper limit of retry times, task %d, hash %s\n", task.ID, memo(hashTransfer, hashRug))
	status := uint64(0)

	s.updateTask(task, memo(hashTransfer, hashRug), &status)
}

func memo(hashes ...common.Hash) string {
	if len(hashes) != 2 {
		return ""
	}
	ret := fmt.Sprintf("%s,%s", hashes[0].Hex(), hashes[1].Hex())
	return ret
}

func (s *Faucet) updateNetwork() {
	if s.uptating.Load() {
		return
	}

	log.Error("faucet module: network is not valid, updating network...")
	s.uptating.Store(true)
	defer s.uptating.Store(false)
	for {
		if s.connect() && s.updateContract() && s.updateNonce() {
			log.Info("faucet module: network is connected")
			return
		}
		time.Sleep(onchain.Reconnect)
	}
}

func (s *Faucet) connect() bool {
	// s.client.Close()

	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("faucet module: connect failed")
		return false
	}
	s.client = c
	return true
}

func (s *Faucet) updateContract() bool {
	instance, err := rug.NewRug(common.HexToAddress(s.cfg.RugAddress), s.client)
	if err != nil {
		log.Error("faucet module: load rug contract failed,", err)
		return false
	}
	s.rug = instance
	return true
}

func (s *Faucet) updateNonce() bool {
	accountAddress := crypto.PubkeyToAddress(*s.publickKey)
	nonce, err := goclient.Client.NonceAt(*s.client, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		log.Error("faucet module: get nonce failed, ", err)
		return false
	}
	s.nonce = nonce
	return true
}
