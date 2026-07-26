package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrProjectNotFound = errors.New("project not found")

type Repository struct {
	db *sql.DB
}

type CreateParams struct {
	CompanyID   uint64
	CreatedByID uint64
	Name        string
	Slug        string
	ProjectKey  string
	Description *string
	Status      string
	StartDate   *string
	EndDate     *string
	MemberIDs   []uint64
}

type CreateResult struct {
	ProjectID uint64
}

type UpdateParams struct {
	ProjectID   uint64
	CompanyID   uint64
	UpdatedByID uint64
	Name        string
	Slug        string
	ProjectKey  *string
	Description *string
	Status      string
	StartDate   *string
	EndDate     *string
	MemberIDs   []uint64
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) HasProjectAccess(
	ctx context.Context,
	projectID uint64,
	userID uint64,
	companyID uint64,
) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM projects p
				INNER JOIN project_members pm
					ON pm.project_id = p.id
				WHERE p.id = ?
				  AND p.company_id = ?
				  AND pm.user_id = ?
				  AND pm.status = 'active'
			)
		`,
		projectID,
		companyID,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check project access: %w",
			err,
		)
	}

	return exists, nil
}

func (r *Repository) AreActiveCompanyMembers(
	ctx context.Context,
	companyID uint64,
	memberIDs []uint64,
) (bool, error) {
	if len(memberIDs) == 0 {
		return false, nil
	}

	query := `
		SELECT COUNT(DISTINCT user_id)
		FROM company_members
		WHERE company_id = ?
		  AND status = 'active'
		  AND user_id IN (
	`

	args := make([]any, 0, len(memberIDs)+1)
	args = append(args, companyID)

	for index, memberID := range memberIDs {
		if index > 0 {
			query += ", "
		}

		query += "?"
		args = append(args, memberID)
	}

	query += ")"

	var activeMemberCount int

	if err := r.db.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(&activeMemberCount); err != nil {
		return false, fmt.Errorf(
			"failed to check company members: %w",
			err,
		)
	}

	return activeMemberCount == len(memberIDs), nil
}

func (r *Repository) SlugExists(
	ctx context.Context,
	companyID uint64,
	slug string,
	excludeProjectID uint64,
) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM projects
				WHERE company_id = ?
				  AND slug = ?
				  AND (? = 0 OR id <> ?)
			)
		`,
		companyID,
		slug,
		excludeProjectID,
		excludeProjectID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check project slug: %w",
			err,
		)
	}

	return exists, nil
}

func (r *Repository) ProjectKeyExists(
	ctx context.Context,
	companyID uint64,
	projectKey string,
	excludeProjectID uint64,
) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM projects
				WHERE company_id = ?
				  AND project_key = ?
				  AND (? = 0 OR id <> ?)
			)
		`,
		companyID,
		projectKey,
		excludeProjectID,
		excludeProjectID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check project key: %w",
			err,
		)
	}

	return exists, nil
}

func (r *Repository) Create(
	ctx context.Context,
	params CreateParams,
) (*CreateResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to begin project creation transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	startedAt, err := nullableDate(params.StartDate)
	if err != nil {
		return nil, err
	}

	endedAt, err := nullableDate(params.EndDate)
	if err != nil {
		return nil, err
	}

	result, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO projects (
				company_id,
				name,
				slug,
				project_key,
				description,
				status,
				started_at,
				ended_at
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		params.CompanyID,
		params.Name,
		params.Slug,
		params.ProjectKey,
		params.Description,
		params.Status,
		startedAt,
		endedAt,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create project: %w",
			err,
		)
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get created project id: %w",
			err,
		)
	}

	if err := insertProjectMembers(
		ctx,
		tx,
		uint64(projectID),
		params.CreatedByID,
		params.MemberIDs,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(
			"failed to commit project creation transaction: %w",
			err,
		)
	}

	return &CreateResult{
		ProjectID: uint64(projectID),
	}, nil
}

func (r *Repository) FindAllByUser(
	ctx context.Context,
	companyID uint64,
	userID uint64,
) ([]Response, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				p.id,
				p.company_id,
				p.name,
				p.slug,
				p.project_key,
				p.description,
				p.status,
				p.started_at,
				p.ended_at,
				COUNT(DISTINCT t.id) AS task_count,
				p.created_at,
				p.updated_at
			FROM projects p
			INNER JOIN project_members current_member
				ON current_member.project_id = p.id
				AND current_member.user_id = ?
				AND current_member.status = 'active'
			LEFT JOIN tasks t
				ON t.project_id = p.id
			WHERE p.company_id = ?
			GROUP BY
				p.id,
				p.company_id,
				p.name,
				p.slug,
				p.project_key,
				p.description,
				p.status,
				p.started_at,
				p.ended_at,
				p.created_at,
				p.updated_at
			ORDER BY p.created_at DESC, p.id DESC
		`,
		userID,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to find projects: %w",
			err,
		)
	}
	defer rows.Close()

	projects := make([]Response, 0)

	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, err
		}

		members, err := r.findMembers(ctx, project.ProjectID)
		if err != nil {
			return nil, err
		}

		project.Members = members
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while iterating projects: %w",
			err,
		)
	}

	return projects, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	projectID uint64,
	companyID uint64,
) (*Response, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				p.id,
				p.company_id,
				p.name,
				p.slug,
				p.project_key,
				p.description,
				p.status,
				p.started_at,
				p.ended_at,
				COUNT(DISTINCT t.id) AS task_count,
				p.created_at,
				p.updated_at
			FROM projects p
			LEFT JOIN tasks t
				ON t.project_id = p.id
			WHERE p.id = ?
			  AND p.company_id = ?
			GROUP BY
				p.id,
				p.company_id,
				p.name,
				p.slug,
				p.project_key,
				p.description,
				p.status,
				p.started_at,
				p.ended_at,
				p.created_at,
				p.updated_at
		`,
		projectID,
		companyID,
	)

	project, err := scanProject(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	members, err := r.findMembers(ctx, project.ProjectID)
	if err != nil {
		return nil, err
	}

	project.Members = members

	return &project, nil
}

func (r *Repository) Update(
	ctx context.Context,
	params UpdateParams,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"failed to begin project update transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	startedAt, err := nullableDate(params.StartDate)
	if err != nil {
		return err
	}

	endedAt, err := nullableDate(params.EndDate)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(
		ctx,
		`
			UPDATE projects
			SET
				name = ?,
				slug = ?,
				project_key = COALESCE(project_key, ?),
				description = ?,
				status = ?,
				started_at = ?,
				ended_at = ?
			WHERE id = ?
			  AND company_id = ?
		`,
		params.Name,
		params.Slug,
		params.ProjectKey,
		params.Description,
		params.Status,
		startedAt,
		endedAt,
		params.ProjectID,
		params.CompanyID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update project: %w",
			err,
		)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to get updated project rows: %w",
			err,
		)
	}

	if affectedRows == 0 {
		return ErrProjectNotFound
	}

	if params.ProjectKey != nil {
		if _, err := tx.ExecContext(
			ctx,
			`
				UPDATE tasks
				SET task_key = CONCAT(?, '-', task_number)
				WHERE project_id = ?
				  AND task_key IS NULL
			`,
			*params.ProjectKey,
			params.ProjectID,
		); err != nil {
			return fmt.Errorf(
				"failed to generate existing task keys: %w",
				err,
			)
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			DELETE FROM project_members
			WHERE project_id = ?
		`,
		params.ProjectID,
	); err != nil {
		return fmt.Errorf(
			"failed to replace project members: %w",
			err,
		)
	}

	if err := insertProjectMembers(
		ctx,
		tx,
		params.ProjectID,
		params.UpdatedByID,
		params.MemberIDs,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"failed to commit project update transaction: %w",
			err,
		)
	}

	return nil
}

func (r *Repository) Delete(
	ctx context.Context,
	projectID uint64,
	companyID uint64,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM projects
			WHERE id = ?
			  AND company_id = ?
		`,
		projectID,
		companyID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to delete project: %w",
			err,
		)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to get deleted project rows: %w",
			err,
		)
	}

	if affectedRows == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (r *Repository) findMembers(
	ctx context.Context,
	projectID uint64,
) ([]MemberResponse, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				pm.user_id,
				u.display_name,
				pm.role,
				pm.status
			FROM project_members pm
			INNER JOIN users u
				ON u.id = pm.user_id
			WHERE pm.project_id = ?
			ORDER BY
				CASE pm.role
					WHEN 'manager' THEN 1
					WHEN 'developer' THEN 2
					ELSE 3
				END,
				pm.id
		`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to find project members: %w",
			err,
		)
	}
	defer rows.Close()

	members := make([]MemberResponse, 0)

	for rows.Next() {
		var member MemberResponse

		if err := rows.Scan(
			&member.UserID,
			&member.DisplayName,
			&member.Role,
			&member.Status,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan project member: %w",
				err,
			)
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while iterating project members: %w",
			err,
		)
	}

	return members, nil
}

func insertProjectMembers(
	ctx context.Context,
	tx *sql.Tx,
	projectID uint64,
	managerUserID uint64,
	memberIDs []uint64,
) error {
	for _, memberID := range memberIDs {
		role := "developer"

		if memberID == managerUserID {
			role = "manager"
		}

		if _, err := tx.ExecContext(
			ctx,
			`
				INSERT INTO project_members (
					project_id,
					user_id,
					role,
					status,
					joined_at
				)
				VALUES (?, ?, ?, 'active', CURRENT_TIMESTAMP)
			`,
			projectID,
			memberID,
			role,
		); err != nil {
			return fmt.Errorf(
				"failed to insert project member: %w",
				err,
			)
		}
	}

	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProject(scanner rowScanner) (Response, error) {
	var project Response
	var projectKey sql.NullString
	var description sql.NullString
	var startedAt sql.NullTime
	var endedAt sql.NullTime

	if err := scanner.Scan(
		&project.ProjectID,
		&project.CompanyID,
		&project.Name,
		&project.Slug,
		&projectKey,
		&description,
		&project.Status,
		&startedAt,
		&endedAt,
		&project.TaskCount,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return Response{}, err
	}

	project.ProjectKey = nullableStringPointer(projectKey)
	project.Description = nullableStringPointer(description)
	project.StartDate = nullableDateString(startedAt)
	project.EndDate = nullableDateString(endedAt)
	project.Members = make([]MemberResponse, 0)

	return project, nil
}

func nullableDate(value *string) (any, error) {
	if value == nil {
		return nil, nil
	}

	normalizedValue := strings.TrimSpace(*value)
	if normalizedValue == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse("2006-01-02", normalizedValue)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse project date: %w",
			err,
		)
	}

	return parsedDate, nil
}

func nullableDateString(value sql.NullTime) *string {
	if !value.Valid {
		return nil
	}

	formattedValue := value.Time.Format("2006-01-02")

	return &formattedValue
}

func nullableStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}
