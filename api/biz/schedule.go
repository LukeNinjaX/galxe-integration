package biz

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/artela-network/galxe-integration/api"
)

func GetFaucetTask(db *sql.DB, limit int) ([]AddressTask, error) {
	return lockTasksForHandler(db, limit, api.Task_Name_GetFaucet)
}

func GetAspectPullTask(db *sql.DB, limit int) ([]AddressTask, error) {
	return lockTasksForHandler(db, limit, api.Task_Name_RugPull)
}

func GetAddLiquidityTask(db *sql.DB, limit int) ([]AddressTask, error) {
	return lockTasksForHandler(db, limit, api.Task_Name_AddLiquidity)
}

func lockTasksForHandler(db *sql.DB, limit int, whereTaskName string) ([]AddressTask, error) {
	SetStatus := string(api.TaskStatusProcessing)
	whereStatus := string(api.TaskStatusPending)
	uuidV4 := uuid.New().String()
	updateQuery := &UpdateTaskQuery{
		TaskStatus:  &SetStatus,
		TaskName:    &whereTaskName,
		JobBatchId:  &uuidV4,
		StatusEqual: &whereStatus,
		LimitNum:    limit}
	err := UpdateTask(db, updateQuery)
	if err != nil {
		return nil, err
	}
	// get tasks for schedule
	query := &TaskQuery{
		TaskName:   whereTaskName,
		TaskStatus: SetStatus,
	}
	return GetTasks(db, query)
}
