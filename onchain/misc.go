package onchain

import "time"

const (
	PullSleep = 100 * time.Millisecond
	PushSleep = 10 * time.Millisecond

	DeQuequeWait = 100 * time.Millisecond

	PullBatchCount = 200
	QueueSize      = 1000

	Reconnect = 200 * time.Millisecond

	CleanDBInterval = 10 * time.Minute
)
