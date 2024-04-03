package rug

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewRug(t *testing.T) {
	cfg := &config.RugConfig{
		OnChain: config.OnChain{
			URL:     "https://betanet-inner3.artela.network",
			KeyFile: "../../rug.txt",
		},
	}
	s, err := NewRug(nil, cfg)
	require.Equal(t, nil, err)
	defer s.Client().Close()

	task := biz.AddressTask{
		ID:             123,
		GMTCreate:      time.Time{},
		GMTModify:      time.Time{},
		AccountAddress: new(string),
		TaskName:       new(string),
		TaskStatus:     new(string),
		Memo:           new(string),
		Txs:            new(string),
		TaskId:         new(string),
		TaskTopic:      new(string),
		JobBatchId:     new(string),
	}

	hash, err := s.send(task)
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(hash))
	fmt.Println(hash[0].Hex())
}

/* mock area */
func MockAddTasks() {
	for i := 0; i < 1000000; i++ {
		if i%100 == 0 {
			// fmt.Println("---------------------------------------------------------")
			// time.Sleep(time.Second)
		}
		address := "0xEeA0A3FE27A1D63C1BF0a356a6462C9aAdB35217"
		task := biz.AddressTask{
			ID:             int64(i),
			GMTCreate:      time.Time{},
			GMTModify:      time.Time{},
			AccountAddress: &address,
			TaskName:       new(string),
			TaskStatus:     new(string),
			Memo:           new(string),
			Txs:            new(string),
			TaskId:         new(string),
			TaskTopic:      new(string),
			JobBatchId:     new(string),
		}
		// time.Sleep(1 * time.Second)
		taskCh <- task
		// fmt.Println(time.Now().Format("2006.01.02 15:04:05"), "adding task", i)
	}
}

var taskCh = make(chan biz.AddressTask, 100)

func (s *Rug) mockGetTasks(count int) ([]biz.AddressTask, error) {
	tasks := make([]biz.AddressTask, count)

	for i := 0; i < count; i++ {
		tasks[i] = <-taskCh
	}
	return tasks, nil
}

func (s *Rug) mockUpdateTask(task biz.AddressTask, hashs []common.Hash, status *uint64) error {
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

func TestStart(t *testing.T) {
	go MockAddTasks()

	ch := make(chan struct{})
	cfg := &config.RugConfig{
		OnChain: config.OnChain{
			URL:     "http://127.0.0.1:8545",
			KeyFile: "./privateKey.txt",
		},
	}
	cfg.FillDefaults()

	s, err := NewRug(nil, cfg)
	if err != nil {
		panic(err)
	}

	s.RegisterGetTasks(s.mockGetTasks)     // base.RegisterGetTasks(f.getTasks)
	s.RegisterUpdateTask(s.mockUpdateTask) // base.RegisterUpdateTask(f.updateTask)

	s.Start()
	<-ch
}
