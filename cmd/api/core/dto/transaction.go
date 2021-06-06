package dto

import (
	"context"
	"fmt"
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/core"
	"github.com/yjmurakami/go-kakeibo/internal/database"
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

// Generated from SelectTransactionById.sql
func SelectTransactionById(db database.DB, id int) (*Transaction, error) {
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
		WHERE
		  t.id = ?
		
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	d := Transaction{}
	err := db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.UserID, &d.Date, &d.Amount, &d.Note, &d.CreatedAt, &d.ModifiedAt, &d.Version, &d.CategoryID, &d.CategoryType, &d.CategoryName)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Generated from SelectTransactions.sql
func SelectTransactions(db database.DB, from time.Time, to time.Time, filter core.Filter) ([]*Transaction, core.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
		  COUNT(*) OVER() total_records
		  , t.id
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
		WHERE
		  (t.date >= ? OR ? IS NULL)
		  AND (t.date <= ? OR ? IS NULL)
		ORDER BY
		  %s %s, t.id
		LIMIT ? OFFSET ?
	`, filter.SortColumn(), filter.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, from, from, to, to, filter.Limit(), filter.Offset())
	if err != nil {
		return nil, core.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	s := []*Transaction{}
	for rows.Next() {
		d := Transaction{}
		err = rows.Scan(&totalRecords, &d.ID, &d.UserID, &d.Date, &d.Amount, &d.Note, &d.CreatedAt, &d.ModifiedAt, &d.Version, &d.CategoryID, &d.CategoryType, &d.CategoryName)
		if err != nil {
			return nil, core.Metadata{}, err
		}
		s = append(s, &d)
	}

	err = rows.Err()
	if err != nil {
		return nil, core.Metadata{}, err
	}

	metadata := core.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)
	return s, metadata, nil
}
