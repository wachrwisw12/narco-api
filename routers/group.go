package routers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoute(app *fiber.App) {
	api := app.Group("/api")

	SetupNarcoticsReport(api.Group("/v1"))
	SetupInnitApp(api.Group("/init"))
	SetupAuth(api.Group("/auth"))
}

var clients = make(map[*websocket.Conn]bool)

func SetupWS(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		clients[conn] = true
		log.Println("client connected")

		defer func() {
			delete(clients, conn)
			conn.Close()
			log.Println("client disconnected")
		}()

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))
}
