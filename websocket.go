package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "strconv"
	"time"

	"github.com/matishsiao/net/websocket"
)

var clientList *list.List

type WsClient struct {
	UId      string
	UserName string
	Channel  []string
	WS       *websocket.Conn
	Closed   bool
}

func (wc *WsClient) Close() {
	wc.Closed = true
	wc.WS.Close()
}

func Echo(ws *websocket.Conn) {
	var err error
	client := new(WsClient)
	client.UId = fmt.Sprintf("U-%d", time.Now().UnixNano())
	client.WS = ws
	item := clientList.PushBack(client)

	defer clientList.Remove(item)

	client.Channel = append(client.Channel, "Public")

	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("Can't receive:" + err.Error())
			var leaveMsg ChatMessage
			leaveMsg.Channel = "Public"
			leaveMsg.Timestamp = time.Now().Unix() * 1000
			leaveMsg.Message = fmt.Sprintf("%s logout GoChat.", client.UserName)
			leaveMsg.User = "GoChat Service"
			leaveMsg.Type = "system"
			processMessage(client, leaveMsg)
			break
		}
		var msg ChatMessage
		err = json.Unmarshal([]byte(reply), &msg)
		if err != nil {
			msg.Message = "Your Message format incorrect."
			msg_str, _ := json.Marshal(msg)
			if string(msg_str) != "" {
				websocket.Message.Send(ws, string(msg_str))
			}
		} else {
			newChannel := true
			for _, v := range client.Channel {
				if v == msg.Channel {
					newChannel = false
					break
				}
			}
			if newChannel {
				client.Channel = append(client.Channel, msg.Channel)
			}

			switch msg.Type {
			case "login":
				//Send Login Message
				client.UserName = msg.User
				var welcomeMsg ChatMessage
				welcomeMsg.Channel = "Public"
				welcomeMsg.Timestamp = time.Now().Unix() * 1000
				welcomeMsg.Message = "Welcome to GoChat."
				welcomeMsg.User = "GoChat Service"
				welcomeMsg.Type = "system"
				msg_str, _ := json.Marshal(welcomeMsg)
				if string(msg_str) != "" {
					websocket.Message.Send(ws, string(msg_str))
				}
				//User Login Messagevar welcomeMsg ChatMessage
				welcomeMsg.Channel = "Public"
				welcomeMsg.Timestamp = time.Now().Unix() * 1000
				welcomeMsg.Message = fmt.Sprintf("<b>%s</b> login GoChat.", msg.User)
				welcomeMsg.User = "GoChat Service"
				welcomeMsg.Type = "system"
				processMessage(client, welcomeMsg)
			case "message":
				processMessage(client, msg)
			case "command":
				processCommand(client, msg)
			}

		}
		if CONFIGS.Debug {
			log.Println("Received back from client: " + reply)
		}
	}
}

func processMessage(client *WsClient, msg ChatMessage) {
	switch CONFIGS.Mode {
	case "chat":
		broadcast(client.UId, msg)
	}
}

func processCommand(client *WsClient, msg ChatMessage) {
	var cmdMsg ChatMessage
	cmdMsg.Channel = "System"
	cmdMsg.Timestamp = time.Now().Unix() * 1000

	cmdMsg.User = "GoChat Service"
	cmdMsg.Type = "system"
	process := false
	//verify admin token
	if msg.Token == CONFIGS.AuthToken {
		process = true
	} else {
		cmdMsg.Message = "This coommand is Administrator only."
	}

	switch msg.Data.Command {
	case "list":
		cmdMsg.Message = "<p>User list</p>"
		for e := clientList.Front(); e != nil; e = e.Next() {
			if e.Value.(*WsClient).UId != client.UId {
				channellist, _ := json.Marshal(e.Value.(*WsClient).Channel)
				cmdMsg.Message += fmt.Sprintf("<b>%s</b> in Channel:%s</br>", e.Value.(*WsClient).UserName, channellist)
			}
		}
	case "?":
		cmdMsg.Message = `
				<label>Command List</label></br>
				<b>/list</b>: get online user list
			`
	case "kick":
		if process {
			if len(msg.Data.Args) == 2 {
				for e := clientList.Front(); e != nil; e = e.Next() {
					if e.Value.(*WsClient).UserName == msg.Data.Args[1] {
						kickErr := e.Value.(*WsClient).WS.Close()
						if kickErr == nil {
							cmdMsg.Message = fmt.Sprintf("<b>%s</b> has kicked.</br>", e.Value.(*WsClient).UserName)
						} else {
							cmdMsg.Message = fmt.Sprintf("<b>%s</b> has kick failed. Reason:%v</br>", e.Value.(*WsClient).UserName, kickErr)
						}
						break
					}
				}
			} else {
				cmdMsg.Message = "Args has empty."
			}
		}
	}

	if cmdMsg.Message != "" {
		msg_str, _ := json.Marshal(cmdMsg)
		if string(msg_str) != "" {
			websocket.Message.Send(client.WS, string(msg_str))
		}
	}
}

func broadcast(uid string, msg ChatMessage) {
	var err error
	msg_str, _ := json.Marshal(msg)
	for e := clientList.Front(); e != nil; e = e.Next() {
		if e.Value.(*WsClient).UId != uid {
			find := false
			for _, v := range e.Value.(*WsClient).Channel {
				if v == msg.Channel {
					find = true
					break
				}
			}
			if find {
				if err = websocket.Message.Send(e.Value.(*WsClient).WS, string(msg_str)); err != nil {
					log.Printf("WsSokect error:%s\n", err.Error())
				}
			}
		}
	}
}

func subProtocolHandshake(config *websocket.Config, req *http.Request) error {
	//if need set rules for protocol

	var protoList []string = []string{"hash", "kv", "test"}

	accept := false
	for _, proto := range config.Protocol {
		//check the protocol is our service protocol
		for _, serverProto := range protoList {
			if proto == serverProto {
				accept = true
				break
			}
		}
	}

	if accept {
		return nil
	}
	return websocket.ErrBadWebSocketProtocol

}

func subProtoServer(ws *websocket.Conn) {
	var err error
	client := new(WsClient)
	client.UId = fmt.Sprintf("U-%d", time.Now().UnixNano())
	client.WS = ws
	item := clientList.PushBack(client)

	defer clientList.Remove(item)

	sub := make(map[string]interface{})
	var subproto []string
	for _, proto := range ws.Config().Protocol {
		subproto = append(subproto, proto)
	}
	sub["subscription"] = subproto
	if err = websocket.JSON.Send(ws, sub); err != nil {
		log.Printf("WsSokect error:%s\n", err.Error())
	}

	for {
		if client.Closed {
			log.Println("Client has closed.", client.UId)
			break
		}
		time.Sleep(time.Second)
	}

}

func pubProtoServer(ws *websocket.Conn) {
	var err error
	client := new(WsClient)
	client.UId = fmt.Sprintf("U-%d", time.Now().UnixNano())
	client.WS = ws
	item := clientList.PushBack(client)

	defer clientList.Remove(item)

	for {
		var reply string
		if err = websocket.JSON.Receive(ws, &reply); err != nil {
			log.Println("Can't receive:" + err.Error())
			if err == io.EOF {
				break
			} else {
				response := make(map[string]interface{})
				response["error"] = err.Error()
				if err = websocket.JSON.Send(ws, response); err != nil {
					log.Printf("WsSokect error:%s\n", err.Error())
					break
				}
			}
		}
	}

}

func ProcessSubRequest(ws *websocket.Conn, message string) {

}

func SendSubMessage(subproto string, msg string) {
	for e := clientList.Front(); e != nil; e = e.Next() {
		wc := e.Value.(*WsClient)
		if !wc.Closed {
			for _, proto := range wc.WS.Config().Protocol {
				if proto == subproto {
					if err := websocket.Message.Send(wc.WS, fmt.Sprintf("submessage:%s %s", msg)); err != nil {
						log.Printf("WsSokect error:%s\n", err.Error())
						if err == io.EOF {
							wc.Close()
							clientList.Remove(e)
						}
					}
					break
				}
			}
		}
	}
}

func SendSubInterface(subproto string, msg interface{}) {
	for e := clientList.Front(); e != nil; e = e.Next() {
		wc := e.Value.(*WsClient)
		if !wc.Closed {
			for _, proto := range wc.WS.Config().Protocol {
				if proto == subproto {
					response := make(map[string]interface{})
					response["protocol"] = subproto
					response["data"] = msg
					if err := websocket.JSON.Send(wc.WS, response); err != nil {
						log.Printf("WsSokect error:%s\n", err.Error())
						if err == io.EOF {
							wc.Close()
							clientList.Remove(e)
						}
					}
					break
				}
			}
		}
	}
}

func WsServer() {
	http.Handle("/socket", websocket.Handler(Echo))
	subproto := websocket.Server{
		Handler: websocket.Handler(subProtoServer),
	}

	http.Handle("/sub", subproto)
	pubproto := websocket.Server{
		Handler: websocket.Handler(pubProtoServer),
	}

	http.Handle("/pub", pubproto)
	clientList = list.New()
	log.Println("Web socket service starting.")

}
