package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func CreateNewUserRepo(userID *uuid.UUID, genericError error) adapter.CreateUserRepo {
	return func(ctx context.Context, user domain.User) (*uuid.UUID, error) {
		return userID, genericError
	}
}

func TestHandleCreateUserError(t *testing.T) {
	var genericError = errors.New("error")

	w := webservice.NewWebservice(
		HandleCreateUserFunc(CreateNewUserRepo(nil, genericError)),
		nil,
		nil,
	)
	userID, err := w.CreateUser(context.TODO(), domain.User{})

	assert.Nil(t, userID, "User id should be null")
	assert.Equal(t, genericError, err)
}

func TestHandleCreateUser(t *testing.T) {
	userID := uuid.New()

	w := webservice.NewWebservice(
		HandleCreateUserFunc(CreateNewUserRepo(&userID, nil)),
		nil,
		nil,
	)
	returnedUserID, err := w.CreateUser(context.TODO(), domain.User{})

	assert.Equal(t, &userID, returnedUserID)
	assert.Nil(t, err)
}
