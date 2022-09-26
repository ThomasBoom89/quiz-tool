package quiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type Overview struct {
	router        fiber.Router
	logger        zerolog.Logger
	communication *Communication
}

func NewOverview(router fiber.Router, logger zerolog.Logger, communication *Communication) *Overview {
	overview := &Overview{router: router, logger: logger, communication: communication}
	overview.attachRoutes()

	return overview
}

func (O *Overview) attachRoutes() {
	O.router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	O.router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		O.handleWebsocket(c)
	}))
}

func (O *Overview) handleWebsocket(c *websocket.Conn) {
	// Register the client
	O.communication.RegisterOverview(c)
}
