package mws

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrelationIDMiddleware(t *testing.T) {
	tests := []struct {
		description            string
		hasCorrelationIDHeader bool
		correlationID          string
	}{
		{
			description:            "Request without Correlation Id should always create and return one",
			hasCorrelationIDHeader: false,
			correlationID:          "", // Not needed
		},

		{
			description:            "Request with Correlation Id should add it to context and return it",
			hasCorrelationIDHeader: true,
			correlationID:          "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://testing", nil)

			if tt.hasCorrelationIDHeader {
				req.Header.Set(CorrelationIDHeader, tt.correlationID)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				correlationID := r.Context().Value(RequestIDKey)

				if correlationID == nil {
					t.Errorf("CorrelationID should never be null.")
				}

				if tt.hasCorrelationIDHeader {
					assert.Equal(t, tt.correlationID, r.Header.Get(CorrelationIDHeader), "The correlation ids should match.")
				}
			})

			middlewareToTest := CorrelationID(nextHandler)
			middlewareToTest.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
