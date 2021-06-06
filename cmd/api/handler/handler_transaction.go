package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
	"github.com/yjmurakami/go-kakeibo/internal/constant"
	"github.com/yjmurakami/go-kakeibo/internal/validator"
)

type transactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *transactionHandler {
	return &transactionHandler{
		service: s,
	}
}

func (h *transactionHandler) V1TransactionsPost() http.HandlerFunc {
	type response struct {
		Data *openapi.V1TransactionsRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()

		oaReq := &openapi.V1TransactionsPostReq{}
		err := decodeJSON(w, r, oaReq)
		if err != nil {
			v.AddError("json", err.Error())
			badRequestError(w, r, v.Errors())
			return
		}

		v.ValidateStruct(oaReq)
		if !v.Valid() {
			badRequestError(w, r, v.Errors())
			return
		}

		oaRes, err := h.service.V1TransactionsPost(r.Context(), oaReq)
		if err != nil {
			var errInvalidParameter core.ErrInvalidParameter
			if errors.As(err, &errInvalidParameter) {
				v.AddError(errInvalidParameter.Key, errInvalidParameter.Message)
				badRequestError(w, r, v.Errors())
			} else {
				serverError(w, r, err)
			}
			return
		}

		header := http.Header{}
		header.Set("Location", fmt.Sprintf("/v1/transactions/%d", oaRes.Id))
		res := response{
			Data: oaRes,
		}
		err = encodeJSON(w, http.StatusCreated, res, header)
		if err != nil {
			serverError(w, r, err)
			return
		}
	}
}

func (h *transactionHandler) V1TransactionsGet() http.HandlerFunc {
	type request struct {
		From time.Time
		To   time.Time
		core.Filter
	}
	type response struct {
		Metadata openapi.Metadata             `json:"metadata"`
		Data     []*openapi.V1TransactionsRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		req := request{
			From: readDate(qs, "from", constant.MinTime(), v),
			To:   readDate(qs, "to", constant.MaxTime(), v),
			Filter: core.Filter{
				Page:     readInt(qs, "page", 1, v),
				PageSize: readInt(qs, "page_size", 20, v),
				Sort:     readString(qs, "sort", "id"),
				SortSafelist: map[string]bool{
					"id":           true,
					"-id":          true,
					"date":         true,
					"-date":        true,
					"category_id":  true,
					"-category_id": true,
				},
			},
		}

		core.ValidateFilter(v, req.Filter)
		if !v.Valid() {
			badRequestError(w, r, v.Errors())
			return
		}

		oaRes, metadata, err := h.service.V1TransactionsGet(r.Context(), req.From, req.To, req.Filter)
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

func (h *transactionHandler) V1TransactionsTransactionIdDelete() http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := readParamId(r, "transaction_id")
		if err != nil {
			NotFoundError(w, r)
			return
		}

		err = h.service.V1TransactionsTransactionIdDelete(r.Context(), id)
		if err != nil {
			if errors.Is(err, core.ErrNoResource) {
				NotFoundError(w, r)
			} else {
				serverError(w, r, err)
			}
			return
		}

		res := response{
			Message: "successfully deleted",
		}
		err = encodeJSON(w, http.StatusOK, res, nil)
		if err != nil {
			serverError(w, r, err)
			return
		}
	}
}

func (h *transactionHandler) V1TransactionsTransactionIdGet() http.HandlerFunc {
	type response struct {
		Data *openapi.V1TransactionsRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := readParamId(r, "transaction_id")
		if err != nil {
			NotFoundError(w, r)
			return
		}

		oaRes, err := h.service.V1TransactionsTransactionIdGet(r.Context(), id)
		if err != nil {
			if errors.Is(err, core.ErrNoResource) {
				NotFoundError(w, r)
			} else {
				serverError(w, r, err)
			}
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

func (h *transactionHandler) V1TransactionsTransactionIdPatch() http.HandlerFunc {
	type response struct {
		Data *openapi.V1TransactionsRes `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()

		id, err := readParamId(r, "transaction_id")
		if err != nil {
			NotFoundError(w, r)
			return
		}

		oaReq := &openapi.V1TransactionsTransactionIdPatchReq{}
		err = decodeJSON(w, r, oaReq)
		if err != nil {
			v.AddError("json", err.Error())
			badRequestError(w, r, v.Errors())
			return
		}

		v.ValidateStruct(oaReq)
		if !v.Valid() {
			badRequestError(w, r, v.Errors())
			return
		}

		oaRes, err := h.service.V1TransactionsTransactionIdPatch(r.Context(), id, oaReq)
		if err != nil {
			if errors.Is(err, core.ErrNoResource) {
				NotFoundError(w, r)
			} else if errors.Is(err, core.ErrEditConflict) {
				conflictError(w, r)
			} else {
				serverError(w, r, err)
			}
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
