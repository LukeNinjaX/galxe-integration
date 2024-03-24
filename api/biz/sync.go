package biz

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const SECWAREX_URL = "https://test-api.secwarex.io/open/v1/task/finish"
const MANAGE_ID = "APP10001"
const CHANNEL_CODE = "APP10001"

// const SECWAREX_URL = "https://api.secwarex.io/open/v1/task/finish"

func SyncStatus(db *sql.DB, addr string) error {
	client := &http.Client{Timeout: time.Second * 20}

	// body

	postBody := getbody(db, addr)
	bytesData, _ := json.Marshal(postBody)
	pstReq, postErr := http.NewRequest("POST", SECWAREX_URL, bytes.NewReader(bytesData))
	if postErr != nil {
		return postErr
	}

	// header
	pstReq.Header.Add("Content-Type", "application/json")
	pstReq.Header.Add("manageId", MANAGE_ID)
	pstReq.Header.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	pstReq.Header.Add("sign", createSign(db, addr))

	resp, doErr := client.Do(pstReq)
	if doErr != nil {
		return doErr
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	fmt.Println(string(body))

	return nil
}

func createSign(db *sql.DB, addr string) string {
	return "xx"
}

type PostBody struct {
	ChannelCode   string `json:"channelCode"`
	ChannelTaskId string `json:"channelTaskId"`
	CompleteTime  string `json:"completeTime"`
	UserAddress   string `json:"userAddress"`
}

func getbody(db *sql.DB, addr string) *PostBody {
	body := &PostBody{}
	return body
}
