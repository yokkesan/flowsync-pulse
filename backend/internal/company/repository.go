package company

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

type CreateParams struct {
	Name string
	Slug string
}

type CreateResult struct {
	CompanyID uint64
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	params CreateParams,
) (CreateResult, error) {
	result, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO companies (
				name,
				slug,
				status
			)
			VALUES (?, ?, 'active')
		`,
		params.Name,
		params.Slug,
	)
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to create company: %w",
			err,
		)
	}

	companyID, err := result.LastInsertId()
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to get company id: %w",
			err,
		)
	}

	return CreateResult{
		CompanyID: uint64(companyID),
	}, nil
}
