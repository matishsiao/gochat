package main

type Configs struct {
	Debug     bool   `json:"debug"`
	HTTP      string `json:"http"`
	HTTPS     string `json:"https"`
	Mode      string `json:"mode"`
	AuthToken string `json:"authtoken"`
	SSL       struct {
		Key string `json:"key"`
		Crt string `json:"crt"`
	} `json:"ssl"`
	Timeout         int64 `json:"timeout"`
	ConnectionLimit int   `json:"connection_limit"`
}

type ChatMessage struct {
	User      string   `json:"user"`
	UUID      string   `json:"uuid"`
	Timestamp int64    `json:"timestamp"`
	Message   string   `json:"message"`
	Type      string   `json:"type"`
	Channel   string   `json:"channel"`
	Token     string   `json:"token"`
	Data      ChatData `json:"data"`
}

type ChatData struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}
