package mws

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

var CorrelationIDHeader = "X-Correlation-Id"
var RequestIDKey = "CorrelationKey"

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, RequestIDKey, correlationID)
		w.Header().Set(CorrelationIDHeader, correlationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
