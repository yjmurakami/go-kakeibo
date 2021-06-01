package server

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/yjmurakami/go-kakeibo/cmd/api/handler"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
)

type handlerConfig struct {
	timer     clock.Clock
	db        *sql.DB
	validator *validator.Validate
	jwt       handler.Jwt
	config    *config
	container container
}

type container struct{}

func newContainer() container {
	c := container{}

	return c
}

func initIncomeHandler(hc handlerConfig) handler.IncomeHandler {
	return nil
}

func initCategoryHandler(hc handlerConfig) handler.CategoryHandler {
	return nil
}

func initSystemHandler(hc handlerConfig) handler.SystemHandler {
	return nil
}
