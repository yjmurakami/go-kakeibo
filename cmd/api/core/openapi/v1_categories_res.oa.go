// Code generated by openapi-generator; DO NOT EDIT.
// API version: 0.0.1

package openapi

type V1CategoriesRes struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"` // Transaction type:   * 1 - income   * 2 - expense
	TypeName string `json:"type_name"`
	Name     string `json:"name"`
}
