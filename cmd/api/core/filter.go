package core

import (
	"math"
	"strings"

	"github.com/yjmurakami/go-kakeibo/internal/validator"
)

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist map[string]bool
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func ValidateFilter(v *validator.Validator, f Filter) {
	v.Check(f.Page >= 1, "page", "must be more than or equal to 1")
	v.Check(f.Page <= 1_000_000, "page", "must be less than or equal to 1 million")
	v.Check(f.PageSize >= 1, "page_size", "must be more than or equal to 1")
	v.Check(f.PageSize <= 100, "page_size", "must be less than or equal to 100")
	v.Check(validator.InMap(f.Sort, f.SortSafelist), "sort", "invalid sort value")
}

func (f Filter) SortColumn() string {
	_, ok := f.SortSafelist[f.Sort]
	if ok {
		return strings.TrimPrefix(f.Sort, "-")
	}

	return ""
}

func (f Filter) SortDirection() string {
	_, ok := f.SortSafelist[f.Sort]
	if ok {
		if strings.HasPrefix(f.Sort, "-") {
			return "DESC"
		}
		return "ASC"
	}

	return ""
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func CalculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
