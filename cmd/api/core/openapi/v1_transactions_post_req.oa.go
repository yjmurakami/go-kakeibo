// Code generated by openapi-generator; DO NOT EDIT.
// API version: 0.0.1

package openapi

type V1TransactionsPostReq struct {
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	CategoryId int    `json:"category_id" validate:"min=1,max=999999999"`
	Amount     int    `json:"amount" validate:"min=1,max=999999999"`
	Note       string `json:"note" validate:"max=1000"`
}
