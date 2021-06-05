// Code generated by openapi-generator; DO NOT EDIT.
// API version: 0.0.1

package openapi

type V1TransactionsPostReq struct {
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	Type       int    `json:"type" validate:"required,oneof=1 2"` // Transaction type:   * 1 - income   * 2 - expense
	CategoryId int    `json:"categoryId" validate:"min=1,max=999999999"`
	Amount     int    `json:"amount" validate:"min=1,max=999999999"`
	Note       string `json:"note" validate:"max=1000"`
}
