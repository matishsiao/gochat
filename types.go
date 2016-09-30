package main

type Configs struct {
	Debug bool   `json:"debug"`
	HTTP  string `json:"http"`
	HTTPS string `json:"https"`
	SSL   struct {
		Key string `json:"key"`
		Crt string `json:"crt"`
	} `json:"ssl"`
	Timeout         int64 `json:"timeout"`
	ConnectionLimit int   `json:"connection_limit"`
}

type ChatMessage struct {
	User      string `json:"user"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	Token     string `json:"token"`
}
