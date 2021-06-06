package dto

import (
	"time"
)

// Generated from SelectTransactionById.sql
type Transaction struct {
	ID           int       // id
	UserID       int       // user_id
	Date         time.Time // date
	Amount       int       // amount
	Note         string    // note
	CreatedAt    time.Time // created_at
	ModifiedAt   time.Time // modified_at
	Version      int       // version
	CategoryID   int       // category_id
	CategoryType int       // category_type
	CategoryName string    // category_name
}
