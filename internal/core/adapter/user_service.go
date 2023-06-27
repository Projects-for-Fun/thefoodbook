package adapter

import (
	"context"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

type CreateUserService func(ctx context.Context, user domain.User) (*uuid.UUID, error)

type LoginUserService func(ctx context.Context, username, password string) (*domain.User, error)
