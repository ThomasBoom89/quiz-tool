package main

import (
	alp "github.com/ThomasBoom89/accept-language-parser"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"os"
	quiz "quiz-tool/src"
	"regexp"
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

	app.Get("/:filename?/**", handleStatic(logger))

	logger.Fatal().Err(app.Listen(":8898"))
}

func handleStatic(logger zerolog.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		dirs, err := os.ReadDir("./public")
		if err != nil {
			logger.Panic().Err(err)
		}

		filename := ctx.Params("filename", "")
		re, _ := regexp.Compile(`[a-zA-Z0-9]*\.(js|css|txt|html)`)
		if !re.MatchString(filename) {
			filename = ""
		}

		headers := ctx.GetReqHeaders()
		acceptedLanguageHeader := headers["Accept-Language"]
		languages, err := alp.Parse(acceptedLanguageHeader)
		logger.Debug().Msgf("languages detected: ", languages)
		if err != nil {
			return ctx.SendFile("./public/en-US/" + filename)
		}

		for _, dir := range dirs {
			if dir.IsDir() && dir.Name() == languages[0].Name {
				return ctx.SendFile("./public/" + languages[0].Name + "/" + filename)
			}
		}

		return ctx.SendFile("./public/en-US/" + filename)
	}
}
