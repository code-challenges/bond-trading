package database

import (
	"context"
	"errors"
	"fmt"

	. "github.com/asalvi0/bond-trading/internal/models"
)

func (db *Database) CreateUser(ctx context.Context, user *User) error {
	sql := `INSERT INTO users (email, password_hash, active, created_at)
			VALUES ($1, $2, $3, $4) RETURNING "id"`

	row := db.dbPool.QueryRow(ctx, sql, user.Email, user.Password, user.Active, user.CreatedAt)

	var id uint
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("Unable to retrieve record ID: %w", err)
	}
	user.ID = id

	return nil
}

func (db *Database) UpdateUser(ctx context.Context, user *User) error {
	sql := `UPDATE users SET
				email = $2,
				password_hash = $3,
				active = $4
				updated_at = $5
			WHERE id = $1`

	tag, err := db.dbPool.Exec(ctx, sql, user.ID, user.Email, user.Password, user.Active, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Unable to update record: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("No records were updated")
	}

	return nil
}

func (db *Database) GetUserByID(ctx context.Context, id uint) (*User, error) {
	sql := `SELECT * FROM users WHERE id = $1`

	row := db.dbPool.QueryRow(ctx, sql, id)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to select record: %w", err)
	}

	return &user, nil
}

func (db *Database) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	sql := `SELECT * FROM users WHERE email = $1`

	row := db.dbPool.QueryRow(ctx, sql, email)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to select record: %w", err)
	}

	return &user, nil
}
