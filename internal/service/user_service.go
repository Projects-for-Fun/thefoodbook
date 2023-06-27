package service

import (
	"context"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

	"github.com/google/uuid"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

func CreateUserServiceFunc(create adapter.CreateUserRepo) adapter.CreateUserService {
	return func(ctx context.Context, user domain.User) (*uuid.UUID, error) {
		logger := logging.GetLogger(ctx)

		userID, err := create(ctx, user)

		if err == nil {
			logger.Info().
				Str("user-id", userID.String()).
				Str("username", user.Username).
				Msg("Created new user")
		}

		return userID, err
	}
}

func LoginUserServiceFunc(getUserByUsername adapter.GetUserByUsernameRepo,
	verifyPassword func(password, hash string) bool,
	setUserLastLogin adapter.SetUserLastLoginRepo) adapter.LoginUserService {
	return func(ctx context.Context, username, password string) (*domain.User, error) {
		logger := logging.GetLogger(ctx)

		userLogged, err := getUserByUsername(ctx, username)

		if err != nil {
			logger.Info().
				AnErr("error", err).
				Str("username", username).
				Msg("Couldn't log user")
			return nil, err
		}

		if !verifyPassword(password, userLogged.Password) {
			logger.Info().
				Str("username", username).
				Msg("Couldn't log user - invalid username or password")
			return nil, domain.ErrInvalidUsernameOrPassword
		}

		err = setUserLastLogin(ctx, username)
		if err != nil {
			return nil, err
		}

		return userLogged, err
	}
}
