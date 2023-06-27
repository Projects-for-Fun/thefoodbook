package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
