package routes

import (
	"sapasora/platform/server"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// ChannelRouter is the router for the WebSocket, usually everything that is not an API or will be consumed by this application only
// @wired:provide(group=router)
func ChannelRouter(
	channel *Channel,
) server.Router {
	cfg := websocket.Config{
		RecoverHandler: func(c *websocket.Conn) {
			if err := recover(); err != nil {
				c.WriteJSON(fiber.Map{
					"message": "error occurred",
					"error":   fmt.Sprintf("%v", err),
				})
			}
		},
	}

	return server.NewRouter("/ws", func(r fiber.Router) {
		channel := r.Group("/").Name("channel.")
		channel.Get("/ping", websocket.New(func(c *websocket.Conn) {
			c.WriteMessage(websocket.TextMessage, []byte("pong!"))
		}, cfg)).Name("ping")
	}, channel.PipeLine()...)
}
