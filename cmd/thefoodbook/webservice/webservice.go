package webservice

import (
	"net/http"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/internal/repository"
	"github.com/Projects-for-Fun/thefoodbook/internal/service"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/mws"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func RunWebservice(config *configs.Config, db neo4j.DriverWithContext, logger zerolog.Logger) error {
	logger.Info().Msg("Running Webservice")
	w := webservice.NewWebservice(service.HandleCreateUserFunc(repository.CreateUserRepoFunc(db)))

	router := chi.NewRouter()

	router.Use(middleware.Timeout(50 * time.Second))
	router.Use(mws.CorrelationID)
	router.Use(mws.LoggerWithRecoverer(logger))

	router.Get("/status", func(w http.ResponseWriter, r *http.Request) { /* Empty status function. */ })

	router.Post("/sign-up", w.HandleSignUp)

	return http.ListenAndServe(":"+config.ServicePort, router)
}
