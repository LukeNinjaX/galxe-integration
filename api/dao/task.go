package dao

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UpdateTaskQuery struct {
	Account    string `json:"address" xml:"address" binding:"required"`
	TaskName   string `json:"taskName" xml:"taskName" binding:"required"`
	TaskStatus string `json:"taskStatus" xml:"taskStatus" binding:"required"`
	Memo       string `json:"memo" xml:"memo" binding:"required"`
	TxHash     string `json:"txHash" xml:"txHash" binding:"required"`
}

type AddressTask struct {
	ID             int64
	GMTCreate      time.Time
	GMTModify      time.Time
	AccountAddress string
	TaskName       string
	TaskStatus     string
	Memo           string
	TxHash         string
}

func UpdateTask(db *sql.DB, query UpdateTaskQuery) error {
	// 生成 UPDATE 语句

	if query.Account == "" || query.TaskName == "" {
		return fmt.Errorf("address or taskName cannot be empty")
	}

	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString("UPDATE address_tasks SET ")
	if query.TaskStatus != "" {
		queryBuilder.WriteString("task_status = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.TaskStatus)
	}
	if query.Memo != "" {
		queryBuilder.WriteString("memo = memo || $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.Memo)
	}
	if query.TxHash != "" {
		queryBuilder.WriteString("tx_hash = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.TxHash)

	}

	// 去除末尾的逗号和空格
	queryBuilder.WriteString(" WHERE account_address = $ and task_name = $")
	queryBuilder.WriteString(fmt.Sprintf("%d", len(args)+1))
	args = append(args, query.Account)
	queryBuilder.WriteString(fmt.Sprintf("%d", len(args)+1))
	args = append(args, query.TaskName)

	// 去除末尾的逗号和空格
	updateSql := strings.TrimSuffix(queryBuilder.String(), ", ")

	// 执行 UPDATE 语句
	_, err := db.Exec(updateSql, args...)
	return err
}

func GetTasks(db *sql.DB, addr string) ([]AddressTask, error) {
	rows, err := db.Query("SELECT id,gmt_create,gmt_modify,account_address,task_name,task_status,memo,tx_hash  FROM address_tasks WHERE account_address = $1", addr)
	if err != nil {
		log.Errorf("Failed to getTasks: %v", err)
		return nil, err
	}
	// 解析 query 到 struct 类型中
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Errorf("Failed to close sql at getTasks: %v", closeErr)
		}
	}(rows)

	// 遍历结果集
	var addressTasks []AddressTask
	for rows.Next() {
		var addressTask AddressTask
		err := rows.Scan(
			&addressTask.ID,
			&addressTask.GMTCreate,
			&addressTask.GMTModify,
			&addressTask.AccountAddress,
			&addressTask.TaskName,
			&addressTask.TaskStatus,
			&addressTask.Memo,
			&addressTask.TxHash,
		)
		if err != nil {
			log.Fatal(err)
		}
		addressTasks = append(addressTasks, addressTask)
	}
	return addressTasks, nil
}
