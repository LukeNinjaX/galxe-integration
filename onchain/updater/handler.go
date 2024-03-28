package updater

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/goclient"
	"github.com/ethereum/go-ethereum/common"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	log "github.com/sirupsen/logrus"
)

type Updater struct {
	url    string
	db     *sql.DB
	client *goclient.Client

	cfg *config.UpdaterConfig

	queue *llq.Queue
}

func NewUpdater(db *sql.DB, cfg *config.UpdaterConfig) (*Updater, error) {
	url := cfg.URL

	c, err := goclient.NewClient(url)
	if err != nil {
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
		if s.queue.Size() > s.cfg.QueueMaxSize {
			log.Debugf("updater module: queue is full, try get tasks later\n")
			time.Sleep(time.Duration(s.cfg.PullInterval) * time.Millisecond)
			continue
		}

		tasks, err := s.getTasks(s.cfg.PullBatchCount)
		if err != nil {
			log.Error("updater module: getTasks failed", err)
			time.Sleep(time.Duration(s.cfg.PullInterval) * time.Millisecond)
			continue
		}

		log.Debugf("updater module: get %d tasks\n", len(tasks))
		if len(tasks) == 0 {
			time.Sleep(time.Duration(s.cfg.PullInterval) * time.Millisecond)
			continue
		}

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
			time.Sleep(time.Duration(s.cfg.PushInterval) * time.Millisecond)
			continue
		}
		task, ok := elem.(biz.AddressTask)
		if !ok {
			log.Debugf("updater module: element is not a AddressTask\n")
			continue
		}
		log.Debugf("updater module: handing task %d, hash: %s\n", task.ID, *task.Txs)
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

	log.Debugf("update addliquidity task: %d, hash %s, status %s\n", req.ID, *task.Txs, *req.TaskStatus)
	return biz.UpdateTask(s.db, req)
}

func (s *Updater) getReceipt(task biz.AddressTask) {
	log.Debugf("updater module: get Receipt for task %d, hash %s\n", task.ID, *task.Txs)
	hash := common.HexToHash(*task.Txs)
	receipt, err := s.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			// client is disconnected
			s.connect()
		}
		log.Debug("updater module: get receipt failed and put back into queue", "task", task.ID, "hash", hash.Hex(), err)
		s.queue.Enqueue(task)
	}

	if receipt == nil {
		s.updateTask(task, 0)
		return
	}

	s.updateTask(task, receipt.Status)
}

func (s *Updater) connect() {
	c, err := goclient.NewClient(s.url)
	if err != nil {
		log.Error("connect failed")
		return
	}
	s.client = c
}
