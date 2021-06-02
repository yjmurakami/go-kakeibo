package service

import (
	"context"

	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type MiddlewareService interface {
	Authenticate(ctx context.Context, userID int) (*entity.User, error)
}
