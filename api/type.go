package api

type TaskStatus string

var (
	TaskStatusNew        TaskStatus = "0"
	TaskStatusPending    TaskStatus = "1"
	TaskStatusProcessing TaskStatus = "2" // manual handling
	TaskStatusSuccess    TaskStatus = "3"
	TaskStatusFail       TaskStatus = "4"
)
