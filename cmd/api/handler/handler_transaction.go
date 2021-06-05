package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/cmd/api/core/openapi"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
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
		oaReq := &openapi.V1TransactionsPostReq{}
		err := decodeJSON(w, r, oaReq)
		if err != nil {
			badRequestError(w, r, err)
			return
		}

		v := validator.New()
		v.ValidateStruct(oaReq)
		if !v.Valid() {
			unprocessableEntityError(w, r, v.Errors())
			return
		}

		oaRes, err := h.service.V1TransactionsPost(r.Context(), oaReq)
		if err != nil {
			var errInvalidParameter core.ErrInvalidParameter
			if errors.As(err, &errInvalidParameter) {
				v.AddError(errInvalidParameter.Key, errInvalidParameter.Message)
				unprocessableEntityError(w, r, v.Errors())
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
	panic("not implemented") // TODO: Implement
}

func (h *transactionHandler) V1TransactionsTransactionIdDelete() http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getParamId(r, "transactionId")
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
		id, err := getParamId(r, "transactionId")
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
		id, err := getParamId(r, "transactionId")
		if err != nil {
			NotFoundError(w, r)
			return
		}

		oaReq := &openapi.V1TransactionsTransactionIdPatchReq{}
		err = decodeJSON(w, r, oaReq)
		if err != nil {
			badRequestError(w, r, err)
			return
		}

		v := validator.New()
		v.ValidateStruct(oaReq)
		if !v.Valid() {
			unprocessableEntityError(w, r, v.Errors())
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
