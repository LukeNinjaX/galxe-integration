package faucet

import "time"

const (
	PullInterval   = 1 * time.Second
	PullBatchCount = 20

	PushInterval   = 2 * time.Second
	PushBatchCount = 50

	QueueMaxSize = 200

	BlockTime          = 600 * time.Millisecond
	GetReceiptInterval = 100 * time.Millisecond

	TransferAmount int64 = 1
)
