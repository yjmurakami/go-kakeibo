package server

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/yjmurakami/go-kakeibo/cmd/api/handler"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
	"github.com/yjmurakami/go-kakeibo/internal/repository"
)

type handlerConfig struct {
	clock     clock.Clock
	db        *sql.DB
	validator *validator.Validate
	jwt       handler.Jwt
	config    *config
	container container
}

type container struct {
	userRepository repository.UserRepository
}

func newContainer() container {
	c := container{
		userRepository: repository.NewUserRepository(),
	}

	return c
}

func initIncomeHandler(hc handlerConfig) handler.IncomeHandler {
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
