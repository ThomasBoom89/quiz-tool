package quiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"strconv"
)

type Overview struct {
	router         fiber.Router
	logger         zerolog.Logger
	connectionPool *ConnectionPool
}

func NewOverview(router fiber.Router, logger zerolog.Logger, connectionPool *ConnectionPool) *Overview {
	overview := &Overview{router: router, logger: logger, connectionPool: connectionPool}
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
	// When the function returns, unregister the client and close the connection
	defer func() {
		O.connectionPool.Unregister(c)
		c.Close()
	}()

	// Register the client
	O.connectionPool.Register(c)

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				O.logger.Error().Str("read error", err.Error())
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			O.connectionPool.Broadcast(string(message))
		} else {
			O.logger.Error().Str("websocket message received of type", strconv.Itoa(messageType))
		}
	}
}
