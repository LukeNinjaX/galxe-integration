package slack

type Config struct {
	Timeout  string `json:"timeout"`
	Token    string `json:"token"`
	Channel  string `json:"channel"`
	Throttle string `json:"throttle"`
}
