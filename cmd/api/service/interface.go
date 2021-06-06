package service

import (
	"context"
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type CategoryService interface {
	V1CategoriesGet(ctx context.Context, categoryType int, filter core.Filter) ([]*openapi.V1CategoriesRes, openapi.Metadata, error)
}

type MiddlewareService interface {
	Authenticate(ctx context.Context, userID int) (*entity.User, error)
}

type SystemService interface {
	V1HealthGet(ctx context.Context) (*openapi.V1HealthRes, error)
}

type TransactionService interface {
	V1TransactionsGet(ctx context.Context, from time.Time, to time.Time, filter core.Filter) ([]*openapi.V1TransactionsRes, openapi.Metadata, error)
	V1TransactionsPost(ctx context.Context, oaReq *openapi.V1TransactionsPostReq) (*openapi.V1TransactionsRes, error)
	V1TransactionsTransactionIdDelete(ctx context.Context, transactionId int) error
	V1TransactionsTransactionIdGet(ctx context.Context, transactionId int) (*openapi.V1TransactionsRes, error)
	V1TransactionsTransactionIdPatch(ctx context.Context, transactionId int, oaReq *openapi.V1TransactionsTransactionIdPatchReq) (*openapi.V1TransactionsRes, error)
}
