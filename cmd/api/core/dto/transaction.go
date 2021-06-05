package dto

import (
	"context"
	"fmt"
	"time"

	"github.com/yjmurakami/go-kakeibo/internal/database"
)

// Generated from Transaction1.sql
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

func SelectTransactionById(db database.DB, id int) (*Transaction, error) {
	query := fmt.Sprintf(`
		%s
		WHERE t.id = ?
	`, getTransactionQuery())

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	d := Transaction{}
	err := db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.UserID, &d.Date, &d.Amount, &d.Note, &d.CreatedAt, &d.ModifiedAt, &d.Version, &d.CategoryID, &d.CategoryType, &d.CategoryName)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Generated from Transaction1.sql
func getTransactionQuery() string {
	query := `
		SELECT
		  t.id
		  , t.user_id
		  , t.date
		  , t.amount
		  , t.note
		  , t.created_at
		  , t.modified_at
		  , t.version
		  , t.category_id
		  , c.type category_type
		  , c.name category_name
		FROM
		  kakeibo.transactions t
		  INNER JOIN kakeibo.categories c 
		    ON  t.category_id = c.id
	`

	return query
}
