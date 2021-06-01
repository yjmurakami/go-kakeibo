package service

import "github.com/yjmurakami/go-kakeibo/internal/entity"

type MiddlewareService interface {
	Authenticate(userID int) (*entity.User, error)
}
