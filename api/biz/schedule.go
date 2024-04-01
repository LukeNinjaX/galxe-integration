package biz

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/artela-network/galxe-integration/api/types"
)

func GetFaucetTask(db *sql.DB, limit int) ([]AddressTask, error) {
	return lockTasksForHandler(db, limit, types.Task_Name_GetFaucet)
}

func GetAspectPullTask(db *sql.DB, limit int) ([]AddressTask, error) {
	return lockTasksForHandler(db, limit, types.Task_Name_AspectPull)
}

func GetAddLiquidityTask(db *sql.DB, limit int) ([]AddressTask, error) {
	return lockTasksForHandler(db, limit, types.Task_Name_AddLiquidity)
}

func lockTasksForHandler(db *sql.DB, limit int, whereTaskName string) ([]AddressTask, error) {

	SetStatus := string(types.TaskStatusProcessing)
	whereStatus := string(types.TaskStatusPending)
	uuidV4 := uuid.New().String()

	if limit == 0 && whereTaskName == "" {
		return nil, fmt.Errorf("limit or whereTaskName cannot be empty")
	}

	limitSql := " LIMIT $5) "
	if strings.EqualFold(whereTaskName, types.Task_Name_AddLiquidity) {
		limitSql = "and txs IS NOT NULL LIMIT $5) "
	}
	querySql := "UPDATE address_tasks SET task_status = $1, job_batch_id = $2, gmt_modify = CURRENT_TIMESTAMP WHERE id IN (SELECT id FROM address_tasks WHERE task_name=$3 and task_status = $4 " + limitSql

	// 执行 UPDATE 语句
	_, err := db.Exec(querySql, SetStatus, uuidV4, whereTaskName, whereStatus, limit)
	if err != nil {
		return nil, err
	}
	// get tasks for schedule
	query := &TaskQuery{
		TaskName:   whereTaskName,
		TaskStatus: SetStatus,
		JobBatchId: uuidV4,
	}
	return GetTasks(db, query)
}

// let timeout data retry
func LetTimeoutRecordRetry(db *sql.DB) (int64, error) {
	selectSql := "UPDATE address_tasks SET task_status = '1', job_batch_id = null , gmt_modify = CURRENT_TIMESTAMP where gmt_modify < current_timestamp - interval '10 minutes' and (task_status='2')"
	res, err := db.Exec(selectSql)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected, nil
}
