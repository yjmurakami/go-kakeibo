package handler

import "net/http"

type MiddlewareHandler interface {
	RecoverPanic(next http.Handler) http.HandlerFunc
	HandleError(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc
	Authenticate(next http.Handler) http.HandlerFunc
}
