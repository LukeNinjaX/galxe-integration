package biz

import "database/sql"

func GetFaucetTask(db *sql.DB, limit int) ([]AddressTask, error) {
	whereTaskName := "GetFaucet"
	return lockTasksForHandler(db, limit, whereTaskName)
}

func GetAspectPullTask(db *sql.DB, limit int) ([]AddressTask, error) {
	whereTaskName := "RugPull"
	return lockTasksForHandler(db, limit, whereTaskName)
}

func GetAddLiquidityTask(db *sql.DB, limit int) ([]AddressTask, error) {
	whereTaskName := "AddLiquidity"
	return lockTasksForHandler(db, limit, whereTaskName)
}

func lockTasksForHandler(db *sql.DB, limit int, whereTaskName string) ([]AddressTask, error) {
	SetStatus := "2"
	whereStatus := "1"
	updateQuery := &UpdateTaskQuery{
		TaskStatus:  &SetStatus,
		TaskName:    &whereTaskName,
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
