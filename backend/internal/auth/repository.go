package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrLoginUserNotFound = errors.New("login user not found")

type Repository struct {
	db *sql.DB
}

type LoginUserRecord struct {
	UserID       uint64
	DisplayName  string
	Email        string
	PasswordHash string
	CompanyID    uint64
	CompanyName  string
	Role         string
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindLoginUserByEmail(
	ctx context.Context,
	email string,
) (LoginUserRecord, error) {
	var record LoginUserRecord

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				u.id,
				u.display_name,
				u.email,
				u.password_hash,
				c.id,
				c.name,
				cm.role
			FROM users AS u
			INNER JOIN company_members AS cm
				ON cm.user_id = u.id
			INNER JOIN companies AS c
				ON c.id = cm.company_id
			WHERE u.email = ?
			  AND u.status = 'active'
			  AND cm.status = 'active'
			  AND c.status = 'active'
			LIMIT 1
		`,
		email,
	).Scan(
		&record.UserID,
		&record.DisplayName,
		&record.Email,
		&record.PasswordHash,
		&record.CompanyID,
		&record.CompanyName,
		&record.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginUserRecord{}, ErrLoginUserNotFound
		}

		return LoginUserRecord{}, fmt.Errorf(
			"failed to find login user: %w",
			err,
		)
	}

	return record, nil
}
