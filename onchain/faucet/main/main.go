package main

import (
	"fmt"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/onchain/faucet"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	go MockAddTasks()

	ch := make(chan struct{}, 0)
	cfg := &config.FaucetConfig{
		OnChain: config.OnChain{
			URL:          "http://127.0.0.1:8545",
			KeyFile:      "./privateKey.txt",
			SendInterval: 40,
		},
		TransferAmount: 1,
		RugAddress:     "0x1f9c0A770a25e37698E54ffbAc0a4AfBa84d2a02",
	}
	cfg.FillDefaults()

	s, err := faucet.NewFaucet(nil, cfg)
	if err != nil {
		panic(err)
	}
	s.Start()
	s.RegisterGetTasks(mockGetTasks)     // base.RegisterGetTasks(f.getTasks)
	s.RegisterUpdateTask(mockUpdateTask) // base.RegisterUpdateTask(f.updateTask)

	<-ch
}

// /* mock area */
// func MockAddTasks() {
// 	for i := 0; i < 10000; i++ {
// 		if i%100 == 0 {
// 			// fmt.Println("---------------------------------------------------------")
// 			// time.Sleep(time.Second)
// 		}
// 		address := "0xEeA0A3FE27A1D63C1BF0a356a6462C9aAdB35217"
// 		task := biz.AddressTask{
// 			ID:             int64(i),
// 			GMTCreate:      time.Time{},
// 			GMTModify:      time.Time{},
// 			AccountAddress: &address,
// 			TaskName:       new(string),
// 			TaskStatus:     new(string),
// 			Memo:           new(string),
// 			Txs:            new(string),
// 			TaskId:         new(string),
// 			TaskTopic:      new(string),
// 			JobBatchId:     new(string),
// 		}
// 		taskCh <- task
// 		// fmt.Println(time.Now().Format("2006.01.02 15:04:05"), "adding task", i)
// 	}
// }

// var taskCh = make(chan biz.AddressTask, 10000)

// func (s *Faucet) mockGetTasks(count int) ([]biz.AddressTask, error) {
// 	tasks := make([]biz.AddressTask, count)

// 	for i := 0; i < count; i++ {
// 		tasks[i] = <-taskCh
// 	}
// 	return tasks, nil
// }

// func (s *Faucet) mockUpdateTask(task biz.AddressTask, memo string, status *uint64) error {
// 	req := &biz.UpdateTaskQuery{}
// 	req.ID = task.ID
// 	req.Txs = &memo
// 	if status != nil {
// 		taskStatus := string(types.TaskStatusFail)
// 		if *status == 1 {
// 			taskStatus = string(types.TaskStatusSuccess)
// 		}
// 		req.TaskStatus = &taskStatus
// 		fmt.Printf("faucet module: updating task, %d, hash %s, status %s", req.ID, *req.Txs, *req.TaskStatus)
// 	} else {
// 		fmt.Printf("faucet module: updating task, %d, hash %s", req.ID, *req.Txs)
// 	}

// 	return nil
// }

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
		taskCh <- task
		// fmt.Println(time.Now().Format("2006.01.02 15:04:05"), "adding task", i)
	}
}

var taskCh = make(chan biz.AddressTask, 100)

func mockGetTasks(count int) ([]biz.AddressTask, error) {
	tasks := make([]biz.AddressTask, count)

	for i := 0; i < count; i++ {
		tasks[i] = <-taskCh
	}
	return tasks, nil
}

func mockUpdateTask(task biz.AddressTask, hashs []common.Hash, status *uint64) error {
	hexs := make([]string, len(hashs))
	for i, hash := range hashs {
		hexs[i] = hash.Hex()
	}
	// memo := strings.Join(hexs, ",")
	if status != nil {
		fmt.Println("updating task", task.ID, *status)
	} else {
		fmt.Println("updating task", task.ID)
	}

	return nil
}
