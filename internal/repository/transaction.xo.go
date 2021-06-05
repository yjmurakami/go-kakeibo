// Code generated by xo. DO NOT EDIT.

package repository

import (
	"database/sql"

	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

// Generated from 'kakeibo.transactions'.
type transactionRepository struct{}

func (r *transactionRepository) SelectAll(db database.DB) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, date, category_id, amount, note, created_at, modified_at
		FROM kakeibo.transactions
		ORDER BY id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.Transaction{}
	for rows.Next() {
		e := entity.Transaction{}
		err = rows.Scan(&e.ID, &e.UserID, &e.Date, &e.CategoryID, &e.Amount, &e.Note, &e.CreatedAt, &e.ModifiedAt)
		if err != nil {
			return nil, err
		}
		s = append(s, &e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *transactionRepository) Insert(db database.DB, e *entity.Transaction) error {
	query := `
		INSERT INTO kakeibo.transactions (
			user_id, date, category_id, amount, note, created_at, modified_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
	`

	res, err := db.Exec(query, e.UserID, e.Date, e.CategoryID, e.Amount, e.Note, e.CreatedAt, e.ModifiedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = int(id)
	return nil
}

func (r *transactionRepository) Update(db database.DB, e *entity.Transaction) error {
	query := `
		UPDATE kakeibo.transactions SET
			user_id = ?, date = ?, category_id = ?, amount = ?, note = ?, created_at = ?, modified_at = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, e.UserID, e.Date, e.CategoryID, e.Amount, e.Note, e.CreatedAt, e.ModifiedAt, e.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (r *transactionRepository) Delete(db database.DB, e *entity.Transaction) error {
	query := `
		DELETE FROM kakeibo.transactions
		WHERE id = ?
	`

	result, err := db.Exec(query, e.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

// Generated from index 'category_id'.
func (r *transactionRepository) SelectByCategoryID(db database.DB, categoryID int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, date, category_id, amount, note, created_at, modified_at
		FROM kakeibo.transactions
		WHERE category_id = ?
		ORDER BY id
	`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.Transaction{}
	for rows.Next() {
		e := entity.Transaction{}
		err = rows.Scan(&e.ID, &e.UserID, &e.Date, &e.CategoryID, &e.Amount, &e.Note, &e.CreatedAt, &e.ModifiedAt)
		if err != nil {
			return nil, err
		}
		s = append(s, &e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Generated from index 'transactions_id_pkey'.
func (r *transactionRepository) SelectByID(db database.DB, id int) (*entity.Transaction, error) {
	query := `
		SELECT id, user_id, date, category_id, amount, note, created_at, modified_at
		FROM kakeibo.transactions
		WHERE id = ?
	`
	e := entity.Transaction{}
	err := db.QueryRow(query, id).Scan(&e.ID, &e.UserID, &e.Date, &e.CategoryID, &e.Amount, &e.Note, &e.CreatedAt, &e.ModifiedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Generated from index 'user_id'.
func (r *transactionRepository) SelectByUserID(db database.DB, userID int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, date, category_id, amount, note, created_at, modified_at
		FROM kakeibo.transactions
		WHERE user_id = ?
		ORDER BY id
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.Transaction{}
	for rows.Next() {
		e := entity.Transaction{}
		err = rows.Scan(&e.ID, &e.UserID, &e.Date, &e.CategoryID, &e.Amount, &e.Note, &e.CreatedAt, &e.ModifiedAt)
		if err != nil {
			return nil, err
		}
		s = append(s, &e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return s, nil
}
