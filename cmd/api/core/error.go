package core

import "errors"

var (
	ErrNoResource           = errors.New("resource does not exist")
	ErrAuthenticationFailed = errors.New("authentication failed")
)
