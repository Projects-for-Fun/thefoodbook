package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/rs/zerolog"
)

func HandleCreateUserFunc(create adapter.CreateUserRepo) adapter.CreateUser {
	return func(ctx context.Context, logger zerolog.Logger, user domain.User) (*uuid.UUID, error) {
		return create(ctx, user)
	}
}
