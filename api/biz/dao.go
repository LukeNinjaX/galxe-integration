package biz

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UpdateTaskQuery struct {
	AccountAddress string `json:"address" xml:"address" binding:"required"`
	TaskName       string `json:"taskName" xml:"taskName" binding:"required"`
	TaskStatus     string `json:"taskStatus" xml:"taskStatus" binding:"required"`
	Memo           string `json:"memo" xml:"memo" binding:"required"`
	TxHash         string `json:"txHash" xml:"txHash" binding:"required"`
}

type AddressTask struct {
	ID             int64
	GMTCreate      time.Time
	GMTModify      time.Time
	AccountAddress string
	TaskName       string
	// 0 init, 1 front done, 2 blockchain  check
	TaskStatus string
	Memo       string
	TxHash     string
}
type TaskInfo struct {
	TaskName   string `json:"taskName,omitempty"`
	TaskStatus int8   `json:"taskStatus,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type AccountTaskInfo struct {
	AccountAddress string `json:"accountAddress,omitempty"`
	// 0:no task 1:part finish 2:completed
	Status    int8       `json:"status,omitempty"`
	CanSync   bool       `json:"canSync,omitempty"`
	TaskInfos []TaskInfo `json:"taskInfos,omitempty"`
}

func UpdateTask(db *sql.DB, query UpdateTaskQuery) error {
	// 生成 UPDATE 语句

	if query.AccountAddress == "" || query.TaskName == "" {
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
	args = append(args, query.AccountAddress)
	queryBuilder.WriteString(fmt.Sprintf("%d", len(args)+1))
	args = append(args, query.TaskName)

	// 去除末尾的逗号和空格
	updateSql := strings.TrimSuffix(queryBuilder.String(), ", ")

	// 执行 UPDATE 语句
	_, err := db.Exec(updateSql, args...)
	return err
}
func GetAccountTaskInfo(db *sql.DB, addr string) (AccountTaskInfo, error) {

	taskInfos, err := GetTasks(db, addr)
	if err != nil {
		return AccountTaskInfo{}, err
	}
	return AccountTaskInfo{
		AccountAddress: addr,
		Status:         calculateStatus(taskInfos),
		CanSync:        calculateSyncCondition(taskInfos),
		TaskInfos:      taskInfo(taskInfos),
	}, nil
}

func taskDescription(taskName string) TaskInfo {

	pull := TaskInfo{
		TaskName:    "RugPull",
		TaskStatus:  0,
		Title:       "Rug Pull",
		Description: "\"Rug Pull\" is a malicious act where liquidity providers on decentralized exchanges like Uniswap suddenly withdraw their provided liquidity, causing a sharp drop in liquidity for a trading pair. This action can severely affect users' ability to trade or result in financial losses. While Uniswap strives to mitigate this risk through community oversight and contract audits, Rug Pulls remain a concern, emphasizing the importance of cautious participation and thorough due diligence on projects and teams before engaging with liquidity pools.",
	}
	aspect := TaskInfo{
		TaskName:    "AspectPull",
		TaskStatus:  0,
		Title:       "Aspect Work",
		Description: "Aspect is a tool for detecting and blocking suspicious transactions, ideal for preventing Rug Pulls. It operates by continuously monitoring transaction behavior in real-time, identifying potential risk patterns and anomalies. For instance, Rug Pulls often involve sudden, substantial fund withdrawals or transfers. Aspect tracks abrupt spikes in transaction volume, comparing them with historical data. Upon detecting unusual changes, Aspect issues alerts and reviews transactions to confirm any risks.",
	}
	addLiquidity := TaskInfo{
		TaskName:    "AddLiquidity",
		TaskStatus:  0,
		Title:       "Add Liquidity",
		Description: "When users add liquidity to a pool, they typically provide an equal value of two different tokens in the pair. For example, in the ETH/DAI trading pair, a user might add an equal value of Ethereum and DAI tokens to the liquidity pool. In return for providing liquidity, users receive liquidity tokens representing their share of the pool. These tokens can be used to withdraw their portion of the liquidity at any time, along with any accumulated trading fees.",
	}

	// 把上面3个taskInfo加入到一个map中，key是TaskName

	taskMap := map[string]TaskInfo{
		pull.TaskName:         pull,
		aspect.TaskName:       aspect,
		addLiquidity.TaskName: addLiquidity,
	}
	return taskMap[taskName]

}
func taskInfo(tasks []AddressTask) []TaskInfo {
	var taskInfos []TaskInfo
	for _, task := range tasks {
		// 将字符串转换为int64类型
		num, err := strconv.ParseInt(task.TaskStatus, 10, 8)
		if err != nil {
			fmt.Println("Failed to ParseInt:", err)
		}
		// 将int64类型转换为int8类型
		intValue := int8(num)

		description := taskDescription(task.TaskName)

		taskInfos = append(taskInfos, TaskInfo{
			TaskName:    task.TaskName,
			TaskStatus:  intValue,
			Title:       description.Title,
			Description: description.Description,
		})
	}
	return taskInfos

}

func calculateSyncCondition(tasks []AddressTask) bool {
	status := calculateStatus(tasks)
	// 2:completed
	return status == 2
}
func calculateStatus(tasks []AddressTask) int8 {
	status := 0
	if len(tasks) == 0 {
		return int8(status)
	}
	DoneStatus := "1"
	count := 0
	for _, task := range tasks {
		if task.TaskStatus == DoneStatus {
			count += 1
		}
	}
	if count == 3 {
		// 2:completed
		status = 2
	} else if count > 0 && count < 3 {
		// 1:part finish
		status = 1
	}
	return int8(status)
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

func GetTask(db *sql.DB, addr string, taskName string) (AddressTask, error) {
	var addressTask AddressTask
	err := db.QueryRow("SELECT id,gmt_create,gmt_modify,account_address,task_name,task_status,memo,tx_hash  FROM address_tasks WHERE account_address = $1 and task_name = $2", addr, taskName).Scan(
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
		return addressTask, err
	}
	return addressTask, nil
}
