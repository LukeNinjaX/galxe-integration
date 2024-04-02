package updater

import (
	"context"
	"database/sql"
	"strings"
	"sync/atomic"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/artela-network/galxe-integration/onchain"
	"github.com/ethereum/go-ethereum/common"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

type Updater struct {
	url    string
	db     *sql.DB
	client *goclient.Client

	cfg *config.UpdaterConfig

	queue    *llq.Queue
	uptating atomic.Bool
}

func NewUpdater(db *sql.DB, cfg *config.UpdaterConfig) (*Updater, error) {
	url := cfg.URL

	c, err := goclient.NewClient(url)
	if err != nil {
		log.Error("update module: connect to rpc failed", err)
		return nil, err
	}

	cfg.FillDefaults()

	return &Updater{
		url:    url,
		db:     db,
		client: c,
		queue:  llq.New(),
		cfg:    cfg,
	}, nil
}

func (s *Updater) Start() {
	go s.pullTasks()
	go s.handleTasks()
}

func (s *Updater) pullTasks() {
	log.Debug("updater module: starting grab Updater task service...")
	for {
		if s.queue.Size() > onchain.QueueSize {
			log.Debugf("updater module: queue is full, try get tasks later")
			time.Sleep(onchain.PullSleep)
			continue
		}

		tasks, err := s.getTasks(onchain.PullBatchCount)
		if err != nil {
			log.Error("updater module: getTasks failed", err)
			time.Sleep(onchain.PullSleep)
			continue
		}

		if len(tasks) == 0 {
			time.Sleep(onchain.PullSleep)
			continue
		}

		log.Debugf("updater module: get %d tasks", len(tasks))
		for _, task := range tasks {
			s.queue.Enqueue(task)
		}
	}
}

func (s *Updater) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAddLiquidityTask(s.db, count)
}

func (s *Updater) handleTasks() {
	log.Debug("updater module: starting handling Updater task service...")

	ch := make(chan struct{}, s.cfg.Concurrency)
	for {
		elem, ok := s.queue.Dequeue()
		if !ok || elem == nil {
			time.Sleep(onchain.DeQuequeWait)
			continue
		}
		task, ok := elem.(biz.AddressTask)
		if !ok {
			log.Debugf("updater module: element is not a AddressTask")
			continue
		}

		if task.Txs == nil {
			log.Debugf("updater module: task.txs cannot be empty")
			continue
		}

		log.Debugf("updater module: handing task %d, hash: %s", task.ID, *task.Txs)
		ch <- struct{}{}
		go func(task biz.AddressTask) {
			s.getReceipt(task)
			<-ch
		}(task)
	}
}

func (s *Updater) updateTask(task biz.AddressTask, status uint64) error {
	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	taskStatus := *task.TaskStatus
	if status == 1 {
		taskStatus = string(types.TaskStatusSuccess)
	} else {
		taskStatus = string(types.TaskStatusFail)
	}
	req.TaskStatus = &taskStatus

	log.Debugf("updater moduler: update addliquidity task: %d, hash %s, status %s\n", req.ID, *task.Txs, *req.TaskStatus)
	return biz.UpdateTask(s.db, req)
}

func (s *Updater) getReceipt(task biz.AddressTask) {
	log.Debugf("updater module: get Receipt for task %d, hash %s", task.ID, *task.Txs)

	var networkErr = func(err error) bool {
		// log.Error("updater module: get receipt err", err)
		if strings.Contains(err.Error(), "connection refused") {
			// client is disconnected
			s.updateNetwork()
			return true
		}
		return false
	}

	hash := common.HexToHash(*task.Txs)
	for i := 0; i < 50; i++ {
		receipt, err := s.client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				// client is disconnected
				s.updateNetwork()
			}
			log.Debug("updater module: get receipt failed and put back into queue", "task", task.ID, "hash", hash.Hex(), err)
			time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
			if networkErr(err) {
				i--
			}
			continue
		}

		if receipt == nil {
			time.Sleep(time.Duration(s.cfg.GetReceiptInterval) * time.Millisecond)
			continue
		}
		s.updateTask(task, receipt.Status)
		return
	}
	log.Errorf("updater module: failed to get receipt after reaching the upper limit of retry times, task %d, hash %s\n", task.ID, hash.Hex())
	status := uint64(0)
	s.updateTask(task, status)
}

func (s *Updater) updateNetwork() {
	if s.uptating.Load() {
		return
	}

	log.Error("updater module: network is not valid, updating network...")
	s.uptating.Store(true)
	defer s.uptating.Store(false)
	for {
		if s.connect() {
			log.Info("updater module: network is connected")
			return
		}
		time.Sleep(onchain.Reconnect)
	}
}

func (s *Updater) connect() bool {
	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("connect failed")
		return false
	}
	s.client = c
	return true
}
