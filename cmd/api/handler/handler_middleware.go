package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
)

type middlewareHandler struct {
	service service.MiddlewareService
	jwt     Jwt
	clock   clock.Clock
}

func NewMiddlewareHandler(s service.MiddlewareService, j Jwt, c clock.Clock) *middlewareHandler {
	return &middlewareHandler{
		service: s,
		jwt:     j,
		clock:   c,
	}
}

func (m *middlewareHandler) RecoverPanic(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.handleServerError(w, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (m *middlewareHandler) Authenticate(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie(cookieNameToken)
		if err != nil {
			m.handleClientError(w, newUnauthorizedError())
			return
		}

		userId, err := m.jwt.DecodeToken(ck.Value)
		if err != nil {
			if errors.Is(err, errJWTInvalid) || errors.Is(err, errJWTExpired) {
				m.handleClientError(w, newUnauthorizedError())
				return
			}
			m.handleServerError(w, err)
			return
		}

		user, err := m.service.Authenticate(r.Context(), userId)
		if err != nil {
			if errors.Is(err, core.ErrAuthenticationFailed) {
				m.handleClientError(w, newUnauthorizedError())
				return
			}
			m.handleServerError(w, err)
			return
		}

		token, err := m.jwt.NewToken(userId)
		if err != nil {
			m.handleServerError(w, err)
			return
		}

		exp := m.clock.Now().Add(m.jwt.Expiration())
		http.SetCookie(w, newCookie(cookieNameToken, token, exp))

		ctx := r.Context()
		ctx = core.SetContextUser(ctx, *user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *middlewareHandler) HandleError(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			var clientErr clientError
			if errors.As(err, &clientErr) {
				m.handleClientError(w, clientErr)
				return
			}
			m.handleServerError(w, err)
		}
	}
}

func (m *middlewareHandler) handleServerError(w http.ResponseWriter, serverErr error) {

	// TODO エラーログ

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, http.StatusText(http.StatusInternalServerError)) // TODO JSONに変更
}

func (m *middlewareHandler) handleClientError(w http.ResponseWriter, clientErr clientError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(clientErr.status)
	encodeJSON(w, clientErr.body)
}
