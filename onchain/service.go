package onchain

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/goclient"
	coretypes "github.com/ethereum/go-ethereum/core/types"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

type Base struct {
	sync.Mutex

	url        string
	db         *sql.DB
	client     *goclient.Client
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	nonce      uint64

	queue    *llq.Queue
	uptating atomic.Bool
	conf     *config.OnChain

	send           Send
	getTasks       GetTasks
	updateTask     UpdateTask
	refreshNetwork RefreshNetwork

	query bool // if this service is a query service
}

func NewBase(db *sql.DB, conf *config.OnChain, query bool) (*Base, error) {
	url := conf.URL

	conf.FillDefaults()

	c, err := goclient.NewClient(url)
	if err != nil {
		return nil, err
	}

	base := &Base{
		url:    url,
		db:     db,
		client: c,
		queue:  llq.New(),
		conf:   conf,
		query:  query,
	}

	if !query {
		keyfile := conf.KeyFile
		privKey, pubKey, err := goclient.ReadKey(keyfile)
		if err != nil {
			return nil, err
		}

		accountAddress := crypto.PubkeyToAddress(*pubKey)
		nonce, err := goclient.Client.NonceAt(*c, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
		if err != nil {
			return nil, err
		}
		base.privateKey = privKey
		base.publicKey = pubKey
		base.nonce = nonce
	}

	return base, nil
}

func (s *Base) Client() *goclient.Client {
	return s.client
}

func (s *Base) DB() *sql.DB {
	return s.db
}

func (s *Base) Privatekey() *ecdsa.PrivateKey {
	return s.privateKey
}

func (s *Base) Publickey() *ecdsa.PublicKey {
	return s.publicKey
}

func (s *Base) RegisterSend(fn Send) {
	s.send = fn
}

func (s *Base) RegisterGetTasks(fn GetTasks) {
	s.getTasks = fn
}

func (s *Base) RegisterUpdateTask(fn UpdateTask) {

	s.updateTask = fn
}

func (s *Base) RegisterRefreshNetwork(fn RefreshNetwork) {
	s.refreshNetwork = fn
}

func (s *Base) GetNonce() uint64 {
	if s.query {
		return 0
	}

	s.Lock()
	defer s.Unlock()
	ret := s.nonce
	s.nonce++
	return ret
}

func (s *Base) DefaultOpts(txConf *config.TxConfig) *bind.TransactOpts {
	if s.query {
		return nil
	}

	fromAddress := crypto.PubkeyToAddress(*s.publicKey)
	nonce := s.GetNonce()

	opts := s.client.DefaultTxOpts(s.privateKey, fromAddress, txConf)
	opts.Nonce = big.NewInt(int64(nonce)) // we maintance the nonce ourself
	return opts
}

func (s *Base) Start() {
	go s.pullTasks()
	go s.handleTasks()
}

func (s *Base) pullTasks() {
	log.Debug("Base module: starting grab Base task service...")
	for {
		if s.queue.Size() > QueueSize {
			time.Sleep(PullSleep)
			continue
		}

		tasks, err := s.getTasks(PullBatchCount)
		if err != nil {
			log.Error("Base module: getTasks failed, ", err)
			time.Sleep(PullSleep)
			continue
		}

		if len(tasks) == 0 {
			time.Sleep(PullSleep)
			continue
		}

		log.Debugf("Base module: get %d facuet stasks, ", len(tasks))
		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Base) handleTasks() {
	log.Debug("Base module: starting handling Base task service...")

	ch := make(chan struct{}, s.conf.Concurrency)
	for {
		elem, ok := s.queue.Dequeue()
		if !ok || elem == nil {
			time.Sleep(DeQuequeWait)
			continue
		}

		task, ok := elem.(biz.AddressTask)
		if !ok {
			log.Error("Base module: element from queue is not a task")
			continue
		}

		var handleError = func(task biz.AddressTask, err error) {
			// log.Error("Base module: transfer err", err)
			if strings.Contains(err.Error(), "invalid nonce") || strings.Contains(err.Error(), "tx already in mempool") {
				// nonce is not match, try to reconnect and update the nonce
				s.updateNetwork()
			} else if strings.Contains(err.Error(), "connection refused") {
				// client is disconnected
				s.updateNetwork()
			}
			if err != ErrInvalidTask {
				s.queue.Enqueue(task)
			}
		}

		ch <- struct{}{}

		next := time.Now().Add(time.Duration(s.conf.SendInterval) * time.Millisecond)
		hashs, err := s.send(task)
		time.Sleep(time.Until(next))
		if err != nil {
			handleError(task, err)
			<-ch
			continue
		}

		err = s.updateTask(task, hashs, nil)
		if err != nil {
			log.Error("Base module: update task failed", task.ID, err)
		}

		go func(task biz.AddressTask, hashs []common.Hash) {
			s.processReceipt(task, hashs)
			<-ch
		}(task, hashs)
	}
}

func (s *Base) processReceipt(task biz.AddressTask, hashs []common.Hash) {
	var networkErr = func(err error) bool {
		// log.Error("Base module: get receipt err", err)
		if strings.Contains(err.Error(), "connection refused") {
			// client is disconnected
			s.updateNetwork()
			return true
		}
		return false
	}

	time.Sleep(time.Duration(s.conf.BlockTime) * time.Millisecond)

	receipts := make(map[string]*coretypes.Receipt, len(hashs))

	for i := 0; i < 50; i++ {
		allSuccess := true
		for _, hash := range hashs {
			if receipt, ok := receipts[hash.Hex()]; !ok || receipt == nil {
				recp, err := s.client.TransactionReceipt(context.Background(), hash)
				if err != nil {
					time.Sleep(time.Duration(s.conf.GetReceiptInterval) * time.Millisecond)
					if networkErr(err) {
						i--
					}
					allSuccess = false
					break
				}
				receipts[hash.Hex()] = recp
			}
		}

		if allSuccess {
			status := uint64(1)
			for _, receipt := range receipts {
				if receipt.Status != 1 {
					status = 0
					break
				}
			}
			s.updateTask(task, hashs, &status)
			return
		}
	}

	status := uint64(0)
	s.updateTask(task, hashs, &status)
}

func (s *Base) updateNetwork() {
	if s.uptating.Load() {
		return
	}

	log.Error("Base module: network is not valid, updating network...")
	s.uptating.Store(true)
	defer s.uptating.Store(false)
	for {
		if s.connect() && s.refreshNetwork() && s.updateNonce() {
			log.Info("Base module: network is connected")
			return
		}
		time.Sleep(Reconnect)
	}
}

func (s *Base) connect() bool {
	// s.client.Close()

	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("Base module: connect failed")
		return false
	}
	s.client = c
	return true
}

func (s *Base) updateNonce() bool {
	if s.query {
		return true
	}

	accountAddress := crypto.PubkeyToAddress(*s.publicKey)
	nonce, err := goclient.Client.NonceAt(*s.client, context.Background(), accountAddress, big.NewInt(rpc.LatestBlockNumber.Int64()))
	if err != nil {
		log.Error("Base module: get nonce failed, ", err)
		return false
	}
	s.nonce = nonce
	return true
}
