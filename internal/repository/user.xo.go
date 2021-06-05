// Code generated by xo. DO NOT EDIT.

package repository

import (
	"database/sql"

	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/entity"
)

// Generated from 'kakeibo.users'.
type userRepository struct{}

func (r *userRepository) SelectAll(db database.DB) ([]*entity.User, error) {
	query := `
		SELECT id, login_id, login_password, created_at, modified_at, version
		FROM kakeibo.users
		ORDER BY id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s := []*entity.User{}
	for rows.Next() {
		e := entity.User{}
		err = rows.Scan(&e.ID, &e.LoginID, &e.LoginPassword, &e.CreatedAt, &e.ModifiedAt, &e.Version)
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

func (r *userRepository) Insert(db database.DB, e *entity.User) error {
	query := `
		INSERT INTO kakeibo.users (
			login_id, login_password, created_at, modified_at
		) VALUES (
			?, ?, ?, ?
		)
	`

	res, err := db.Exec(query, e.LoginID, e.LoginPassword, e.CreatedAt, e.ModifiedAt)
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

func (r *userRepository) Update(db database.DB, e *entity.User) error {
	query := `
		UPDATE kakeibo.users SET
			login_id = ?, login_password = ?, created_at = ?, modified_at = ?, version = version + 1
		WHERE id = ? AND version = ?
	`

	result, err := db.Exec(query, e.LoginID, e.LoginPassword, e.CreatedAt, e.ModifiedAt, e.ID, e.Version)
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

func (r *userRepository) Delete(db database.DB, e *entity.User) error {
	query := `
		DELETE FROM kakeibo.users
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
	return nil
}

// Generated from index 'login_id_UNIQUE'.
func (r *userRepository) SelectByLoginID(db database.DB, loginID string) (*entity.User, error) {
	query := `
		SELECT id, login_id, login_password, created_at, modified_at, version
		FROM kakeibo.users
		WHERE login_id = ?
	`
	e := entity.User{}
	err := db.QueryRow(query, loginID).Scan(&e.ID, &e.LoginID, &e.LoginPassword, &e.CreatedAt, &e.ModifiedAt, &e.Version)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Generated from index 'users_id_pkey'.
func (r *userRepository) SelectByID(db database.DB, id int) (*entity.User, error) {
	query := `
		SELECT id, login_id, login_password, created_at, modified_at, version
		FROM kakeibo.users
		WHERE id = ?
	`
	e := entity.User{}
	err := db.QueryRow(query, id).Scan(&e.ID, &e.LoginID, &e.LoginPassword, &e.CreatedAt, &e.ModifiedAt, &e.Version)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
