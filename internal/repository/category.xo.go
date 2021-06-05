// Code generated by xo. DO NOT EDIT.

package repository

import (
	"context"
	"database/sql"

	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

// Generated from 'kakeibo.categories'.
type categoryRepository struct{}

func (r *categoryRepository) SelectAll(db database.DB) ([]*entity.Category, error) {
	query := `
		SELECT id, type, name, created_at, modified_at, version
		FROM kakeibo.categories
		ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.Category{}
	for rows.Next() {
		e := entity.Category{}
		err = rows.Scan(&e.ID, &e.Type, &e.Name, &e.CreatedAt, &e.ModifiedAt, &e.Version)
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

func (r *categoryRepository) Insert(db database.DB, e *entity.Category) error {
	query := `
		INSERT INTO kakeibo.categories (
			type, name, created_at, modified_at
		) VALUES (
			?, ?, ?, ?
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	res, err := db.ExecContext(ctx, query, e.Type, e.Name, e.CreatedAt, e.ModifiedAt)
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

func (r *categoryRepository) Update(db database.DB, e *entity.Category) error {
	query := `
		UPDATE kakeibo.categories SET
			type = ?, name = ?, created_at = ?, modified_at = ?, version = version + 1
		WHERE id = ? AND version = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, query, e.Type, e.Name, e.CreatedAt, e.ModifiedAt, e.ID, e.Version)
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

	e.Version += 1
	return nil
}

func (r *categoryRepository) Delete(db database.DB, e *entity.Category) error {
	query := `
		DELETE FROM kakeibo.categories
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, query, e.ID)
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
	return nil
}

// Generated from index 'categories_id_pkey'.
func (r *categoryRepository) SelectByID(db database.DB, id int) (*entity.Category, error) {
	query := `
		SELECT id, type, name, created_at, modified_at, version
		FROM kakeibo.categories
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()
	e := entity.Category{}
	err := db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Type, &e.Name, &e.CreatedAt, &e.ModifiedAt, &e.Version)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Generated from index 'idx_categories_type_name'.
func (r *categoryRepository) SelectByTypeName(db database.DB, typ int, name string) (*entity.Category, error) {
	query := `
		SELECT id, type, name, created_at, modified_at, version
		FROM kakeibo.categories
		WHERE type = ? AND name = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), database.QueryTimeout)
	defer cancel()
	e := entity.Category{}
	err := db.QueryRowContext(ctx, query, typ, name).Scan(&e.ID, &e.Type, &e.Name, &e.CreatedAt, &e.ModifiedAt, &e.Version)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
