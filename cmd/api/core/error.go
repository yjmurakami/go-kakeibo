package core

import (
	"errors"
)

var (
	ErrNoResource           = errors.New("resource does not exist")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

type ErrInvalidParameter struct {
	Key     string
	Message string
}

func (e ErrInvalidParameter) Error() string {
	return e.Message
}
