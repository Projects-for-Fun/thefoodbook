package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/Projects-for-Fun/thefoodbook/cmd/thefoodbook/webservice"
	"github.com/Projects-for-Fun/thefoodbook/pkg/database"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/rs/zerolog"
)

func initializeDependencies(ctx context.Context) (*configs.Config, neo4j.DriverWithContext, zerolog.Logger) {
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
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Caller().Timestamp().Logger()
	}

	logger.Info().Msgf("Loading variables for %s environment.", config.Environment)

	db := database.NewDB(ctx, config.DBURI, config.DBUser, config.DBPass, logger)

	return config, db, logger
}

func main() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, db, logger := initializeDependencies(ctxWithCancel)

	if len(os.Args) < 2 {
		logger.Fatal().Msg("Must provide program argument")
	}

	switch os.Args[1] {
	case "webservice":
		err := webservice.RunWebservice(config, db, logger)

		if err != nil {
			logger.Fatal().Err(err).Msg("Webservice stopped")
			database.CloseDriver(ctxWithCancel, db, logger)
		}
	default:
		logger.Error().Msg("Mistakes were made")
	}
}
