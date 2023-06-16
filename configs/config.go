package configs

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

type Config struct {
	Environment string

	ServiceName string
	ServicePort string

	LogLevel  string
	LogFormat string
}

func NewConfig() (*Config, error) {
	var config Config

	serviceName, err := loadEnvironmentVariable("SERVICE_NAME")
	if err != nil {
		return nil, err
	}
	config.ServiceName = serviceName

	env, err := loadEnvironmentVariable("ENVIRONMENT")
	if err != nil {
		// Default to local
		env = "local"
	}
	config.Environment = strings.ToLower(env)

	servicePort, err := loadEnvironmentVariable("SERVICE_PORT")
	if err != nil {
		return nil, err
	}
	config.ServicePort = servicePort

	logLevel, err := loadEnvironmentVariable("LOG_LEVEL")
	if err != nil {
		// Default to info
		logLevel = "info"
	}
	config.LogLevel = logLevel
	zerologLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		// Default to INFO
		zerologLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(zerologLevel)

	logFormat, err := loadEnvironmentVariable("LOG_FORMAT")
	if err != nil {
		return nil, err
	}
	config.LogFormat = logFormat

	return &config, nil
}

func loadEnvironmentVariable(variable string) (string, error) {
	if value, ok := os.LookupEnv(variable); ok {
		return value, nil
	}

	return "", fmt.Errorf("environment variable %s could not be set", variable)
}