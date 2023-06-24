package mws

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

var LoggerKey = "LoggerKey"

func LoggerWithRecoverer(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			correlationID := r.Context().Value(RequestIDKey).(string)

			log := logger.With().Str("correlation-id", correlationID).Caller().Logger()
			ctx := context.WithValue(r.Context(), LoggerKey, log)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

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

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func GetLogger(ctx context.Context) zerolog.Logger {
	return ctx.Value(LoggerKey).(zerolog.Logger)
}
