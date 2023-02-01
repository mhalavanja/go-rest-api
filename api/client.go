package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mhalavanja/go-rest-api/util"
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (sub subscription) readPump() {
	conn := sub.conn
	defer func() {
		sub.hub.unregister <- sub
		conn.ws.Close()
	}()
	for {
		_, msg, err := conn.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := message{msg, sub.groupId}
		sub.hub.broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (sub *subscription) writePump() error {
	conn := sub.conn
	ticker := time.NewTicker(time.Minute)
	defer func() {
		ticker.Stop()
		conn.ws.Close()
	}()
	for {
		select {
		case message, ok := <-conn.send:
			if !ok {
				conn.ws.WriteMessage(websocket.CloseMessage, []byte{})
			}
			conn.ws.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			conn.ws.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func NewUpgrader(config *util.Config) websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return config.Client == r.Header.Get("Origin")
		},
	}
}

func (hub *hub) ServeWs(ctx *gin.Context) {
	groupId, _ := strconv.Atoi(ctx.Param("groupId"))

	ws, err := hub.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{
		send: make(chan []byte, 256),
		ws:   ws,
	}
	sub := subscription{
		hub:     hub,
		conn:    c,
		groupId: int64(groupId),
	}
	hub.register <- sub
	go sub.writePump()
	go sub.readPump()
}
