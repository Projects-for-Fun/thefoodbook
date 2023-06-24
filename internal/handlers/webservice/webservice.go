package webservice

import (
	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
)

type Webservice struct {
	CreateUser  adapter.CreateUser
	LoginUser   adapter.LoginUser
	CreateToken adapter.CreateToken
}

func NewWebservice(
	createUser adapter.CreateUser,
	loginUser adapter.LoginUser,
	createToken adapter.CreateToken,
) *Webservice {
	return &Webservice{
		CreateUser:  createUser,
		LoginUser:   loginUser,
		CreateToken: createToken,
	}
}
