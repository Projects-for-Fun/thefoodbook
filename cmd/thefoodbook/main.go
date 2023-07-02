package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

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

	var logger zerolog.Logger
	if config.LogFormat == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Caller().Timestamp().Logger()
	}
	ctx = logging.AttachLogger(ctx, logger)

	logger.Info().Msgf("Loading variables for %s environment.", config.Environment)

	db, err := database.NewDriver(ctx, config.DBURI, config.DBUser, config.DBPass)
	if err != nil {
		logger.Fatal().Err(err).Caller().Msg("Couldn't connect to the db.")
	}
	logger.Info().Msg("Connected to db. Obtained new driver with context.")

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
			logger.Error().Err(err).Msg("Webservice stopped")
			database.CloseDriver(ctxWithCancel, db, logger)
		}
	default:
		logger.Error().Msg("Mistakes were made")
	}
}
