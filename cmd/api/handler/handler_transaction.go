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
				return
			}
			serverError(w, r, err)
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
	panic("not implemented") // TODO: Implement
}

func (h *transactionHandler) V1TransactionsTransactionIdGet() http.HandlerFunc {
	panic("not implemented") // TODO: Implement
}

func (h *transactionHandler) V1TransactionsTransactionIdPut() http.HandlerFunc {
	panic("not implemented") // TODO: Implement
}
