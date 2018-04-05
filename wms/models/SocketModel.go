package models

import (
	"errors"
	"fmt"
	"message_server/utils/logs"
	"net/http"

	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
)

var SocketPool = make(map[string]*websocket.Conn)

func InitSocket() {
	port := beego.AppConfig.String("socket::port")
	logs.Info("Websocket Running on :" + port)

	http.Handle("/ws", websocket.Handler(SocketHandler))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logs.Debug("Websocket run error:", err)
	}
}

func SocketHandler(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}
		SocketPool[reply] = ws
		fmt.Println("Received back from client: " + reply)
		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg + "_server")
		if err = websocket.Message.Send(ws, msg+"_server"); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func Send(uid, msg string) error {
	ws := SocketPool[uid]
	if ws == nil {
		return errors.New("Non-existent websocket connection")
	}
	if err := websocket.Message.Send(ws, msg); err != nil {
		return err
	}
	return nil
}
