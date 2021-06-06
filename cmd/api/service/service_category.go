package service

import (
	"context"
	"database/sql"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/cmd/api/queryservice"
)

type categoryService struct {
	db *sql.DB
	qs queryservice.CategoryQueryService
}

func NewCategoryService(db *sql.DB, qs queryservice.CategoryQueryService) *categoryService {
	return &categoryService{
		db: db,
		qs: qs,
	}
}

func (s *categoryService) V1CategoriesGet(ctx context.Context, categoryType int, filter core.Filter) ([]*openapi.V1CategoriesRes, openapi.Metadata, error) {
	categories, metadata, err := s.qs.SelectCategories(s.db, categoryType, filter)
	if err != nil {
		return nil, openapi.Metadata{}, err
	}

	oaRes := []*openapi.V1CategoriesRes{}
	for _, c := range categories {
		oa := &openapi.V1CategoriesRes{
			Id:   c.ID,
			Type: c.Type,
			Name: c.Name,
		}
		oaRes = append(oaRes, oa)
	}

	return oaRes, metadata, nil
}
