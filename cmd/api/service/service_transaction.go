package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/dto"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/cmd/api/queryservice"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
	"github.com/yjmurakami/go-kakeibo/internal/repository"
)

type transactionService struct {
	db    *sql.DB
	qs    queryservice.TransactionQueryService
	repos repository.Repositories
	clock clock.Clock
}

func NewTransactionService(db *sql.DB, qs queryservice.TransactionQueryService, repos repository.Repositories, clock clock.Clock) *transactionService {
	return &transactionService{
		db:    db,
		qs:    qs,
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

	dtoTransaction, err := s.qs.SelectTransactionByID(s.db, transaction.ID)
	if err != nil {
		return nil, err
	}

	return s.mapTransactionToOpenAPI(dtoTransaction), nil
}

func (s *transactionService) V1TransactionsGet(ctx context.Context, from time.Time, to time.Time, filter core.Filter) ([]*openapi.V1TransactionsRes, openapi.Metadata, error) {
	transactions, metadata, err := s.qs.SelectTransactions(s.db, from, to, filter)
	if err != nil {
		return nil, openapi.Metadata{}, err
	}

	// 権限チェック

	oaRes := []*openapi.V1TransactionsRes{}
	for _, t := range transactions {
		oaRes = append(oaRes, s.mapTransactionToOpenAPI(t))
	}

	return oaRes, metadata, nil
}

func (s *transactionService) V1TransactionsTransactionIdDelete(ctx context.Context, transactionId int) error {
	transaction, err := s.repos.Transaction.SelectByID(s.db, transactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.ErrNoResource
		}
		return err
	}

	// TODO 権限チェック

	err = s.repos.Transaction.Delete(s.db, transaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.ErrNoResource
		}
		return err
	}

	return nil
}

func (s *transactionService) V1TransactionsTransactionIdGet(ctx context.Context, transactionId int) (*openapi.V1TransactionsRes, error) {
	transaction, err := s.qs.SelectTransactionByID(s.db, transactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNoResource
		}
		return nil, err
	}

	// TODO 権限チェック

	return s.mapTransactionToOpenAPI(transaction), nil
}

func (s *transactionService) V1TransactionsTransactionIdPatch(ctx context.Context, transactionId int, oaReq *openapi.V1TransactionsTransactionIdPatchReq) (*openapi.V1TransactionsRes, error) {
	transaction, err := s.repos.Transaction.SelectByID(s.db, transactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNoResource
		}
		return nil, err
	}

	// TODO 権限チェック

	if oaReq.Date != nil {
		dt, err := time.Parse(openapi.DateFormat, *oaReq.Date)
		if err != nil {
			return nil, err
		}

		transaction.Date = dt
	}

	if oaReq.CategoryId != nil {
		category, err := s.repos.Category.SelectByID(s.db, *oaReq.CategoryId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, core.ErrInvalidParameter{
					Key:     "categoryId",
					Message: "the value is invalid",
				}
			}
			return nil, err
		}

		transaction.CategoryID = category.ID
	}
	if oaReq.Amount != nil {
		transaction.Amount = *oaReq.Amount
	}
	if oaReq.Note != nil {
		transaction.Note = *oaReq.Note
	}
	transaction.ModifiedAt = s.clock.Now()

	err = s.repos.Transaction.Update(s.db, transaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrEditConflict
		}
		return nil, err
	}

	dtoTransaction, err := s.qs.SelectTransactionByID(s.db, transaction.ID)
	if err != nil {
		return nil, err
	}

	return s.mapTransactionToOpenAPI(dtoTransaction), nil
}

func (s *transactionService) mapTransactionToOpenAPI(t *dto.Transaction) *openapi.V1TransactionsRes {
	oa := &openapi.V1TransactionsRes{
		Id:           t.ID,
		Date:         t.Date.Format(openapi.DateFormat),
		Type:         t.CategoryType,
		TypeName:     openapi.CategoryTypes[t.CategoryType],
		CategoryId:   t.CategoryID,
		CategoryName: t.CategoryName,
		Amount:       t.Amount,
		Note:         t.Note,
	}
	return oa
}
