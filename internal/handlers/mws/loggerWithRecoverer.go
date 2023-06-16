package mws

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func LoggerWithRecoverer(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			correlationID := fmt.Sprintf("%v", r.Context().Value(RequestIDKey))

			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			t1 := time.Now()

			log := logger.With().Str("correlation-id", correlationID).Logger()

			defer func() {
				t2 := time.Now()

				log.Info().Fields(map[string]interface{}{
					"remote_ip":  r.RemoteAddr,
					"host":       r.Host,
					"proto":      r.Proto,
					"method":     r.Method,
					"user_agent": r.Header.Get("User-Agent"),
					"status":     ww.Status(),
					"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  ww.BytesWritten(),
				}).Msg(r.RequestURI)

				if rvr := recover(); rvr != nil {
					logger.Error().Interface("recover_info", rvr).Bytes("debug_stack", debug.Stack()).Msg("Panic on request")

					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

			}()

			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}
