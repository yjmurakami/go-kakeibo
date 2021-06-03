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

func (h *systemHandler) V1HealthGet() http.HandlerFunc {
	type response struct {
		Data *openapi.V1HealthRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		oaRes, err := h.service.V1HealthGet(r.Context())
		if err != nil {
			serverError(w, r, err)
			return
		}

		res := response{
			Data: oaRes,
		}
		err = encodeJSON(w, http.StatusOK, res, nil)
		if err != nil {
			serverError(w, r, err)
			return
		}
	}
}
