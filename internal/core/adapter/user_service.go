package adapter

import (
	"context"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/rs/zerolog"
)

type CreateUser func(ctx context.Context, logger zerolog.Logger, user domain.User) (*uuid.UUID, error)
