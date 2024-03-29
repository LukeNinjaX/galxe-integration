package rug

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
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
	"github.com/artela-network/galxe-integration/contracts/uniswapv2"
	"github.com/artela-network/galxe-integration/goclient"

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

	queue *llq.Queue
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

		log.Debugf("rug module: get %d rug tasks\n", len(tasks))
		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Rug) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAspectPullTask(s.db, count)
}

func (s *Rug) handleTasks() {
	log.Debug("starting handling Rug task service...")
	for {
		var wg sync.WaitGroup

		for i := 0; i < s.cfg.PushBatchCount; i++ {
			elem, ok := s.queue.Dequeue()
			if !ok || elem == nil {
				continue
			}

			task, ok := elem.(biz.AddressTask)
			if !ok {
				log.Error("rug module: element from queue is not a task")
				continue
			}
			// s.process(task)

			log.Debug("rug module: processing task", task.ID)
			hash, err := s.rug(task)
			if err != nil {
				log.Error("transfer err", err)
				if strings.Contains(err.Error(), "invalid nonce") || strings.Contains(err.Error(), "tx already in mempool") {
					// nonce is not match, update the nonce
					s.updateNonce()
				} else if strings.Contains(err.Error(), "connection refused") {
					// client is disconnected
					s.connect()
				}
				s.queue.Enqueue(task)
				continue
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

func (s *Rug) updateTask(task biz.AddressTask, hash string, status uint64) error {
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

	log.Debugf("update rug task: %d, hash %s, status %s\n", req.ID, *req.Txs, *req.TaskStatus)
	return biz.UpdateTask(s.db, req)
}

func (s *Rug) processReceipt(task biz.AddressTask, hash common.Hash) {
	log.Debug("rug module: getting receipt for", task.ID, hash.Hex())
	time.Sleep(time.Duration(s.cfg.BlockTime) * time.Millisecond)
	// TODO handle timeout
	for i := 0; i < 50; i++ {
		receipt, err := s.client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			log.Debug("rug module: get receipt failed", hash.Hex(), err)
			time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
			continue
		}
		s.updateTask(task, receipt.TxHash.Hex(), receipt.Status)
		return
	}
	log.Error("rug module: failed to get receipt after reaching the upper limit of retry times")
	s.updateTask(task, hash.Hex(), 0)
}

func (s *Rug) updateNonce() {
	log.Debug("rug module: updating nonce")
	accountAddress := crypto.PubkeyToAddress(*s.publickKey)
	nonce, err := goclient.Client.NonceAt(*s.client, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		log.Error("rug module: get nonce failed")
		// try to reconnect the client
		s.connect()
		time.Sleep(100 * time.Millisecond)
	}
	log.Debug("rug module: new nonce", nonce)
	s.nonce = nonce
}

func (s *Rug) connect() {
	log.Debug("rug module: connecting client")
	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("connect failed")
		return
	}
	s.client = c
	s.updateNonce()
	contractAddress := common.HexToAddress(s.cfg.ContractAddress)
	instance, err := uniswapv2.NewUniswapV2(contractAddress, c)
	if err != nil {
		log.Error("connect failed")
		return
	}
	log.Debug("rug module: client is connected")
	s.contract = instance
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
	log.Debugf("rug module: SwapETHForExactTokens, from %s, to %s, amount %d\n", fromAddress.Hex(), toAddress.Hex(), amount)
	tx, err := s.contract.SwapETHForExactTokens(opts, amount, path, toAddress, big.NewInt(int64(time.Now().Second())+10000))
	if err != nil {
		log.Debug("rug module: submit rug failed", task.ID, err)
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}
