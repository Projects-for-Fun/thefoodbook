package webservice

import (
	"log"
	"net/http"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/mws"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func RunWebservice(config *configs.Config, _ neo4j.DriverWithContext, logger zerolog.Logger) error {
	log.Println("Running Webservice")
	_ = webservice.NewWebservice()

	router := chi.NewRouter()

	middleware.RequestIDHeader = "X-Correlation-Id"

	router.Use(middleware.Timeout(5 * time.Second))
	router.Use(mws.CorrelationID)
	router.Use(mws.LoggerWithRecoverer(logger))

	router.Get("/status", func(w http.ResponseWriter, r *http.Request) { /* Empty status function. */ })

	return http.ListenAndServe(":"+config.ServicePort, router)
}
