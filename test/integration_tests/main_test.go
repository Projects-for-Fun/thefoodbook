//go:build integration_tests

package integration_tests

import (
	"os"
	"testing"
)

var ApiURL string

func TestMain(m *testing.M) {
	if ApiURL, ok := os.LookupEnv("API_URL"); ok {
		panic("Invalid API_URL")
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
