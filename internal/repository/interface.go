package repository

import (
	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type UserRepository interface {
	SelectByID(db database.DB, id int) (*entity.User, error)
}
