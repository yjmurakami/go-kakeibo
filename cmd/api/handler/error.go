package handler

import (
	"errors"
	"log"
	"net/http"
)

var (
	errJWTInvalid = errors.New("jwt is invalid")
	errJWTExpired = errors.New("jwt is expired")
)

// TODO
func logError(r *http.Request, err error) {
	log.Println(err)
}

func clientError(w http.ResponseWriter, r *http.Request, status int, message interface{}) {

	// TODO OpenAPI定義
	errContainer := map[string]interface{}{
		"error": message,
	}

	err := encodeJSON(w, status, errContainer, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func unauthorizedError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	clientError(w, r, status, http.StatusText(status))
}

func notFoundError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	clientError(w, r, status, http.StatusText(status))
}

func methodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMethodNotAllowed
	clientError(w, r, status, http.StatusText(status))
}

func serverError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	status := http.StatusInternalServerError
	clientError(w, r, status, http.StatusText(status))
}
