package biz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/artela-network/galxe-integration/config"
)

// init by main
var Recaptcha_Config *config.RecaptchaConfig

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func CheckCaptcha(responseToken string) error {

	data := url.Values{"secret": {Recaptcha_Config.Secret}, "response": {responseToken}}
	resp, err := http.PostForm(Recaptcha_Config.VerifyUrl, data)

	if err != nil {
		return err
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	log.Info("Recaptcha response:  ", string(body))

	responseData := &RecaptchaResponse{}
	jErr := json.Unmarshal(body, responseData)
	if jErr != nil {
		return jErr
	}
	if !responseData.Success {
		return fmt.Errorf("recaptcha failed: %v", responseData.ErrorCodes)
	}
	return nil
}
