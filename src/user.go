package quiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type userLoginRequestBody struct {
	Name   string
	RoomId string
}

type User struct {
	router        fiber.Router
	logger        zerolog.Logger
	communication *Communication
}

func NewUser(router fiber.Router, logger zerolog.Logger, communication *Communication) *User {
	user := &User{router: router, logger: logger, communication: communication}
	user.attachRoutes()

	return user
}

func (U *User) attachRoutes() {

	U.router.Post("/login", func(ctx *fiber.Ctx) error {
		loginRequestBody := new(userLoginRequestBody)
		if err := ctx.BodyParser(loginRequestBody); err != nil {
			U.logger.Debug().Str("error login", err.Error())
			return err
		}
		// todo: return jwt
		if loginRequestBody.RoomId != "" {
			return ctx.JSON(fiber.Map{"token": "secrettoken"})
		}

		return ctx.SendStatus(fiber.StatusUnauthorized)
	})

	U.router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	// todo jwt middleware
	U.router.Get("/ws/:jwt?", websocket.New(func(c *websocket.Conn) {
		jwt := c.Params("jwt", "")
		U.handleWebsocket(c, jwt)
	}))
}

func (U *User) handleWebsocket(c *websocket.Conn, jwt string) {
	// Register the client
	U.communication.RegisterUser(c, jwt)
}
