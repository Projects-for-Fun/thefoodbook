package mws

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

	"github.com/go-chi/chi/v5/middleware"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger(r.Context())

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()

		defer func() {
			t2 := time.Now()

			logger.Info().Fields(map[string]interface{}{
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

		next.ServeHTTP(w, r)
	})
}
