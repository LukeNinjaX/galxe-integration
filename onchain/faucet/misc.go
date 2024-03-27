package faucet

import "time"

const (
	PullInterval   = 1 * time.Second
	PullBatchCount = 20

	PushInterval   = 2 * time.Second
	PushBatchCount = 50

	QueueMaxSize = 200

	BlockTime          = 1600 * time.Millisecond
	GetReceiptInterval = 100 * time.Millisecond
)
