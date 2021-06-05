package queryservice

import (
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
