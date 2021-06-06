package queryservice

import (
	"context"
	"fmt"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type categoryQueryService struct{}

func NewCategoryQueryService() *categoryQueryService {
	return &categoryQueryService{}
}

// Generated from SelectCategorys.sql
func (q *categoryQueryService) SelectCategories(db database.DB, categoryType int, filter core.Filter) ([]*entity.Category, openapi.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) OVER() total_records
			, id
			, type
			, name
			, created_at
			, modified_at
			, version
		FROM
			kakeibo.categories
		WHERE
			(type = ? OR ? = 0)
		ORDER BY
		  %s %s, id
		LIMIT ? OFFSET ?
	`, filter.SortColumn(), filter.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, categoryType, categoryType, filter.Limit(), filter.Offset())
	if err != nil {
		return nil, openapi.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	s := []*entity.Category{}
	for rows.Next() {
		d := entity.Category{}
		err = rows.Scan(&totalRecords, &d.ID, &d.Type, &d.Name, &d.CreatedAt, &d.ModifiedAt, &d.Version)
		if err != nil {
			return nil, openapi.Metadata{}, err
		}
		s = append(s, &d)
	}

	err = rows.Err()
	if err != nil {
		return nil, openapi.Metadata{}, err
	}

	metadata := openapi.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)
	return s, metadata, nil
}
