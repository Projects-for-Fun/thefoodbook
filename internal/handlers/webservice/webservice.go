package webservice

import (
	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
)

type Webservice struct {
	CreateUser adapter.CreateUser
	LoginUser  adapter.LoginUser
}

func NewWebservice(
	createUser adapter.CreateUser,
	loginUser adapter.LoginUser,
) *Webservice {
	return &Webservice{
		CreateUser: createUser,
		LoginUser:  loginUser,
	}
}
