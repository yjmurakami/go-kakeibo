package core

import (
	"context"

	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type contextKey int

const (
	contextKeyUser contextKey = iota
)

func SetContextUser(ctx context.Context, u entity.User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}
