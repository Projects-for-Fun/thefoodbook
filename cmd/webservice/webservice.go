package webservice

import (
	"log"
	"net/http"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/configs"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/mws"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func RunWebservice(config *configs.Config, logger zerolog.Logger) error {
	log.Println("Running Webservice")
	w := webservice.NewWebservice()

	router := chi.NewRouter()

	middleware.RequestIDHeader = "X-Correlation-Id"

	router.Use(middleware.Timeout(5 * time.Second))
	router.Use(mws.CorrelationId)
	router.Use(mws.LoggerWithRecoverer(logger))

	router.Get("/status", func(w http.ResponseWriter, r *http.Request) {})

	router.Route("/dummy", func(router chi.Router) {
		router.With(mws.Paginate).Get("/", w.GetDummyData) // GET /dummy or GET /dummy?id=1

		router.Route("/{id}", func(router chi.Router) {
			router.Get("/", w.GetDummyDataForId) // GET /dummy/1
		})
	})

	return http.ListenAndServe(":"+config.ServicePort, router)
}
