package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"os"
	quiz "quiz-tool/src"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app := fiber.New()
	app.Static("/", "./public")
	app.Use(fiberlogger.New())
	app.Use(cors.New())

	userGroup := app.Group("/user")
	user := quiz.NewUser(userGroup, logger)
	go user.Run()

	adminGroup := app.Group("/admin")
	admin := quiz.NewAdmin(adminGroup, logger)
	go admin.Run()

	overviewGroup := app.Group("/overview")
	overview := quiz.NewAdmin(overviewGroup, logger)
	go overview.Run()

	logger.Fatal().Err(app.Listen(":8898"))
}
