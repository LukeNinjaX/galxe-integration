package api

type TaskStatus string

var (
	TaskStatusNew        TaskStatus = "0"
	TaskStatusPending    TaskStatus = "1"
	TaskStatusProcessing TaskStatus = "2" // manual handling
	TaskStatusSuccess    TaskStatus = "3"
	TaskStatusFail       TaskStatus = "4"
)

const (
	Task_Topic_Goplus = "goplus"
	Task_Topic_Sys    = "sys"
)

const (
	Task_Name_AddLiquidity = "AddLiquidity"
	Task_Name_AspectPull   = "AspectPull"
	Task_Name_RugPull      = "RugPull"
	Task_Name_GetFaucet    = "GetFaucet"
	Task_Name_Sync         = "Sync"
)
