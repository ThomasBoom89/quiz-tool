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

	userConnectionPool := quiz.NewConnectionPool(logger)
	go userConnectionPool.Run()
	userGroup := app.Group("/user")
	quiz.NewUser(userGroup, logger, userConnectionPool)

	adminConnectionPool := quiz.NewConnectionPool(logger)
	go adminConnectionPool.Run()
	adminGroup := app.Group("/admin")
	quiz.NewAdmin(adminGroup, logger, adminConnectionPool)

	overviewConnectionPool := quiz.NewConnectionPool(logger)
	go overviewConnectionPool.Run()
	overviewGroup := app.Group("/overview")
	quiz.NewOverview(overviewGroup, logger, overviewConnectionPool)

	logger.Fatal().Err(app.Listen(":8898"))
}
