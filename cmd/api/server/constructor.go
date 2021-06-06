package server

import (
	"database/sql"
	"log"

	"github.com/yjmurakami/go-kakeibo/cmd/api/handler"
	"github.com/yjmurakami/go-kakeibo/cmd/api/queryservice"
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
	return handler.NewTransactionHandler(
		service.NewTransactionService(
			hc.db,
			queryservice.NewTransactionQueryService(),
			hc.repos,
			hc.clock,
		),
	)
}

func initCategoryHandler(hc handlerConfig) handler.CategoryHandler {
	return handler.NewCategoryHandler(
		service.NewCategoryService(
			hc.db,
			queryservice.NewCategoryQueryService(),
		),
	)
}

func initSystemHandler(hc handlerConfig) handler.SystemHandler {
	return handler.NewSystemHandler(
		service.NewSystemService(
			hc.db,
		),
	)
}
