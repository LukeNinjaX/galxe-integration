package rug

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
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
	"github.com/artela-network/galxe-integration/contracts/uniswapv2"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/artela-network/galxe-integration/onchain"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

const (
	UINT      = 1000000000000000000
	RugAmount = 100000000
)

type Rug struct {
	sync.Mutex

	url        string
	db         *sql.DB
	client     *goclient.Client
	privateKey *ecdsa.PrivateKey
	publickKey *ecdsa.PublicKey
	nonce      uint64

	contract *uniswapv2.UniswapV2
	cfg      *config.RugConfig

	queue    *llq.Queue
	uptating atomic.Bool
}

func NewRug(db *sql.DB, cfg *config.RugConfig) (*Rug, error) {
	url := cfg.URL
	keyfile := cfg.KeyFile

	cfg.FillDefaults()

	c, err := goclient.NewClient(url)
	if err != nil {
		log.Error("rug module: connect to rpc failed", err)
		return nil, err
	}

	privKey, pubKey, err := goclient.ReadKey(keyfile)
	if err != nil {
		log.Error("rug module: read key failed", err)
		return nil, err
	}

	contractAddress := common.HexToAddress(cfg.ContractAddress)
	instance, err := uniswapv2.NewUniswapV2(contractAddress, c)
	if err != nil {
		log.Error("rug module: load uniswapV2 failed", err)
		return nil, err
	}

	accountAddress := crypto.PubkeyToAddress(*pubKey)
	nonce, err := goclient.Client.NonceAt(*c, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		log.Error("rug module: update nonce failed", err)
		return nil, err
	}

	return &Rug{
		url:        url,
		db:         db,
		client:     c,
		privateKey: privKey,
		publickKey: pubKey,
		nonce:      nonce,
		contract:   instance,
		queue:      llq.New(),
		cfg:        cfg,
	}, nil
}

func (s *Rug) getNonce() uint64 {
	s.Lock()
	defer s.Unlock()
	ret := s.nonce
	s.nonce++
	return ret
}

func (s *Rug) Start() {
	go s.pullTasks()
	go s.handleTasks()
}

func (s *Rug) pullTasks() {
	log.Debug("rug module: starting grab Rug task service...")
	for {
		if s.queue.Size() > onchain.QueueSize {
			time.Sleep(onchain.PullSleep)
			continue
		}

		tasks, err := s.getTasks(onchain.PullBatchCount)
		if err != nil {
			log.Error("rug module: getTasks failed", err)
			time.Sleep(onchain.PullSleep)
			continue
		}

		if len(tasks) == 0 {
			time.Sleep(onchain.PullSleep)
			continue
		}

		log.Debugf("rug module: get %d rug tasks", len(tasks))
		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Rug) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAspectPullTask(s.db, count)
}

func (s *Rug) handleTasks() {
	log.Debug("rug module: starting handling Rug task service...")

	ch := make(chan struct{}, s.cfg.Concurrency)
	for {
		elem, ok := s.queue.Dequeue()
		if !ok || elem == nil {
			time.Sleep(onchain.DeQuequeWait)
			continue
		}

		task, ok := elem.(biz.AddressTask)
		if !ok {
			log.Error("rug module: element from queue is not a task")
			continue
		}
		// s.process(task)

		log.Debug("rug module: processing task", task.ID)
		time.Sleep(onchain.PushSleep)
		hash, err := s.rug(task)
		if err != nil {
			log.Error("transfer err", err)
			if strings.Contains(err.Error(), "invalid nonce") || strings.Contains(err.Error(), "tx already in mempool") {
				// nonce is not match, update the nonce
				s.updateNetwork()
			} else if strings.Contains(err.Error(), "connection refused") {
				// client is disconnected
				s.updateNetwork()
			}
			s.queue.Enqueue(task)
			continue
		}

		err = s.updateTask(task, hash.Hex(), nil)
		if err != nil {
			log.Error("rug module: update task failed", task.ID, err)
			// do not return, still try to get the receipt and update to db again
		}

		ch <- struct{}{}
		go func(task biz.AddressTask, hash common.Hash) {
			s.processReceipt(task, hash)
			<-ch
		}(task, hash)
	}
}

func (s *Rug) updateTask(task biz.AddressTask, hash string, status *uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &hash
	if status != nil {
		taskStatus := string(types.TaskStatusFail)
		if *status == 0 {
			taskStatus = string(types.TaskStatusSuccess)
		}
		req.TaskStatus = &taskStatus
		log.Debugf("update rug task: %d, hash %s, status %s\n", req.ID, *req.Txs, *req.TaskStatus)
	} else {
		log.Debugf("update rug task: %d, hash %s\n", req.ID, *req.Txs)
	}

	return biz.UpdateTask(s.db, req)
}

func (s *Rug) processReceipt(task biz.AddressTask, hash common.Hash) {
	log.Debug("rug module: getting receipt for", task.ID, hash.Hex())

	var networkErr = func(err error) bool {
		log.Error("faucet module: transfer err", err)
		if strings.Contains(err.Error(), "connection refused") {
			// client is disconnected
			s.updateNetwork()
			return true
		}
		return false
	}

	time.Sleep(time.Duration(s.cfg.BlockTime) * time.Millisecond)
	// TODO handle timeout
	for i := 0; i < 50; i++ {
		receipt, err := s.client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			log.Debug("rug module: get receipt failed", hash.Hex(), err)
			time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
			if networkErr(err) {
				i--
			}
			continue
		}
		if receipt == nil {
			continue
		}
		s.updateTask(task, hash.Hex(), &receipt.Status)
		return
	}
	log.Errorf("rug module: failed to get receipt after reaching the upper limit of retry times, task %d, hash %s\n", task.ID, hash.Hex())
	status := uint64(0)
	s.updateTask(task, hash.Hex(), &status)
}

func (s *Rug) rug(task biz.AddressTask) (common.Hash, error) {
	log.Debug("rug module: running rug for", task.ID)
	fromAddress := crypto.PubkeyToAddress(*s.publickKey)
	nonce := s.getNonce()

	// send a tx
	opts := s.client.DefaultTxOpts(s.privateKey, fromAddress, &s.cfg.TxConfig)
	opts.Nonce = big.NewInt(int64(nonce)) // we maintance the nonce ourself
	if len(s.cfg.Path) < 2 {
		panic("config .rug.path is not correct")
	}
	path := make([]common.Address, 2)
	path[0] = common.HexToAddress(s.cfg.Path[0])
	path[1] = common.HexToAddress(s.cfg.Path[1])
	toAddress := fromAddress // rug tokens to the sender
	amount := big.NewInt(1).Mul(big.NewInt(RugAmount), big.NewInt(UINT))
	log.Debugf("rug module: SwapETHForExactTokens, from %s, to %s, amount %d", fromAddress.Hex(), toAddress.Hex(), amount)
	tx, err := s.contract.SwapETHForExactTokens(opts, amount, path, toAddress, big.NewInt(int64(time.Now().Second())+10000))
	if err != nil {
		log.Debug("rug module: submit rug failed", task.ID, err)
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (s *Rug) updateNetwork() {
	if s.uptating.Load() {
		return
	}

	log.Error("rug module: network is not valid, updating network...")
	s.uptating.Store(true)
	defer s.uptating.Store(false)
	for {
		if s.connect() && s.updateNonce() && s.updateContract() {
			log.Info("rug module: network is connected")
			return
		}
		time.Sleep(onchain.Reconnect)
	}
}

func (s *Rug) connect() bool {
	// s.client.Close()

	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("rug module: connect failed ", err)
		return false
	}
	s.client = c
	return true
}

func (s *Rug) updateContract() bool {
	contractAddress := common.HexToAddress(s.cfg.ContractAddress)
	instance, err := uniswapv2.NewUniswapV2(contractAddress, s.client)
	if err != nil {
		log.Error("rug module: load uniswapV2 failed ", err)
		return false
	}

	s.contract = instance
	return true
}

func (s *Rug) updateNonce() bool {
	accountAddress := crypto.PubkeyToAddress(*s.publickKey)
	nonce, err := goclient.Client.NonceAt(*s.client, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		log.Error("rug module: get nonce failed ", err)
		return false
	}
	s.nonce = nonce
	return true
}
