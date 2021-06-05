package repository

import (
	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

type CategoryRepository interface {
	SelectAll(db database.DB) ([]*entity.Category, error)
	SelectByID(db database.DB, id int) (*entity.Category, error)
	Insert(db database.DB, e *entity.Category) error
	Update(db database.DB, e *entity.Category) error
	Delete(db database.DB, e *entity.Category) error
}

type UserRepository interface {
	SelectAll(db database.DB) ([]*entity.User, error)
	SelectByID(db database.DB, id int) (*entity.User, error)
	Insert(db database.DB, e *entity.User) error
	Update(db database.DB, e *entity.User) error
	Delete(db database.DB, e *entity.User) error
}

type TransactionRepository interface {
	SelectAll(db database.DB) ([]*entity.Transaction, error)
	SelectByID(db database.DB, id int) (*entity.Transaction, error)
	Insert(db database.DB, e *entity.Transaction) error
	Update(db database.DB, e *entity.Transaction) error
	Delete(db database.DB, e *entity.Transaction) error
}
