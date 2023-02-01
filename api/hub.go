package api

import (
	"github.com/gorilla/websocket"
)

type message struct {
	data    []byte
	groupId int64
}

type subscription struct {
	hub     *hub
	conn    *connection
	groupId int64
}

type hub struct {
	groups map[int64]map[*connection]bool
	// Inbound messages from the connections.
	broadcast chan message
	// Register requests from the connections.
	register chan subscription
	// Unregister requests from connections.
	unregister chan subscription
	upgrader   *websocket.Upgrader
}

func NewHub(upgrader *websocket.Upgrader) *hub {
	return &hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		groups:     make(map[int64]map[*connection]bool),
		upgrader:   upgrader,
	}

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
