package handler

import (
	"net/http"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
)

type systemHandler struct {
	service service.SystemService
}

func NewSystemHandler(s service.SystemService) *systemHandler {
	return &systemHandler{
		service: s,
	}
}

func (h *systemHandler) V1HealthGet() func(w http.ResponseWriter, r *http.Request) error {
	type response struct {
		Data *openapi.V1HealthRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		oaRes, err := h.service.V1HealthGet(r.Context())
		if err != nil {
			return err
		}

		err = encodeJSON(w, response{
			Data: oaRes,
		})
		return err
	}
}
