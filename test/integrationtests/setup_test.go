package integrationtests

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/Projects-for-Fun/thefoodbook/internal/repository"
	"github.com/Projects-for-Fun/thefoodbook/internal/service"
	"github.com/Projects-for-Fun/thefoodbook/pkg/database"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/auth"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
)

func runIntegrationTests() bool {
	shouldRunIntegrationTests, err := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS"))

	if err != nil {
		return false
	}

	return shouldRunIntegrationTests
}

func setup(t *testing.T, ctx context.Context) (*configs.Config, neo4j.DriverWithContext, *webservice.Webservice, zerolog.Logger) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Caller().Timestamp().Logger()
	ctx = logging.AttachLogger(ctx, logger)

	config, err := configs.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewDriver(ctx, config.DBURI, config.DBUser, config.DBPass)
	if err != nil {
		t.Fatal(err)
	}

	w := webservice.NewWebservice(
		service.CreateUserServiceFunc(repository.CreateUserRepoFunc(db)),
		service.LoginUserServiceFunc(repository.GetUserByUsernameRepoFunc(db), auth.VerifyPassword, repository.SetUserLastLoginRepo(db)),
	)

	return config, db, w, logger
}

func LoggerTestMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(logging.AttachLogger(r.Context(), logger)))
		}

		return http.HandlerFunc(fn)
	}
}

func AuthTestMiddleware(_ []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}
