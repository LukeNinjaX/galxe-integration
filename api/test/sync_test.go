package biz

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestInitData(t *testing.T) {
	// open file
	f, err := os.Open("_address.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	postLog, logErr := os.Create("post_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt")
	if logErr != nil {
		log.Fatal(logErr)
		return
	}
	defer func(postLog *os.File) {
		_ = postLog.Close()
	}(postLog)

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)
	taskId := "891b8fbef81c43c7aec3e4bfeea2c752"
	for scanner.Scan() {
		// do something with a line
		text := scanner.Text()
		splitResult := strings.Split(text, ": ")
		s, postErr := posts(splitResult[1], taskId)
		log.Info(s, postErr)

		_, _ = fmt.Fprintln(postLog, text, s)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
func posts(address string, taskId string) (string, error) {
	mapInstance := make(map[string]interface{})
	mapInstance["accountAddress"] = address
	mapInstance["taskId"] = taskId
	jsonData, err := json.Marshal(mapInstance)
	if err != nil {
		return "", err
	}
	taskUrl := "https://campaign.artela.network/api/goplus/new-task"
	req, postErr := http.NewRequest("POST", taskUrl, bytes.NewBuffer(jsonData))
	if postErr != nil {
		return "", postErr
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}
	return string(body), nil
}
