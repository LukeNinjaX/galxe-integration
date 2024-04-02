package biz

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
)

type PostBody struct {
	ChannelCode   string `json:"channelCode"`
	ChannelTaskId string `json:"channelTaskId"`
	CompleteTime  string `json:"completeTime"`
	UserAddress   string `json:"userAddress"`
}
type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Status bool `json:"status"`
	} `json:"result"`
}

// init by main
var GoPlus_Config *config.GoPlusConfig

func SyncStatus(db *sql.DB, input *InitTaskQuery) error {
	compiled, err := CheckAllTaskCompiled(db, input.AccountAddress)
	if err != nil {
		return err
	}
	if compiled == false {
		return fmt.Errorf(" All tasks haven't been completed")
	}

	postBody := &PostBody{
		ChannelCode:   GoPlus_Config.ChannelCode,
		ChannelTaskId: input.TaskId,
		CompleteTime:  strconv.FormatInt(time.Now().Unix(), 10),
		UserAddress:   input.AccountAddress,
	}

	client := &http.Client{Timeout: time.Second * 20}
	// body

	bytesData, _ := json.Marshal(postBody)
	pstReq, postErr := http.NewRequest("POST", GoPlus_Config.SecwarexUrl, bytes.NewReader(bytesData))
	if postErr != nil {
		return postErr
	}

	// header
	pstReq.Header.Add("Content-Type", "application/json")
	pstReq.Header.Add("manageId", GoPlus_Config.ManageId)
	pstReq.Header.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	sign, s, postErr := createSign(postBody)
	if postErr != nil {
		return postErr
	}
	pstReq.Header.Add("sign", sign)
	log.Info("goplus|request|", GoPlus_Config.SecwarexUrl, "|body|", string(bytesData), "|sign|", sign, "|plaintext|", s)

	resp, doErr := client.Do(pstReq)
	if doErr != nil {
		return doErr
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	log.Info("goplus|response|", string(body), "|address|", input.AccountAddress)

	responseData := &ResponseData{}
	jErr := json.Unmarshal(body, responseData)
	if jErr != nil {
		return jErr
	}

	// update db
	if responseData.Result.Status == true {
		topic := types.Task_Topic_Sys
		status := string(types.TaskStatusSuccess)
		taskName := types.Task_Name_Sync
		result := string(body)
		updateTaskQuery := &UpdateTaskQuery{
			TaskTopic:      &topic,
			TaskStatus:     &status,
			Memo:           &result,
			TaskId:         &input.TaskId,
			AccountAddress: &input.AccountAddress,
			TaskName:       &taskName,
		}

		upErr := UpdateTask(db, updateTaskQuery)
		if upErr != nil {
			return upErr
		}
	}
	return nil
}

func createSign(body *PostBody) (string, string, error) {

	var queryBuilder strings.Builder
	queryBuilder.WriteString("channelCode")
	queryBuilder.WriteString(body.ChannelCode)

	queryBuilder.WriteString("channelTaskId")
	queryBuilder.WriteString(body.ChannelTaskId)

	queryBuilder.WriteString("completeTime")
	queryBuilder.WriteString(body.CompleteTime)

	queryBuilder.WriteString("manageKey")
	queryBuilder.WriteString(GoPlus_Config.ManageKey)

	queryBuilder.WriteString("timestamp")
	queryBuilder.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))

	queryBuilder.WriteString("userAddress")
	queryBuilder.WriteString(body.UserAddress)

	query := queryBuilder.String()
	data := []byte(queryBuilder.String())
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)

	return md5str, query, nil
}

// Check all task compiled
func CheckAllTaskCompiled(db *sql.DB, addr string) (bool, error) {
	tasks, getErr := GetTask(db, addr, types.Task_Name_Sync, 0)
	if getErr != nil {
		return false, getErr
	}
	if tasks.TaskStatus != nil && strings.EqualFold(*tasks.TaskStatus, string(types.TaskStatusSuccess)) {
		return false, fmt.Errorf("Sync task have been completed ", strconv.FormatInt(tasks.ID, 10), "  ", *tasks.TaskStatus, "  ")
	}
	// check that all four tasks have been completed；
	countSql := "select count(*) from address_tasks where LOWER(account_address)=LOWER($1) and task_status=$2 and task_topic=$3"
	rows, err := db.Query(countSql, addr, string(types.TaskStatusSuccess), types.Task_Topic_Goplus)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if countErr := rows.Scan(&count); countErr != nil {
			return false, countErr
		}
	}
	if count == 4 {
		return true, nil
	} else {
		return false, nil
	}
}
