package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
	"github.com/yjmurakami/go-kakeibo/internal/repository"
)

type middlewareService struct {
	db    *sql.DB
	repos repository.Repositories
}

func NewMiddlewareService(db *sql.DB, repos repository.Repositories) *middlewareService {
	return &middlewareService{
		db:    db,
		repos: repos,
	}
}

func (m *middlewareService) Authenticate(ctx context.Context, userID int) (*entity.User, error) {
	u, err := m.repos.User.SelectByID(m.db, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrAuthenticationFailed
		}
		return nil, err
	}

	return u, nil
}
