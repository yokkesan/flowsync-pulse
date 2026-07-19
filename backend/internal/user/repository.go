package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrCompanyNotFound = errors.New("company not found")

type Repository struct {
	db *sql.DB
}

type RegisterParams struct {
	CompanyID    uint64
	DisplayName  string
	Email        string
	PasswordHash string
}

type RegisterResult struct {
	UserID uint64
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Register(
	ctx context.Context,
	params RegisterParams,
) (RegisterResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return RegisterResult{}, fmt.Errorf(
			"failed to begin transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var companyExists bool

	err = tx.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM companies
				WHERE id = ?
				  AND status = 'active'
			)
		`,
		params.CompanyID,
	).Scan(&companyExists)
	if err != nil {
		return RegisterResult{}, fmt.Errorf(
			"failed to check company: %w",
			err,
		)
	}

	if !companyExists {
		return RegisterResult{}, ErrCompanyNotFound
	}

	userResult, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO users (
				display_name,
				email,
				password_hash,
				avatar_type,
				avatar_key,
				status
			)
			VALUES (?, ?, ?, 'preset', NULL, 'active')
		`,
		params.DisplayName,
		params.Email,
		params.PasswordHash,
	)
	if err != nil {
		return RegisterResult{}, fmt.Errorf(
			"failed to create user: %w",
			err,
		)
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		return RegisterResult{}, fmt.Errorf(
			"failed to get user id: %w",
			err,
		)
	}

	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO company_members (
				company_id,
				user_id,
				role,
				status,
				joined_at
			)
			VALUES (?, ?, 'owner', 'active', CURRENT_TIMESTAMP)
		`,
		params.CompanyID,
		userID,
	)
	if err != nil {
		return RegisterResult{}, fmt.Errorf(
			"failed to create company membership: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return RegisterResult{}, fmt.Errorf(
			"failed to commit transaction: %w",
			err,
		)
	}

	return RegisterResult{
		UserID: uint64(userID),
	}, nil
}
