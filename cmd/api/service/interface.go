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

type TransactionService interface {
	V1TransactionsGet(ctx context.Context) ([]*openapi.V1TransactionsRes, error)
	V1TransactionsPost(ctx context.Context, oaReq *openapi.V1TransactionsPostReq) (*openapi.V1TransactionsRes, error)
	V1TransactionsTransactionIdDelete(ctx context.Context, transactionId int) error
	V1TransactionsTransactionIdGet(ctx context.Context, transactionId int) (*openapi.V1TransactionsRes, error)
	V1TransactionsTransactionIdPatch(ctx context.Context, transactionId int, oaReq *openapi.V1TransactionsTransactionIdPatchReq) (*openapi.V1TransactionsRes, error)
}
