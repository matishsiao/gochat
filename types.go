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
