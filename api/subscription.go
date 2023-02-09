package api

import (
	"log"

	"github.com/gorilla/websocket"
)

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

type subscription struct {
	hub     *hub
	conn    *connection
	groupId int64
}

func (sub subscription) readPump() {
	conn := sub.conn
	defer func() {
		sub.hub.unregister <- sub
		conn.ws.Close()
	}()
	for {
		_, msg, err := conn.ws.ReadMessage()
		if err != nil {
			log.Println("ERROR: ", err)
			break
		}
		m := message{msg, sub.groupId}
		sub.hub.broadcast <- m
	}
}

func (sub *subscription) writePump() {
	conn := sub.conn
	for msg := range conn.send {
		conn.ws.WriteMessage(websocket.TextMessage, msg)
	}
	conn.ws.Close()
}
