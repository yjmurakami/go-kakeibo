package handler

import (
	"errors"
	"fmt"
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

func (m *middleware) handleServerError(w http.ResponseWriter, serverErr error) {

	// TODO エラーログ

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, http.StatusText(http.StatusInternalServerError)) // TODO JSONに変更
}

func (m *middleware) handleClientError(w http.ResponseWriter, clientErr clientError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(clientErr.status)
	encodeJSON(w, clientErr.body)
}
