package domain

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
var ErrUnauthorized = errors.New("unauthorized")
