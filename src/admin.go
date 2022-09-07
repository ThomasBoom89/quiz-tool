package quiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"strconv"
)

type adminLoginRequestBody struct {
	Username string
	Password string
}

type Admin struct {
	router         fiber.Router
	logger         zerolog.Logger
	connectionPool *ConnectionPool
}

func NewAdmin(router fiber.Router, logger zerolog.Logger, connectionPool *ConnectionPool) *Admin {
	admin := &Admin{router: router, logger: logger, connectionPool: connectionPool}
	admin.attachRoutes()

	return admin
}

func (A *Admin) attachRoutes() {

	A.router.Post("/login", func(ctx *fiber.Ctx) error {
		loginRequestBody := new(adminLoginRequestBody)
		if err := ctx.BodyParser(loginRequestBody); err != nil {
			A.logger.Debug().Str("error login", err.Error())
			return err
		}
		if loginRequestBody.Password != "" {
			return ctx.SendStatus(fiber.StatusOK)
		}

		return ctx.SendStatus(fiber.StatusUnauthorized)
	})

	A.router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	A.router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		A.handleWebsocket(c)
	}))
}

func (A *Admin) handleWebsocket(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		A.connectionPool.Unregister(c)
		c.Close()
	}()

	// Register the client
	A.connectionPool.Register(c)

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				A.logger.Error().Str("read error", err.Error())
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			A.connectionPool.Broadcast(string(message))
		} else {
			A.logger.Error().Str("websocket message received of type", strconv.Itoa(messageType))
		}
	}
}
