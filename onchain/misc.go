package onchain

import (
	"errors"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/ethereum/go-ethereum/common"
)

type (
	Send           func(task biz.AddressTask) (hashs []common.Hash, err error)
	GetTasks       func(count int) ([]biz.AddressTask, error)
	UpdateTask     func(task biz.AddressTask, hashs []common.Hash, status *uint64) error
	RefreshNetwork func() bool
)

const (
	UINT      = 1000000000000000000
	RugAmount = 10000000

	PullSleep = 100 * time.Millisecond
	PushSleep = 50 * time.Millisecond

	DeQuequeWait = 100 * time.Millisecond

	PullBatchCount = 100
	QueueSize      = 1000

	Reconnect = 200 * time.Millisecond

	CleanDBInterval = 10 * time.Minute
	MaxRetry        = 50
)

var (
	ErrInvalidTask = errors.New("task is not valid")
)
