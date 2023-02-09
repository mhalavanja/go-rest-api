package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mhalavanja/go-rest-api/util"
)

type message struct {
	data    []byte
	groupId int64
}

type hub struct {
	groups     map[int64]map[*connection]bool
	broadcast  chan message
	register   chan subscription
	unregister chan subscription
	upgrader   *websocket.Upgrader
}

func NewUpgrader(config *util.Config) *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return config.Client == r.Header.Get("Origin")
		},
	}
}

func NewHub(config *util.Config) *hub {
	return &hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		groups:     make(map[int64]map[*connection]bool),
		upgrader:   NewUpgrader(config),
	}

}

func (hub *hub) ServeWs(ctx *gin.Context) {
	groupId, _ := strconv.Atoi(ctx.Param("groupId"))

	ws, err := hub.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	conn := &connection{
		send: make(chan []byte, 256),
		ws:   ws,
	}
	sub := subscription{
		hub:     hub,
		conn:    conn,
		groupId: int64(groupId),
	}
	hub.register <- sub

	go sub.writePump()
	go sub.readPump()
}

func (hub *hub) Run() {
	for {
		select {
		case sub := <-hub.register:
			connections := hub.groups[sub.groupId]
			if connections == nil {
				connections = make(map[*connection]bool)
				hub.groups[sub.groupId] = connections
			}
			hub.groups[sub.groupId][sub.conn] = true

		case sub := <-hub.unregister:
			connections := hub.groups[sub.groupId]
			if connections != nil {
				if _, ok := connections[sub.conn]; ok {
					delete(connections, sub.conn)
					close(sub.conn.send)
					if len(connections) == 0 {
						delete(hub.groups, sub.groupId)
					}
				}
			}

		case msg := <-hub.broadcast:
			connections := hub.groups[msg.groupId]
			for conn := range connections {
				select {
				case conn.send <- msg.data:
				default:
					close(conn.send)
					delete(connections, conn)
					if len(connections) == 0 {
						delete(hub.groups, msg.groupId)
					}
				}
			}
		}
	}
}
