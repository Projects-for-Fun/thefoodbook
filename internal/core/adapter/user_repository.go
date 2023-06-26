package adapter

import (
	"context"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

type CreateUserRepo func(ctx context.Context, user domain.User) (*uuid.UUID, error)

type ValidateLoginUserRepo func(ctx context.Context, username string) (*domain.User, error)

type SetUserLastLoginRepo func(ctx context.Context, username string) error
