package biz

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/artela-network/galxe-integration/api/types"
)

type UpdateTaskQuery struct {
	// update set
	TaskTopic  *string `json:"taskTopic" xml:"taskTopic" `
	TaskStatus *string `json:"taskStatus" xml:"taskStatus"`
	Memo       *string `json:"memo" xml:"memo" `
	Txs        *string `json:"txs" xml:"txs" `
	TaskId     *string `json:"taskId" xml:"taskId"`
	JobBatchId *string `json:"jobBatchId" xml:"jobBatchId"`

	// where condition
	ID             int64   `json:"id" xml:"id"`
	AccountAddress *string `json:"accountAddress" xml:"address"`
	TaskName       *string `json:"taskName" xml:"taskName" `
	StatusEqual    *string `json:"statusEqual" xml:"statusEqual"`
}
type InitTaskQuery struct {
	AccountAddress string `json:"accountAddress" xml:"address" binding:"required"`
	TaskId         string `json:"taskId" xml:"taskId"`
	TaskTopic      string `json:"taskTopic" xml:"taskTopic"`
}
type TaskQuery struct {
	ID             int64  `json:"id" xml:"id" binding:"required"`
	AccountAddress string `json:"accountAddress" xml:"address" binding:"required"`
	TaskId         string `json:"taskId" xml:"taskId"`
	TaskStatus     string `json:"taskStatus" xml:"taskStatus" `
	TaskTopic      string `json:"taskTopic" xml:"taskTopic"`
	TaskName       string `json:"taskName" xml:"taskName"`
	JobBatchId     string `json:"jobBatchId" xml:"jobBatchId"`
	LimitNum       int    `json:"limitNum" xml:"limitNum"`
}

type AddressTask struct {
	ID             int64
	GMTCreate      time.Time
	GMTModify      time.Time
	AccountAddress *string `db:"account_address"`
	TaskName       *string `db:"task_name"`
	// 0 init, 1 front done, 2 blockchain  check
	TaskStatus *string `db:"task_status"`
	Memo       *string `db:"memo"`
	Txs        *string `db:"txs"`
	TaskId     *string `db:"task_id"`
	TaskTopic  *string `db:"task_topic"`
	JobBatchId *string `db:"job_batch_id"`
}
type TaskInfo struct {
	ID         int64  `json:"id,omitempty"`
	TaskName   string `json:"taskName,omitempty"`
	TaskStatus int8   `json:"taskStatus"`

	Title string `json:"title"`
	Memo  string `json:"memo"`
	Txs   string `json:"txs"`
}

type AccountTaskInfo struct {
	AccountAddress string `json:"accountAddress"`
	// 0:no task 1:part finish 2:completed
	Status    int8       `json:"status"`
	CanSync   bool       `json:"canSync"`
	TaskInfos []TaskInfo `json:"taskInfos,omitempty"`
}

func InitTask(db *sql.DB, query *InitTaskQuery) error {
	if query.AccountAddress == "" || query.TaskId == "" {
		return fmt.Errorf("address or TaskId cannot be empty")
	}
	// insert db
	insertSql := "INSERT INTO address_tasks (account_address, task_name,task_status,task_id,task_topic) VALUES " +
		"($1, $5, '0',$2,$3)," +
		"($1, $6, '0',$2,$3)," +
		"($1, $7, '0',$2,$3)," +
		"($1, $8, '0',$2,$3)," +
		"($1, $9, '0',$2,$4);"

	_, err := db.Exec(insertSql,
		query.AccountAddress,
		query.TaskId,
		types.Task_Topic_Goplus,
		types.Task_Topic_Sys,
		types.Task_Name_GetFaucet,
		types.Task_Name_AddLiquidity,
		types.Task_Name_RugPull,
		types.Task_Name_AspectPull,
		types.Task_Name_Sync,
	)
	if err != nil {
		return err
	}
	return nil
}
func UpdateTask(db *sql.DB, query *UpdateTaskQuery) error {
	// 生成 UPDATE 语句

	if query.AccountAddress == nil && query.ID == 0 {
		return fmt.Errorf("address or id cannot be empty")
	}
	if query.ID != 0 && query.TaskStatus != nil {
		task, err2 := GetTask(db, "", "", query.ID)
		if err2 != nil {
			return fmt.Errorf("Task not found.")
		}
		if query.TaskStatus != nil && task.TaskStatus != nil {
			// 更新到1 ，原有状态必须是 0,或 4
			if *query.TaskStatus == "1" && (!(*task.TaskStatus == "0" || *task.TaskStatus == "4")) {
				return fmt.Errorf("when update  status to 1 Task status is not 0 or 4.")
			}
			if *query.TaskStatus == "2" && (*task.TaskStatus != "1") {
				return fmt.Errorf("when update status to 2 Task status is not 1.")
			}
			if *query.TaskStatus == "3" && (!(*task.TaskStatus == "1" || *task.TaskStatus == "2")) {
				return fmt.Errorf("when update status  to 3  Task status is not 2.")
			}
			if *query.TaskStatus == "4" && (*task.TaskStatus != "2") {
				return fmt.Errorf("when update status to 4 Task status is not 2.")
			}
		}
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
	if query.TaskId != nil {
		queryBuilder.WriteString("task_id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.TaskId)
	}
	if query.JobBatchId != nil {
		queryBuilder.WriteString("job_batch_id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d, ", len(args)+1))
		args = append(args, query.JobBatchId)
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
	if query.StatusEqual != nil {
		queryBuilder.WriteString(" and task_status = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.StatusEqual)
	}
	// 去除末尾的逗号和空格
	querySql := strings.TrimSuffix(queryBuilder.String(), ", ")

	// 执行 UPDATE 语句
	_, err := db.Exec(querySql, args...)
	return err
}
func GetAccountTaskInfo(db *sql.DB, query *TaskQuery) (AccountTaskInfo, error) {
	if db == nil || query == nil {
		return AccountTaskInfo{}, fmt.Errorf("address cannot be empty")
	}

	taskInfos, err := GetTasks(db, query)
	if err != nil {
		return AccountTaskInfo{}, err
	}
	if len(taskInfos) == 0 {
		return AccountTaskInfo{
			AccountAddress: query.AccountAddress,
			Status:         0,
			CanSync:        false,
		}, nil
	}
	return AccountTaskInfo{
		AccountAddress: query.AccountAddress,
		Status:         calculateStatus(taskInfos),
		CanSync:        calculateSyncCondition(taskInfos),
		TaskInfos:      ConvertTaskInfo(taskInfos),
	}, nil
}

func taskDescription(taskName string) TaskInfo {

	pull := TaskInfo{
		TaskName:   types.Task_Name_RugPull,
		TaskStatus: 0,
		Title:      "Rug Pull",
	}
	aspect := TaskInfo{
		TaskName:   types.Task_Name_AspectPull,
		TaskStatus: 0,
		Title:      "Aspect Work",
	}
	addLiquidity := TaskInfo{
		TaskName:   types.Task_Name_AddLiquidity,
		TaskStatus: 0,
		Title:      "Add Liquidity",
	}
	getFaucet := TaskInfo{
		TaskName:   types.Task_Name_GetFaucet,
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
func ConvertTaskInfo(tasks []AddressTask) []TaskInfo {
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

		taskItem := TaskInfo{
			ID:         task.ID,
			TaskStatus: intValue,
			Title:      description.Title,
		}
		if task.TaskName != nil {
			taskItem.TaskName = *task.TaskName
		}
		if task.Memo != nil {
			taskItem.Memo = *task.Memo
		}
		if task.Txs != nil {
			taskItem.Txs = *task.Txs
		}
		taskInfos = append(taskInfos, taskItem)

	}
	// sort.SliceStable(taskInfos, func(i, j int) bool {
	// 	return taskInfos[i].ID < taskInfos[j].ID
	// })
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

	queryBuilder.WriteString("SELECT id,gmt_create,gmt_modify,account_address,task_name,task_status,memo,txs,task_id,task_topic,job_batch_id FROM address_tasks ")

	queryBuilder.WriteString(" WHERE 1=1 ")
	if query.ID > 0 {
		queryBuilder.WriteString(" and id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.ID)
	}
	if query.AccountAddress != "" {
		queryBuilder.WriteString(" and account_address = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.AccountAddress)
	}

	if query.TaskId != "" {
		queryBuilder.WriteString(" and task_id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskId)
	}
	if query.TaskStatus != "" {
		queryBuilder.WriteString(" and task_status = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskStatus)
	}
	if query.TaskTopic != "" {
		queryBuilder.WriteString(" and task_topic = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskTopic)
	}
	if query.TaskName != "" {
		queryBuilder.WriteString(" and task_name = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.TaskName)
	}
	if query.JobBatchId != "" {
		queryBuilder.WriteString(" and job_batch_id = $")
		queryBuilder.WriteString(fmt.Sprintf("%d ", len(args)+1))
		args = append(args, query.JobBatchId)
	}
	// 去除末尾的逗号和空格
	querySql := strings.TrimSuffix(queryBuilder.String(), ", ")

	querySql = querySql + " ORDER BY ID ASC "
	if query.LimitNum > 0 {
		querySql = querySql + fmt.Sprintf(" LIMIT %d", query.LimitNum)
	}

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
			&addressTask.TaskId,
			&addressTask.TaskTopic,
			&addressTask.JobBatchId,
		)
		if err != nil {
			log.Fatal(err)
		}
		addressTasks = append(addressTasks, addressTask)
	}
	return addressTasks, nil
}

func GetTask(db *sql.DB, addr string, taskName string, id int64) (AddressTask, error) {
	addressTask := AddressTask{}
	tasks, err := GetTasks(db, &TaskQuery{
		AccountAddress: addr,
		TaskName:       taskName,
		ID:             id,
	})
	if err != nil || len(tasks) == 0 {
		return addressTask, err
	}
	return tasks[0], nil
}
