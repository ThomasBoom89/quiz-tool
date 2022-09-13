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
	app.Use(fiberlogger.New())
	app.Use(cors.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	userConnectionPool := quiz.NewConnectionPool(logger)
	go userConnectionPool.Run()
	userGroup := v1.Group("/user")
	quiz.NewUser(userGroup, logger, userConnectionPool)

	adminConnectionPool := quiz.NewConnectionPool(logger)
	go adminConnectionPool.Run()
	adminGroup := v1.Group("/admin")
	quiz.NewAdmin(adminGroup, logger, adminConnectionPool)

	overviewConnectionPool := quiz.NewConnectionPool(logger)
	go overviewConnectionPool.Run()
	overviewGroup := v1.Group("/overview")
	quiz.NewOverview(overviewGroup, logger, overviewConnectionPool)

	app.Static("/", "./public")
	app.Static("**", "./public")

	logger.Fatal().Err(app.Listen(":8898"))
}
