package database

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
)

func NewDB(ctx context.Context, DBURI, DBUser, DBPass string, logger zerolog.Logger) neo4j.DriverWithContext {
	driver, err := newDriver(ctx, DBURI, DBUser, DBPass)

	if err != nil {
		logger.Fatal().Err(err).Caller().Msg("Couldn't connect to the db.")
	}

	logger.Info().Msg("Connected to db. Obtained new driver with context.")
	return driver
}

func newDriver(ctx context.Context, uri, username, password string) (neo4j.DriverWithContext, error) {
	// Create Driver
	driverWithContext, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))

	// Handle any driver creation errors
	if err != nil {
		return nil, err
	}

	// Verify Connectivity
	err = driverWithContext.VerifyConnectivity(ctx)

	// If connectivity fails, handle the error
	if err != nil {
		return nil, err
	}

	return driverWithContext, nil
}

// CloseDriver call on application exit
func CloseDriver(ctx context.Context, driver neo4j.DriverWithContext, logger zerolog.Logger) {
	logger.Info().Msg("Closing the driver")
	err := driver.Close(ctx)

	if err != nil {
		logger.Err(err).Msg("Couldn't close the db.")
	}
}
