package user

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

type CreateCompanyOwnerParams struct {
	CompanyName  string
	CompanySlug  string
	DisplayName  string
	Email        string
	PasswordHash string
	AvatarKey    string
}

type CreateCompanyOwnerResult struct {
	UserID    uint64
	CompanyID uint64
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateCompanyOwner(
	ctx context.Context,
	params CreateCompanyOwnerParams,
) (CreateCompanyOwnerResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
			"failed to begin transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	companyResult, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO companies (
				name,
				slug,
				status
			)
			VALUES (?, ?, 'active')
		`,
		params.CompanyName,
		params.CompanySlug,
	)
	if err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
			"failed to create company: %w",
			err,
		)
	}

	companyID, err := companyResult.LastInsertId()
	if err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
			"failed to get company id: %w",
			err,
		)
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
			VALUES (?, ?, ?, 'preset', NULLIF(?, ''), 'active')
		`,
		params.DisplayName,
		params.Email,
		params.PasswordHash,
		params.AvatarKey,
	)
	if err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
			"failed to create user: %w",
			err,
		)
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
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
		companyID,
		userID,
	)
	if err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
			"failed to create company membership: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return CreateCompanyOwnerResult{}, fmt.Errorf(
			"failed to commit transaction: %w",
			err,
		)
	}

	return CreateCompanyOwnerResult{
		UserID:    uint64(userID),
		CompanyID: uint64(companyID),
	}, nil
}
