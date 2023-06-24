package tracing

import (
	"context"

	"github.com/google/uuid"
)

// CorrelationIDKey The key used in context for correlation id
var CorrelationIDKey = "CorrelationKey"

var CorrelationIDTag = "correlation-id"

var NullCorrelationID = "00000000-0000-0000-0000-000000000000"

func GenerateCorrelationID() string {
	return uuid.NewString()
}

func AttachCorrelationID(ctx context.Context, corrID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, corrID)
}

func GetCorrelationID(ctx context.Context) string {
	corrID, ok := ctx.Value(CorrelationIDKey).(string)

	if !ok {
		return NullCorrelationID
	}

	return corrID
}
