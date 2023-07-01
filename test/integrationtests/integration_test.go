package integrationtests

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/Projects-for-Fun/thefoodbook/internal/repository"
	"github.com/Projects-for-Fun/thefoodbook/internal/service"
	"github.com/Projects-for-Fun/thefoodbook/pkg/database"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/auth"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"
	"github.com/rs/zerolog"
)

var (
	//nolint:all
	db     neo4j.DriverWithContext
	logger zerolog.Logger

	w *webservice.Webservice
)

func runIntegrationTests() bool {
	shouldRunIntegrationTests, err := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS"))

	if err != nil {
		return false
	}

	return shouldRunIntegrationTests
}

func setup(t *testing.T, ctx context.Context) {
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Caller().Timestamp().Logger()

	config, err := configs.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewDriver(ctx, config.DBURI, config.DBUser, config.DBPass, logger)
	if err != nil {
		t.Fatal(err)
	}

	w = webservice.NewWebservice(
		service.CreateUserServiceFunc(repository.CreateUserRepoFunc(db)),
		service.LoginUserServiceFunc(repository.GetUserByUsernameRepoFunc(db), auth.VerifyPassword, repository.SetUserLastLoginRepo(db)),
	)
}

func LoggerTestMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(logging.AttachLogger(r.Context(), logger)))
		}

		return http.HandlerFunc(fn)
	}
}
