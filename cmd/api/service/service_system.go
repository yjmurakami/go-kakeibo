package service

import (
	"context"
	"database/sql"
)

type systemService struct {
	db *sql.DB
}

func NewSystemService(db *sql.DB) *systemService {
	return &systemService{
		db: db,
	}
}

func (s *systemService) V1HealthGet(ctx context.Context) error {
	return s.db.Ping()
}
