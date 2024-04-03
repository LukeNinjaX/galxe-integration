package updater

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/common"
)

func TestStart(t *testing.T) {
	go MockAddTasks()

	ch := make(chan struct{})
	cfg := &config.UpdaterConfig{
		OnChain: config.OnChain{
			URL: "http://47.251.61.27:8545",
		},
	}
	cfg.FillDefaults()

	s, err := NewUpdater(nil, cfg)
	if err != nil {
		panic(err)
	}

	s.RegisterGetTasks(s.mockGetTasks)     // base.RegisterGetTasks(f.getTasks)
	s.RegisterUpdateTask(s.mockUpdateTask) // base.RegisterUpdateTask(f.updateTask)

	s.Start()
	<-ch
}

/* mock area */
func MockAddTasks() {
	for i := 0; i < 1000000; i++ {
		if i%100 == 0 {
			// fmt.Println("---------------------------------------------------------")
			// time.Sleep(time.Second)
		}
		address := "0xEeA0A3FE27A1D63C1BF0a356a6462C9aAdB35217"
		txs := "0x0239f148d7b6d5d5b0bd862f794df321d4c8a05fd5c8dfc7423ac1029acafb6b"
		task := biz.AddressTask{
			ID:             int64(i),
			GMTCreate:      time.Time{},
			GMTModify:      time.Time{},
			AccountAddress: &address,
			TaskName:       new(string),
			TaskStatus:     new(string),
			Memo:           new(string),
			Txs:            &txs,
			TaskId:         new(string),
			TaskTopic:      new(string),
			JobBatchId:     new(string),
		}
		taskCh <- task
		// fmt.Println(time.Now().Format("2006.01.02 15:04:05"), "adding task", i)
	}
}

var taskCh = make(chan biz.AddressTask, 100)

func (s *Updater) mockGetTasks(count int) ([]biz.AddressTask, error) {
	tasks := make([]biz.AddressTask, count)

	for i := 0; i < count; i++ {
		tasks[i] = <-taskCh
	}
	return tasks, nil
}

func (s *Updater) mockUpdateTask(task biz.AddressTask, hashs []common.Hash, status *uint64) error {
	hexs := make([]string, len(hashs))
	for i, hash := range hashs {
		hexs[i] = hash.Hex()
	}
	memo := strings.Join(hexs, ",")
	if status != nil {
		fmt.Println("updating task ", task.ID, memo, *status)
	} else {
		fmt.Println("updating task ", task.ID, memo)
	}

	return nil
}
