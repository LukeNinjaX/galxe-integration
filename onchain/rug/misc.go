package rug

import "time"

const (
	PullInterval   = 1 * time.Second
	PullBatchCount = 20

	PushInterval   = 2 * time.Second
	PushBatchCount = 10

	QueueMaxSize = 200
)

var (
	TaskStatusFail    = "fail"
	TaskStatusSuccess = "success"
)
