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
				serverError(w, r, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (m *middlewareHandler) Authenticate(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie(cookieNameToken)
		if err != nil {
			unauthorizedError(w, r)
			return
		}

		userId, err := m.jwt.DecodeToken(ck.Value)
		if err != nil {
			if errors.Is(err, errJWTInvalid) || errors.Is(err, errJWTExpired) {
				unauthorizedError(w, r)
				return
			}
			serverError(w, r, err)
			return
		}

		user, err := m.service.Authenticate(r.Context(), userId)
		if err != nil {
			if errors.Is(err, core.ErrAuthenticationFailed) {
				unauthorizedError(w, r)
				return
			}
			serverError(w, r, err)
			return
		}

		token, err := m.jwt.NewToken(userId)
		if err != nil {
			serverError(w, r, err)
			return
		}

		exp := m.clock.Now().Add(m.jwt.Expiration())
		http.SetCookie(w, newCookie(cookieNameToken, token, exp))

		ctx := r.Context()
		ctx = core.SetContextUser(ctx, *user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
