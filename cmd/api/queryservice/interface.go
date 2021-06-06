package queryservice

import (
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/dto"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type CategoryQueryService interface {
	SelectCategories(db database.DB, categoryType int, filter core.Filter) ([]*entity.Category, openapi.Metadata, error)
}

type TransactionQueryService interface {
	SelectTransactionByID(db database.DB, id int) (*dto.Transaction, error)
	SelectTransactions(db database.DB, from time.Time, to time.Time, filter core.Filter) ([]*dto.Transaction, openapi.Metadata, error)
}
