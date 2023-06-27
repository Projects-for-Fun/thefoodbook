package auth

import (
	"time"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

func CreateTokenForUser(user domain.User, jwtKey []byte) (time.Time, string, error) {
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

	return createToken(claims, expirationTime, jwtKey)
}

func CreateTokenFromExistingClaims(claims *domain.Claims, jwtKey []byte) (time.Time, string, error) {
	expirationTime := time.Now().Add(30 * time.Second)

	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)

	return createToken(claims, expirationTime, jwtKey)
}

func createToken(claims *domain.Claims, expirationTime time.Time, jwtKey []byte) (time.Time, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return expirationTime, tokenString, err
}
