package service

import (
	"context"
	"database/sql"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
)

type systemService struct {
	db *sql.DB
}

func NewSystemService(db *sql.DB) *systemService {
	return &systemService{
		db: db,
	}
}

func (s *systemService) V1HealthGet(ctx context.Context) (*openapi.V1HealthRes, error) {
	oaRes := &openapi.V1HealthRes{
		Status:  "available",
		Version: "0.0.1", // TODO 変数
	}

	err := s.db.Ping()
	if err != nil {
		oaRes.Status = "unavailable"
	}

	return oaRes, nil
}
