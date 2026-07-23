package main

import (
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/config"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	appbuilder "github.com/Tihomir-Tinkov/cooking-site-project/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("failed to load .env file, using system environment")
	}

	conf := config.Config{}

	if err := conf.Parse(); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}

	logger.Init(conf.Logger)

	app, err := appbuilder.NewApp(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build app")
	}

	defer func(app *appbuilder.App) {
		if err := app.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close app")
		}
	}(app)

	log.Info().Str("addr", app.Server.Addr).Msg("server started")
	err = app.ServeAndListen()
	if err != nil {
		log.Fatal().Err(err).Msg("server failed to start")
	}
}
