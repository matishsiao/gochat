package main

import (
	"github.com/matishsiao/net/websocket"
	"container/list"
	"fmt"
	"log"
	"net/http"
	_ "strconv"
	"strings"
	"time"
	"io"
)

var clientList *list.List
type WsClient struct {
	UId string
	WS  *websocket.Conn
	Closed bool
}

func (wc *WsClient) Close() {
	wc.Closed = true
	wc.WS.Close()
}

func EchoServer(ws *websocket.Conn) {
    io.Copy(ws, ws)
}

func Echo(ws *websocket.Conn) {
	var err error
	client := new(WsClient)
	client.UId = fmt.Sprintf("U-%d",time.Now().UnixNano())
	client.WS = ws
	item := clientList.PushBack(client)

	defer clientList.Remove(item)
	
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("Can't receive:" + err.Error())
			break
		}

		msgs := strings.Split(reply, ",")
		if msgs[0] == "0" {
			client.UId = msgs[1]
		}

		fmt.Println("Received back from client: " + reply)
		broadcast(client.UId, reply)

		if msgs[0] == "0" {
			msg := "Received:  " + reply
			//fmt.Println("Sending to client: " + msg)

			if err = websocket.Message.Send(ws, msg); err != nil {
				fmt.Printf("WsSokect error:%s\n", err.Error())
				break
			}
		}
	}
}

func broadcast(uid string, msg string) {
	var err error
	//fmt.Println("WsServer broadcast:" + msg)
	
	for e := clientList.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value.(*WsClient).UId + " now:" + uid + "\n")
		if e.Value.(*WsClient).UId != uid {
			if err = websocket.Message.Send(e.Value.(*WsClient).WS, msg); err != nil {
				fmt.Printf("WsSokect error:%s\n", err.Error())
			}
		}
	}
}

func subProtocolHandshake(config *websocket.Config, req *http.Request) error {
	//if need set rules for protocol
	
	var protoList []string = []string{"hash","kv","test"}
	
	accept := false
	for _, proto := range config.Protocol {
		//check the protocol is our service protocol
		for _,serverProto := range protoList {
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
	client.UId = fmt.Sprintf("U-%d",time.Now().UnixNano())
	client.WS = ws
	item := clientList.PushBack(client)

	defer clientList.Remove(item)
	
	sub := make(map[string]interface{})
	var subproto []string
	for _, proto := range ws.Config().Protocol {
		subproto = append(subproto,proto)
	}
	sub["subscription"] = subproto
	if err = websocket.JSON.Send(ws, sub); err != nil {
		log.Printf("WsSokect error:%s\n", err.Error())
	}
	
	for {
		if client.Closed {
			log.Println("Client has closed.",client.UId)
			break
		}
		time.Sleep(time.Second)
	}
	
}

func pubProtoServer(ws *websocket.Conn) {
	var err error
	client := new(WsClient)
	client.UId = fmt.Sprintf("U-%d",time.Now().UnixNano())
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

func ProcessSubRequest(ws *websocket.Conn,message string) {
	 
}

func SendSubMessage(subproto string, msg string) {
	for e := clientList.Front(); e != nil; e = e.Next() {
		wc := e.Value.(*WsClient)
		if !wc.Closed {
			for _, proto := range wc.WS.Config().Protocol {
				if proto == subproto {
					if err := websocket.Message.Send(wc.WS, fmt.Sprintf("submessage:%s %s",msg)); err != nil {
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
		Handler:   websocket.Handler(subProtoServer),
	}
		
	http.Handle("/sub", subproto)
	pubproto := websocket.Server{
		Handler:   websocket.Handler(pubProtoServer),
	}
		
	http.Handle("/pub", pubproto)
	clientList = list.New()
	log.Println("Web socket service starting.") 

}
