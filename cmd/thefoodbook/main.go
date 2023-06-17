package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/cmd/thefoodbook/webservice"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/rs/zerolog"
)

func initializeDependencies(_ context.Context) (*configs.Config, zerolog.Logger) {
	config, err := configs.NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to set project configuration: %v", err))
	}

	// TODO: update this
	// https://github.com/rs/zerolog#pass-a-sub-logger-by-context
	// https://github.com/rs/zerolog#contextcontext-integration
	// https://github.com/rs/zerolog#integration-with-nethttp
	var logger zerolog.Logger
	if config.LogFormat == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Timestamp().Logger()
	}

	logger.Info().Msgf("Loading variables for %s environment.", config.Environment)
	logger.Info().Msgf("Running on port %s.", config.ServicePort)

	return config, logger
}

func main() {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	config, logger := initializeDependencies(ctxWithTimeout)

	if len(os.Args) < 2 {
		logger.Fatal().Msg("Must provide program argument")
	}

	switch os.Args[1] {
	case "webservice":
		err := webservice.RunWebservice(config, logger)

		if err != nil {
			logger.Err(err).Msg("Webservice stopped")
		}
	default:
		logger.Error().Msg("Mistakes were made")
	}
}
