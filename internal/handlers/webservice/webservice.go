package webservice

import (
	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
)

type Webservice struct {
	CreateUser adapter.CreateUser
}

func NewWebservice(createUser adapter.CreateUser) *Webservice {
	return &Webservice{
		CreateUser: createUser,
	}
}
