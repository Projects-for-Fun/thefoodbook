package webservice

import (
	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
)

type Webservice struct {
	CreateUser adapter.CreateUserService
	LoginUser  adapter.LoginUserService
}

func NewWebservice(
	createUser adapter.CreateUserService,
	loginUser adapter.LoginUserService,
) *Webservice {
	return &Webservice{
		CreateUser: createUser,
		LoginUser:  loginUser,
	}
}
