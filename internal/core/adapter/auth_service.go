package adapter

import (
	"context"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

type CreateToken func(ctx context.Context, user domain.User) (time.Time, string, error)
