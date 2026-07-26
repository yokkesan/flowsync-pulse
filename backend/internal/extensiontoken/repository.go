package extensiontoken

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrExtensionTokenNotFound = errors.New(
	"extension token not found",
)

type Repository struct {
	db *sql.DB
}

type CreateParams struct {
	UserID    uint64
	CompanyID uint64
	Name      string
	TokenHash string
	ExpiresAt time.Time
}

type CreateResult struct {
	TokenID uint64
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
			INSERT INTO extension_tokens (
				user_id,
				company_id,
				name,
				token_hash,
				expires_at
			)
			VALUES (?, ?, ?, ?, ?)
		`,
		params.UserID,
		params.CompanyID,
		params.Name,
		params.TokenHash,
		params.ExpiresAt,
	)
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to create extension token: %w",
			err,
		)
	}

	tokenID, err := result.LastInsertId()
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to get extension token id: %w",
			err,
		)
	}

	return CreateResult{
		TokenID: uint64(tokenID),
	}, nil
}

func (r *Repository) FindAllByUser(
	ctx context.Context,
	userID uint64,
	companyID uint64,
) ([]Response, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				name,
				last_used_at,
				expires_at,
				revoked_at,
				created_at
			FROM extension_tokens
			WHERE user_id = ?
			  AND company_id = ?
			ORDER BY created_at DESC, id DESC
		`,
		userID,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to find extension tokens: %w",
			err,
		)
	}
	defer rows.Close()

	tokens := make([]Response, 0)

	for rows.Next() {
		var response Response
		var lastUsedAt sql.NullTime
		var expiresAt sql.NullTime
		var revokedAt sql.NullTime

		if err := rows.Scan(
			&response.TokenID,
			&response.Name,
			&lastUsedAt,
			&expiresAt,
			&revokedAt,
			&response.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan extension token: %w",
				err,
			)
		}

		if !expiresAt.Valid {
			return nil, errors.New(
				"extension token expiration is not set",
			)
		}

		response.ExpiresAt = expiresAt.Time

		if lastUsedAt.Valid {
			value := lastUsedAt.Time
			response.LastUsedAt = &value
		}

		if revokedAt.Valid {
			value := revokedAt.Time
			response.RevokedAt = &value
		}

		tokens = append(tokens, response)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while iterating extension tokens: %w",
			err,
		)
	}

	return tokens, nil
}

func (r *Repository) Revoke(
	ctx context.Context,
	tokenID uint64,
	userID uint64,
	companyID uint64,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE extension_tokens
			SET revoked_at = COALESCE(
				revoked_at,
				CURRENT_TIMESTAMP
			)
			WHERE id = ?
			  AND user_id = ?
			  AND company_id = ?
		`,
		tokenID,
		userID,
		companyID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to revoke extension token: %w",
			err,
		)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to get revoked extension token rows: %w",
			err,
		)
	}

	if affectedRows == 0 {
		return ErrExtensionTokenNotFound
	}

	return nil
}

func (r *Repository) FindValidByHash(
	ctx context.Context,
	tokenHash string,
	now time.Time,
) (
	uint64,
	uint64,
	uint64,
	error,
) {
	var tokenID uint64
	var userID uint64
	var companyID uint64

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				user_id,
				company_id
			FROM extension_tokens
			WHERE token_hash = ?
			  AND revoked_at IS NULL
			  AND expires_at IS NOT NULL
			  AND expires_at > ?
			LIMIT 1
		`,
		tokenHash,
		now,
	).Scan(
		&tokenID,
		&userID,
		&companyID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, 0,
				ErrExtensionTokenNotFound
		}

		return 0, 0, 0, fmt.Errorf(
			"failed to authenticate extension token: %w",
			err,
		)
	}

	return tokenID, userID, companyID, nil
}

func (r *Repository) UpdateLastUsedAt(
	ctx context.Context,
	tokenID uint64,
	usedAt time.Time,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE extension_tokens
			SET last_used_at = ?
			WHERE id = ?
			  AND revoked_at IS NULL
			  AND expires_at IS NOT NULL
			  AND expires_at > ?
		`,
		usedAt,
		tokenID,
		usedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update extension token last used time: %w",
			err,
		)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to get updated extension token rows: %w",
			err,
		)
	}

	if affectedRows == 0 {
		return ErrExtensionTokenNotFound
	}

	return nil
}
