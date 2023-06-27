package service

import (
	"context"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

func HandleCreateTokenFunc(jwtKey []byte) adapter.CreateToken {
	return func(ctx context.Context, user domain.User) (time.Time, string, error) {
		expirationTime := time.Now().Add(30 * time.Second)

		claims := &domain.Claims{
			UserID:   user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   user.ID.String(),
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		return expirationTime, tokenString, err
	}
}
