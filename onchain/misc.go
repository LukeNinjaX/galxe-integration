package onchain

import "time"

const (
	PullSleep = 100 * time.Millisecond
	PushSleep = 20 * time.Millisecond

	DeQuequeWait = 100 * time.Millisecond

	PullBatchCount = 100
	QueueSize      = 1000

	Reconnect = 200 * time.Millisecond

	CleanDBInterval = 10 * time.Minute
	MaxRetry        = 50
)
