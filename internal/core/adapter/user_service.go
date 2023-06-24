package adapter

import (
	"context"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

type CreateUser func(ctx context.Context, user domain.User) (*uuid.UUID, error)
