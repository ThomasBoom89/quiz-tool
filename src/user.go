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
		generatedUuid, err := uuid.NewRandom()
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"name":    loginRequestBody.Name,
			"id":      generatedUuid.String(),
			"isAdmin": false,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}
		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		encodedToken, err := token.SignedString([]byte("secret"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if loginRequestBody.RoomId != "" {
			return ctx.JSON(fiber.Map{"token": encodedToken})
		}

		return ctx.SendStatus(fiber.StatusUnauthorized)
	})

	// JWT Middleware
	U.router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	U.router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	U.router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		U.handleWebsocket(c)
	}))
}

func (U *User) handleWebsocket(c *websocket.Conn) {
	// Register the client
	U.communication.RegisterUser(c)
}
