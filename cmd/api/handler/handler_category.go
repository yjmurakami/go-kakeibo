package handler

import (
	"net/http"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
	"github.com/yjmurakami/go-kakeibo/internal/validator"
)

type categoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *categoryHandler {
	return &categoryHandler{
		service: s,
	}
}

func (h *categoryHandler) V1CategoriesGet() http.HandlerFunc {
	type request struct {
		Type int
		core.Filter
	}
	type response struct {
		Metadata openapi.Metadata           `json:"metadata"`
		Data     []*openapi.V1CategoriesRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		req := request{
			Type: readInt(qs, "type", 0, v),
			Filter: core.Filter{
				Page:     readInt(qs, "page", 1, v),
				PageSize: readInt(qs, "page_size", 20, v),
				Sort:     readString(qs, "sort", "id"),
				SortSafelist: map[string]bool{
					"id":   true,
					"type": true,
					"name": true,
				},
			},
		}

		core.ValidateFilter(v, req.Filter)
		if !v.Valid() {
			badRequestError(w, r, v.Errors())
			return
		}

		oaRes, metadata, err := h.service.V1CategoriesGet(r.Context(), req.Type, req.Filter)
		if err != nil {
			serverError(w, r, err)
			return
		}

		res := response{
			Metadata: metadata,
			Data:     oaRes,
		}
		err = encodeJSON(w, http.StatusOK, res, nil)
		if err != nil {
			serverError(w, r, err)
			return
		}
	}
}
