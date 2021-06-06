package queryservice

import (
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/dto"
	"github.com/yjmurakami/go-kakeibo/internal/database"
)

type transactionQueryService struct{}

func NewTransactionQueryService() *transactionQueryService {
	return &transactionQueryService{}
}

func (q *transactionQueryService) SelectTransactionByID(db database.DB, id int) (*dto.Transaction, error) {
	return dto.SelectTransactionById(db, id)
}

func (q *transactionQueryService) SelectTransactions(db database.DB, from time.Time, to time.Time, filter core.Filter) ([]*dto.Transaction, core.Metadata, error) {
	return dto.SelectTransactions(db, from, to, filter)
}
