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
	container := map[string]interface{}{
		"error": message,
	}

	err := encodeJSON(w, status, container, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusBadRequest
	clientError(w, r, status, err.Error())
}

func unauthorizedError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	clientError(w, r, status, http.StatusText(status))
}

func NotFoundError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	clientError(w, r, status, http.StatusText(status))
}

func MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMethodNotAllowed
	clientError(w, r, status, http.StatusText(status))
}

func conflictError(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the resource due to an edit conflict, please try again"
	clientError(w, r, http.StatusConflict, message)
}

func unprocessableEntityError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	status := http.StatusUnprocessableEntity
	clientError(w, r, status, errors)
}

func serverError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	status := http.StatusInternalServerError
	clientError(w, r, status, http.StatusText(status))
}
