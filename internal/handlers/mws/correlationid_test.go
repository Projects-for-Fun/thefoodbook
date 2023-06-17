package mws

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrelationIDMiddleware(t *testing.T) {
	tests := []struct {
		description            string
		hasCorrelationIdHeader bool
		correlationId          string
	}{
		{
			description:            "Request without Correlation Id should always create and return one",
			hasCorrelationIdHeader: false,
			correlationId:          "", // Not needed
		},

		{
			description:            "Request with Correlation Id should add it to context and return it",
			hasCorrelationIdHeader: true,
			correlationId:          "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://testing", nil)

			if tt.hasCorrelationIdHeader {
				req.Header.Set(CorrelationIDHeader, tt.correlationId)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				correlationId := r.Context().Value(RequestIDKey)

				if correlationId == nil {
					t.Errorf("CorrelationID should never be null.")
				}

				if tt.hasCorrelationIdHeader {
					assert.Equal(t, tt.correlationId, r.Header.Get(CorrelationIDHeader), "The correlation ids should match.")

				}
			})

			middlewareToTest := CorrelationID(nextHandler)
			middlewareToTest.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
