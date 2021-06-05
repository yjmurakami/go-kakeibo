package core

import (
	"errors"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrEditConflict         = errors.New("edit conflict")
	ErrNoResource           = errors.New("resource does not exist")
)

type ErrInvalidParameter struct {
	Key     string
	Message string
}

func (e ErrInvalidParameter) Error() string {
	return e.Message
}
