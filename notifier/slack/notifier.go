package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/artela-network/galxe-integration/common"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

const NotifierName = "slack"

type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type Notifier struct {
	token    string
	throttle time.Duration
	channel  string

	httpClient        *http.Client
	ctx               context.Context
	lastNotifiedTimes sync.Map
}

func NewSlackNotifier(ctx context.Context, rawConf json.RawMessage) common.Notifier {
	conf := &Config{}
	if err := json.Unmarshal(rawConf, conf); err != nil {
		log.Panic("failed to parse wechat notifier config", err)
	}

	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		timeout = 5 * time.Second
	}

	throttle, err := time.ParseDuration(conf.Throttle)
	if err != nil {
		throttle = 5 * time.Second
	}

	return &Notifier{
		token:   conf.Token,
		channel: conf.Channel,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		throttle: throttle,
		ctx:      ctx,
	}
}

func (n *Notifier) Notify(msg, from string, throttle bool) {
	log.Info("sending slack notification...")
	log.Debugf("slack payload: %v", msg)

	if throttle {
		currentTime := time.Now().UnixNano()
		lastNotifiedTime, ok := n.lastNotifiedTimes.Load(from)
		if ok && lastNotifiedTime.(int64)+int64(n.throttle) > currentTime {
			log.Infof("slack notification throttled, ignore message")
			return
		}
		n.lastNotifiedTimes.Store(from, currentTime)
	}

	payload := SlackMessage{
		Channel: n.channel,
		Text:    msg,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("failed to marshal slack message: %v", err)
		return
	}

	request, err := http.NewRequestWithContext(n.ctx, http.MethodPost, "https://slack.com/api/chat.postMessage", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Errorf("failed to create request: %v", err)
		return
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+n.token)

	resp, err := n.httpClient.Do(request)
	if err != nil {
		log.Errorf("failed to send notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMsg, _ := io.ReadAll(resp.Body)
		log.Errorf("failed to send notification: %s", errorMsg)
	}
	log.Info("successfully notified via slack")
}
