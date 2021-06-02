package service

import (
	"context"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type MiddlewareService interface {
	Authenticate(ctx context.Context, userID int) (*entity.User, error)
}

type SystemService interface {
	V1HealthGet(ctx context.Context) (*openapi.V1HealthRes, error)
}
