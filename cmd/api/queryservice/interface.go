package queryservice

import (
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/dto"
	"github.com/yjmurakami/go-kakeibo/internal/database"
)

type TransactionQueryService interface {
	SelectTransactionByID(db database.DB, id int) (*dto.Transaction, error)
}
