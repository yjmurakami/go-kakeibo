package handler

import (
	"errors"
	"fmt"
)

var (
	errJWTInvalid = errors.New("jwt is invalid")
	errJWTExpired = errors.New("jwt is expired")
)

type clientError struct {
	status int
	body   string
}

func (e clientError) Error() string {
	return fmt.Sprintf("status : %v, body : %v", e.status, e.body)
}
