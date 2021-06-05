package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
	"github.com/yjmurakami/go-kakeibo/internal/repository"
)

type transactionService struct {
	db    *sql.DB
	repos repository.Repositories
	clock clock.Clock
}

func NewTransactionService(db *sql.DB, repos repository.Repositories, clock clock.Clock) *transactionService {
	return &transactionService{
		db:    db,
		repos: repos,
		clock: clock,
	}
}

func (s *transactionService) V1TransactionsPost(ctx context.Context, oaReq *openapi.V1TransactionsPostReq) (*openapi.V1TransactionsRes, error) {
	category, err := s.repos.Category.SelectByID(s.db, oaReq.CategoryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrInvalidParameter{
				Key:     "categoryId",
				Message: "the value is invalid",
			}
		}
		return nil, err
	}

	dt, err := time.Parse(openapi.DateFormat, oaReq.Date)
	if err != nil {
		return nil, err
	}

	now := s.clock.Now()
	transaction := &entity.Transaction{
		UserID:     1, // TODO Context ログインユーザー
		Date:       dt,
		CategoryID: category.ID,
		Amount:     oaReq.Amount,
		Note:       oaReq.Note,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	err = s.repos.Transaction.Insert(s.db, transaction)
	if err != nil {
		return nil, err
	}

	oaRes := &openapi.V1TransactionsRes{
		Id:         transaction.ID,
		Date:       transaction.Date.Format(openapi.DateFormat),
		Type:       category.Type,
		CategoryId: transaction.CategoryID,
		Amount:     transaction.Amount,
		Note:       transaction.Note,
	}
	return oaRes, nil
}

func (s *transactionService) V1TransactionsGet(ctx context.Context) ([]*openapi.V1TransactionsRes, error) {
	panic("not implemented") // TODO: Implement
}

func (s *transactionService) V1TransactionsTransactionIdDelete(ctx context.Context, transactionId int) error {
	panic("not implemented") // TODO: Implement
}

func (s *transactionService) V1TransactionsTransactionIdGet(ctx context.Context, transactionId int) (*openapi.V1TransactionsRes, error) {
	panic("not implemented") // TODO: Implement
}

func (s *transactionService) V1TransactionsTransactionIdPut(ctx context.Context, transactionId int, oaReq *openapi.V1TransactionsTransactionIdPutReq) (*openapi.V1TransactionsRes, error) {
	panic("not implemented") // TODO: Implement
}
