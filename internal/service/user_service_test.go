package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createNewUserRepoMock(userID *uuid.UUID, genericError error) adapter.CreateUserRepo {
	return func(ctx context.Context, user domain.User) (*uuid.UUID, error) {
		return userID, genericError
	}
}

func TestHandleCreateUser(t *testing.T) {
	tests := []struct {
		description    string
		w              webservice.Webservice
		expectedUserID *uuid.UUID
		expectedErr    error
	}{
		{
			description: "Should error when creating a new user",
			w: *webservice.NewWebservice(
				HandleCreateUserFunc(createNewUserRepoMock(nil, fmt.Errorf("error"))),
				nil,
				nil,
			),
			expectedUserID: nil,
			expectedErr:    fmt.Errorf("error"),
		},
		{
			description: "Should create a new user",
			w: *webservice.NewWebservice(
				HandleCreateUserFunc(createNewUserRepoMock(&uuid.Nil, nil)),
				nil,
				nil,
			),
			expectedUserID: &uuid.Nil,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			userID, err := tt.w.CreateUser(context.TODO(), domain.User{})

			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func validateLoginUserRepoMock(user *domain.User, err error) adapter.ValidateLoginUserRepo {
	return func(ctx context.Context, username string) (*domain.User, error) {
		return user, err
	}
}

func setUserLastLoginRepoMock(err error) adapter.SetUserLastLoginRepo {
	return func(ctx context.Context, username string) error {
		return err
	}
}

func TestHandleLoginUserFunc(t *testing.T) {

	tests := []struct {
		description  string
		w            webservice.Webservice
		expectedUser *domain.User
		expectedErr  error
	}{
		{
			description: "Should fail the validation",
			w: *webservice.NewWebservice(nil,
				HandleLoginUserFunc(validateLoginUserRepoMock(nil, fmt.Errorf("error")),
					func(password, hash string) bool { return false },
					setUserLastLoginRepoMock(nil)),
				nil),
			expectedUser: nil,
			expectedErr:  fmt.Errorf("error"),
		},
		{
			description: "Should fail when verifying the password",
			w: *webservice.NewWebservice(nil,
				HandleLoginUserFunc(validateLoginUserRepoMock(&domain.User{}, nil),
					func(password, hash string) bool { return false },
					setUserLastLoginRepoMock(fmt.Errorf("error"))),
				nil),
			expectedUser: nil,
			expectedErr:  domain.ErrInvalidUsernameOrPassword,
		},
		{
			description: "Should fail when set user as logged in",
			w: *webservice.NewWebservice(nil,
				HandleLoginUserFunc(validateLoginUserRepoMock(&domain.User{}, nil),
					func(password, hash string) bool { return true },
					setUserLastLoginRepoMock(fmt.Errorf("error"))),
				nil),
			expectedUser: nil,
			expectedErr:  fmt.Errorf("error"),
		},
		{
			description: "Should succeed login the user",
			w: *webservice.NewWebservice(nil,
				HandleLoginUserFunc(validateLoginUserRepoMock(&domain.User{}, nil),
					func(password, hash string) bool { return true },
					setUserLastLoginRepoMock(nil)),
				nil),
			expectedUser: &domain.User{},
			expectedErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			user, err := tt.w.LoginUser(context.TODO(), "username", "password")

			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
