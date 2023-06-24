package mws

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/tracing"

	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			corrID := tracing.GetCorrelationID(r.Context())

			logger := logger.With().Str(tracing.CorrelationIDTag, corrID).Caller().Logger()

			next.ServeHTTP(w, r.WithContext(logging.AttachLogger(r.Context(), logger)))
		}

		return http.HandlerFunc(fn)
	}
}
