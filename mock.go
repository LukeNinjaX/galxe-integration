package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/ethereum/go-ethereum/common"
)

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
	memo := strings.Join(hexs, ",")
	if status != nil {
		fmt.Println("updating task ", task.ID, memo, *status)
	} else {
		fmt.Println("updating task ", task.ID, memo)
	}

	return nil
}
