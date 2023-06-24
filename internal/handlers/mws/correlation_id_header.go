package mws

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/tracing"
)

var CorrelationIDHeader = "X-Correlation-Id"

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := r.Header.Get(CorrelationIDHeader)
		if corrID == "" {
			corrID = tracing.GenerateCorrelationID()
			r.Header.Set(CorrelationIDHeader, corrID)
		}

		w.Header().Set(CorrelationIDHeader, corrID)

		next.ServeHTTP(w, r.WithContext(tracing.AttachCorrelationID(r.Context(), corrID)))
	})
}
