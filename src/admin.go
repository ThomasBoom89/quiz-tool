package quiz

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
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
		generatedUuid, err := uuid.NewRandom()
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"name":    loginRequestBody.Username,
			"id":      generatedUuid.String(),
			"isAdmin": true,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}
		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		encodedToken, err := token.SignedString([]byte("secret"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if loginRequestBody.Password != "" {
			return ctx.JSON(fiber.Map{"token": encodedToken})
		}

		return ctx.SendStatus(fiber.StatusUnauthorized)
	})

	// JWT Middleware
	A.router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

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
	// Register the client
	A.communication.RegisterAdmin(c)
}
