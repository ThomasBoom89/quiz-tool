package quiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type adminLoginRequestBody struct {
	Username string
	Password string
}

type Admin struct {
	router        fiber.Router
	logger        zerolog.Logger
	communication *Communication
}

func NewAdmin(router fiber.Router, logger zerolog.Logger, communication *Communication) *Admin {
	admin := &Admin{router: router, logger: logger, communication: communication}
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
		// todo: return jwt
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

	// todo jwt middleware
	A.router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		jwt := c.Params("jwt", "")
		A.handleWebsocket(c, jwt)
	}))
}

func (A *Admin) handleWebsocket(c *websocket.Conn, jwt string) {
	// Register the client
	A.communication.RegisterAdmin(c, jwt)
}
