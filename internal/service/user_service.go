package service

import (
	"context"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

func HandleCreateUserFunc(create adapter.CreateUserRepo) adapter.CreateUser {
	return func(ctx context.Context, user domain.User) (*uuid.UUID, error) {
		logger := logging.GetLogger(ctx)

		userID, err := create(ctx, user)

		if err == nil {
			logger.Info().Str("user-id", userID.String()).Msg("Created new user")
		}

		return userID, err
	}
}
