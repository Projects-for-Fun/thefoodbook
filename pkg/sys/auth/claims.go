package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// TokenKey The key used in context for claims
var TokenKey = "TokenKey"

func AttachToken(ctx context.Context, tknStr *jwt.Token) context.Context {
	return context.WithValue(ctx, TokenKey, tknStr)
}

func GetToken(ctx context.Context) *jwt.Token {
	tknStr, ok := ctx.Value(TokenKey).(*jwt.Token)

	if !ok {
		return nil
	}

	return tknStr
}
