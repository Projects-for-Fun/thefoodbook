package mws

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

var CorrelationIDHeader = "X-Correlation-Id"
var RequestIDKey = "CorrelationKey"

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, RequestIDKey, correlationID)
		rw.Header().Set(CorrelationIDHeader, correlationID)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
