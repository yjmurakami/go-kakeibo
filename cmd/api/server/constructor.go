package server

import (
	"database/sql"
	"log"

	"github.com/yjmurakami/go-kakeibo/cmd/api/handler"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
	"github.com/yjmurakami/go-kakeibo/internal/repository"
)

type handlerConfig struct {
	logger *log.Logger
	clock  clock.Clock
	db     *sql.DB
	jwt    handler.Jwt
	config *config
	repos  repository.Repositories
}

func initTransactionHandler(hc handlerConfig) handler.TransactionHandler {
	return nil
}

func initCategoryHandler(hc handlerConfig) handler.CategoryHandler {
	return nil
}

func initSystemHandler(hc handlerConfig) handler.SystemHandler {
	return handler.NewSystemHandler(
		service.NewSystemService(
			hc.db,
		),
	)
}
