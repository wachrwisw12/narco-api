package ws

import "github.com/gofiber/websocket/v2"

var Clients = make(map[*websocket.Conn]bool)

func Broadcast(msg string) {
	for c := range Clients {
		c.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
