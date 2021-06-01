package handler

import (
	"errors"
	"net/http"

	"github.com/yjmurakami/go-kakeibo/internal/clock"
)

type middleware struct {
	service service.MiddlewareService
	jwt     Jwt
	clock   clock.Clock
}

func (m *middleware) HandleError(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			var cerr clientError
			if errors.As(err, &cerr) {
				m.handleClientError(w, cerr)
				return
			}
			m.handleServerError(w, r, err)
		}
	}
}
