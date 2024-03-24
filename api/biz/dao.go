package biz

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	Task_Group_Normal = "normal"
	Task_Group_Sys    = "sys"
)

type UpdateTaskQuery struct {
	// update set
	TaskGroup     *string `json:"taskGroup" xml:"taskGroup" `
	TaskStatus    *string `json:"taskStatus" xml:"taskStatus"`
	Memo          *string `json:"memo" xml:"memo" `
	Txs           *string `json:"txs" xml:"txs" `
	ChannelTaskId *string `json:"channelTaskId" xml:"channelTaskId"`

	// where condition
	ID             int64   `json:"id" xml:"id" binding:"required"`
	AccountAddress *string `json:"accountAddress" xml:"address" binding:"required"`
	TaskName       *string `json:"taskName" xml:"taskName" `
}
type InitTaskQuery struct {
	AccountAddress string `json:"accountAddress" xml:"address" binding:"required"`
	ChannelTaskId  string `json:"channelTaskId" xml:"channelTaskId" binding:"required"`
}
type TaskQuery struct {
	AccountAddress string `json:"accountAddress" xml:"address" binding:"required"`
	ChannelTaskId  string `json:"channelTaskId" xml:"channelTaskId" binding:"required"`
	TaskStatus     string `json:"taskStatus" xml:"taskStatus" binding:"required"`
	TaskGroup      string `json:"taskGroup" xml:"taskGroup" binding:"required"`
	TaskName       string `json:"taskName" xml:"taskName" binding:"required"`
}

type AddressTask struct {
	ID             int64
	GMTCreate      time.Time
	GMTModify      time.Time
	AccountAddress *string `db:"account_address"`
	TaskName       *string `db:"task_name"`
	// 0 init, 1 front done, 2 blockchain  check
	TaskStatus    *string `db:"task_status"`
	Memo          *string `db:"memo"`
	Txs           *string `db:"txs"`
	ChannelTaskId *string `db:"channel_task_id"`
	TaskGroup     *string `db:"task_group"`
}
type TaskInfo struct {
	ID         int64  `json:"id,omitempty"`
	TaskName   string `json:"taskName,omitempty"`
	TaskStatus int8   `json:"taskStatus"`

	Title string `json:"title"`
}

type AccountTaskInfo struct {
	AccountAddress string `json:"accountAddress"`
	// 0:no task 1:part finish 2:completed
	Status    int8       `json:"status"`
	CanSync   bool       `json:"canSync"`
	TaskInfos []TaskInfo `json:"taskInfos,omitempty"`
}

func InitTask(db *sql.DB, query *InitTaskQuery) error {
	if query.AccountAddress == "" || query.ChannelTaskId == "" {
		return fmt.Errorf("address or ChannelTaskId cannot be empty")
	}
	// insert db
	insertSql := "INSERT INTO address_tasks (account_address, task_name,task_status,channel_task_id,task_group) VALUES " +
		"($1, 'AddLiquidity','0',$2,$3)," +
		"($1, 'AspectPull', '0',$2,$3)," +
		"($1, 'RugPull', '0',$2,$3)," +
		"($1, 'GetFaucet', '0',$2,$3)," +
		"($1, 'Sync', '0',$2,$4);"

	_, err := db.Exec(insertSql, query.AccountAddress, query.ChannelTaskId, Task_Group_Normal, Task_Group_Sys)
	if err != nil {
		return err
	}
	return nil
}
func UpdateTask(db *sql.DB, query UpdateTaskQuery) error {
	// 生成 UPDATE 语句

	if query.AccountAddress == nil || query.ID == 0 {
		return fmt.Errorf("address or id cannot be empty")
	}

	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString("UPDATE address_tasks SET ")

	if query.TaskStatus != nil {
		queryBuilder.WriteString("task_status = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.TaskStatus)
	}
	if query.Memo != nil {
		queryBuilder.WriteString("memo = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.Memo)
	}
	if query.Txs != nil {
		queryBuilder.WriteString("txs = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.Txs)
	}
	if query.ChannelTaskId != nil {
		queryBuilder.WriteString("channel_task_id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.ChannelTaskId)
	}
	queryBuilder.WriteString(" gmt_modify = CURRENT_TIMESTAMP ")

	queryBuilder.WriteString(" WHERE 1=1 ")
	if query.ID > 0 {
		queryBuilder.WriteString(" and id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.ID)
	}
	if query.AccountAddress != nil {
		queryBuilder.WriteString(" and account_address = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.AccountAddress)
	}
	if query.TaskName != nil {
		queryBuilder.WriteString(" and task_name = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskName)
	}
	// 去除末尾的逗号和空格
	querySql := strings.TrimSuffix(queryBuilder.String(), ", ")

	// 执行 UPDATE 语句
	_, err := db.Exec(querySql, args...)
	return err
}
func GetAccountTaskInfo(db *sql.DB, addr string) (AccountTaskInfo, error) {

	taskInfos, err := GetTasks(db, &TaskQuery{
		AccountAddress: addr,
		TaskGroup:      Task_Group_Normal,
	})
	if err != nil {
		return AccountTaskInfo{}, err
	}
	if len(taskInfos) == 0 {
		return AccountTaskInfo{
			AccountAddress: addr,
			Status:         0,
			CanSync:        false,
		}, nil
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
		TaskName:   "RugPull",
		TaskStatus: 0,
		Title:      "Rug Pull",
	}
	aspect := TaskInfo{
		TaskName:   "AspectPull",
		TaskStatus: 0,
		Title:      "Aspect Work",
	}
	addLiquidity := TaskInfo{
		TaskName:   "AddLiquidity",
		TaskStatus: 0,
		Title:      "Add Liquidity",
	}
	getFaucet := TaskInfo{
		TaskName:   "GetFaucet",
		TaskStatus: 0,
		Title:      "Get Faucet",
	}
	// 把上面3个taskInfo加入到一个map中，key是TaskName

	taskMap := map[string]TaskInfo{
		pull.TaskName:         pull,
		aspect.TaskName:       aspect,
		addLiquidity.TaskName: addLiquidity,
		getFaucet.TaskName:    getFaucet,
	}
	return taskMap[taskName]

}
func taskInfo(tasks []AddressTask) []TaskInfo {
	var taskInfos []TaskInfo
	for _, task := range tasks {
		// 将字符串转换为int64类型
		num, err := strconv.ParseInt(*task.TaskStatus, 10, 8)
		if err != nil {
			fmt.Println("Failed to ParseInt:", err)
		}
		// 将int64类型转换为int8类型
		intValue := int8(num)

		description := taskDescription(*task.TaskName)

		taskInfos = append(taskInfos, TaskInfo{
			ID:         task.ID,
			TaskName:   *task.TaskName,
			TaskStatus: intValue,
			Title:      description.Title,
		})
	}
	return taskInfos

}

// 是否
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
		if *task.TaskStatus == DoneStatus {
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

func GetTasks(db *sql.DB, query *TaskQuery) ([]AddressTask, error) {
	// 遍历结果集
	var addressTasks []AddressTask
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString("SELECT id,gmt_create,gmt_modify,account_address,task_name,task_status,memo,txs,channel_task_id,task_group FROM address_tasks ")

	if query.AccountAddress == "" {
		return addressTasks, fmt.Errorf("address cannot be empty")
	}
	queryBuilder.WriteString(" WHERE account_address = $1 ")
	args = append(args, query.AccountAddress)

	if query.ChannelTaskId != "" {
		queryBuilder.WriteString(" and channel_task_id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.ChannelTaskId)
	}
	if query.TaskStatus != "" {
		queryBuilder.WriteString(" and task_status = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskStatus)
	}
	if query.TaskGroup != "" {
		queryBuilder.WriteString(" and task_group = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskGroup)
	}
	if query.TaskName != "" {
		queryBuilder.WriteString(" and task_name = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskName)
	}
	// 去除末尾的逗号和空格
	querySql := strings.TrimSuffix(queryBuilder.String(), ", ")

	rows, err := db.Query(querySql, args...)
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
			&addressTask.Txs,
			&addressTask.ChannelTaskId,
			&addressTask.TaskGroup,
		)
		if err != nil {
			log.Fatal(err)
		}
		addressTasks = append(addressTasks, addressTask)
	}
	return addressTasks, nil
}

func GetTask(db *sql.DB, addr string, taskName string) (AddressTask, error) {
	addressTask := AddressTask{}
	tasks, err := GetTasks(db, &TaskQuery{
		AccountAddress: addr,
		TaskName:       taskName,
	})
	if err != nil || len(tasks) == 0 {
		return addressTask, err
	}
	return tasks[0], nil
}
