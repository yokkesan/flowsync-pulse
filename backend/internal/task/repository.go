package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrProjectKeyNotSet = errors.New(
	"project key is not set",
)

type Repository struct {
	db *sql.DB
}

type CreateParams struct {
	ProjectID      uint64
	Name           string
	Description    *string
	AssigneeUserID uint64
	BranchName     string
	Status         string
	Priority       string
	StartDate      *string
	DueDate        *string
	CompletedAt    *time.Time
}

type CreateResult struct {
	TaskID uint64
}

type UpdateParams struct {
	ProjectID      uint64
	TaskID         uint64
	Name           string
	Description    *string
	AssigneeUserID uint64
	BranchName     string
	Status         string
	Priority       string
	StartDate      *string
	DueDate        *string
	CompletedAt    *time.Time
}

type FindByIDParams struct {
	ProjectID uint64
	TaskID    uint64
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

func (r *Repository) IsActiveProjectMember(
	ctx context.Context,
	projectID uint64,
	userID uint64,
) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM project_members
				WHERE project_id = ?
				  AND user_id = ?
				  AND status = 'active'
			)
		`,
		projectID,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check project member: %w",
			err,
		)
	}

	return exists, nil
}

func (r *Repository) BranchExistsInProject(
	ctx context.Context,
	projectID uint64,
	branchName string,
	excludeTaskID uint64,
) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM task_branches tb
				INNER JOIN tasks t
					ON t.id = tb.task_id
				WHERE t.project_id = ?
				  AND tb.branch_name = ?
				  AND (? = 0 OR t.id <> ?)
			)
		`,
		projectID,
		branchName,
		excludeTaskID,
		excludeTaskID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check branch duplication: %w",
			err,
		)
	}

	return exists, nil
}

func (r *Repository) Create(
	ctx context.Context,
	params CreateParams,
) (CreateResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to begin task transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var projectKey sql.NullString
	var taskNumber uint64

	err = tx.QueryRowContext(
		ctx,
		`
			SELECT
				project_key,
				next_task_number
			FROM projects
			WHERE id = ?
			FOR UPDATE
		`,
		params.ProjectID,
	).Scan(
		&projectKey,
		&taskNumber,
	)
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to lock project task sequence: %w",
			err,
		)
	}

	if !projectKey.Valid || projectKey.String == "" {
		return CreateResult{}, ErrProjectKeyNotSet
	}

	taskKey := fmt.Sprintf(
		"%s-%d",
		projectKey.String,
		taskNumber,
	)

	result, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO tasks (
				project_id,
				task_number,
				task_key,
				name,
				description,
				assignee_user_id,
				status,
				priority,
				start_date,
				due_date,
				completed_at
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		params.ProjectID,
		taskNumber,
		taskKey,
		params.Name,
		params.Description,
		params.AssigneeUserID,
		params.Status,
		params.Priority,
		params.StartDate,
		params.DueDate,
		params.CompletedAt,
	)
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to create task: %w",
			err,
		)
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to get task id: %w",
			err,
		)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			UPDATE projects
			SET next_task_number = next_task_number + 1
			WHERE id = ?
		`,
		params.ProjectID,
	); err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to update project task sequence: %w",
			err,
		)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO task_branches (
				task_id,
				branch_name
			)
			VALUES (?, ?)
		`,
		taskID,
		params.BranchName,
	); err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to create task branch: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return CreateResult{}, fmt.Errorf(
			"failed to commit task transaction: %w",
			err,
		)
	}

	return CreateResult{
		TaskID: uint64(taskID),
	}, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	params FindByIDParams,
) (Response, error) {
	var response Response
	var taskKey sql.NullString
	var description sql.NullString
	var startDate sql.NullTime
	var dueDate sql.NullTime
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				t.id,
				t.project_id,
				p.name,
				t.task_number,
				t.task_key,
				t.name,
				t.description,
				t.assignee_user_id,
				u.display_name,
				tb.branch_name,
				t.status,
				t.priority,
				t.start_date,
				t.due_date,
				t.completed_at,
				t.created_at,
				t.updated_at
			FROM tasks t
			INNER JOIN projects p
				ON p.id = t.project_id
			INNER JOIN users u
				ON u.id = t.assignee_user_id
			INNER JOIN task_branches tb
				ON tb.task_id = t.id
			WHERE t.id = ?
			  AND t.project_id = ?
		`,
		params.TaskID,
		params.ProjectID,
	).Scan(
		&response.TaskID,
		&response.ProjectID,
		&response.ProjectName,
		&response.TaskNumber,
		&taskKey,
		&response.Name,
		&description,
		&response.AssigneeUserID,
		&response.AssigneeName,
		&response.BranchName,
		&response.Status,
		&response.Priority,
		&startDate,
		&dueDate,
		&completedAt,
		&response.CreatedAt,
		&response.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, ErrTaskNotFound
		}

		return Response{}, fmt.Errorf(
			"failed to find task by id: %w",
			err,
		)
	}

	setNullableTaskFields(
		&response,
		taskKey,
		description,
		startDate,
		dueDate,
		completedAt,
	)

	return response, nil
}

func (r *Repository) FindAllByProjectID(
	ctx context.Context,
	projectID uint64,
) ([]Response, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				t.id,
				t.project_id,
				p.name,
				t.task_number,
				t.task_key,
				t.name,
				t.description,
				t.assignee_user_id,
				u.display_name,
				tb.branch_name,
				t.status,
				t.priority,
				t.start_date,
				t.due_date,
				t.completed_at,
				t.created_at,
				t.updated_at
			FROM tasks t
			INNER JOIN projects p
				ON p.id = t.project_id
			INNER JOIN users u
				ON u.id = t.assignee_user_id
			INNER JOIN task_branches tb
				ON tb.task_id = t.id
			WHERE t.project_id = ?
			ORDER BY
				t.created_at DESC,
				t.id DESC
		`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to find tasks by project id: %w",
			err,
		)
	}
	defer rows.Close()

	tasks := make([]Response, 0)

	for rows.Next() {
		response, err := scanTaskResponse(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, response)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate tasks: %w",
			err,
		)
	}

	return tasks, nil
}

func (r *Repository) FindAllAccessible(
	ctx context.Context,
	companyID uint64,
) ([]Response, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				t.id,
				t.project_id,
				p.name,
				t.task_number,
				t.task_key,
				t.name,
				t.description,
				t.assignee_user_id,
				u.display_name,
				tb.branch_name,
				t.status,
				t.priority,
				t.start_date,
				t.due_date,
				t.completed_at,
				t.created_at,
				t.updated_at
			FROM tasks t
			INNER JOIN projects p
				ON p.id = t.project_id
			INNER JOIN users u
				ON u.id = t.assignee_user_id
			INNER JOIN task_branches tb
				ON tb.task_id = t.id
			WHERE p.company_id = ?
			ORDER BY
				t.created_at DESC,
				t.id DESC
		`,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to find accessible tasks: %w",
			err,
		)
	}
	defer rows.Close()

	tasks := make([]Response, 0)

	for rows.Next() {
		response, err := scanTaskResponse(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, response)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate accessible tasks: %w",
			err,
		)
	}

	return tasks, nil
}

func (r *Repository) Update(
	ctx context.Context,
	params UpdateParams,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"failed to begin task update transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var exists bool

	err = tx.QueryRowContext(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1
				FROM tasks
				WHERE id = ?
				  AND project_id = ?
			)
		`,
		params.TaskID,
		params.ProjectID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf(
			"failed to check task existence: %w",
			err,
		)
	}

	if !exists {
		return ErrTaskNotFound
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE tasks
			SET
				name = ?,
				description = ?,
				assignee_user_id = ?,
				status = ?,
				priority = ?,
				start_date = ?,
				due_date = ?,
				completed_at = ?
			WHERE id = ?
			  AND project_id = ?
		`,
		params.Name,
		params.Description,
		params.AssigneeUserID,
		params.Status,
		params.Priority,
		params.StartDate,
		params.DueDate,
		params.CompletedAt,
		params.TaskID,
		params.ProjectID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update task: %w",
			err,
		)
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE task_branches
			SET branch_name = ?
			WHERE task_id = ?
		`,
		params.BranchName,
		params.TaskID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update task branch: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"failed to commit task update transaction: %w",
			err,
		)
	}

	return nil
}

func (r *Repository) Delete(
	ctx context.Context,
	projectID uint64,
	taskID uint64,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM tasks
			WHERE id = ?
			  AND project_id = ?
		`,
		taskID,
		projectID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to delete task: %w",
			err,
		)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to get deleted task rows: %w",
			err,
		)
	}

	if affectedRows == 0 {
		return ErrTaskNotFound
	}

	return nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTaskResponse(
	scanner taskScanner,
) (Response, error) {
	var response Response
	var taskKey sql.NullString
	var description sql.NullString
	var startDate sql.NullTime
	var dueDate sql.NullTime
	var completedAt sql.NullTime

	if err := scanner.Scan(
		&response.TaskID,
		&response.ProjectID,
		&response.ProjectName,
		&response.TaskNumber,
		&taskKey,
		&response.Name,
		&description,
		&response.AssigneeUserID,
		&response.AssigneeName,
		&response.BranchName,
		&response.Status,
		&response.Priority,
		&startDate,
		&dueDate,
		&completedAt,
		&response.CreatedAt,
		&response.UpdatedAt,
	); err != nil {
		return Response{}, fmt.Errorf(
			"failed to scan task: %w",
			err,
		)
	}

	setNullableTaskFields(
		&response,
		taskKey,
		description,
		startDate,
		dueDate,
		completedAt,
	)

	return response, nil
}

func setNullableTaskFields(
	response *Response,
	taskKey sql.NullString,
	description sql.NullString,
	startDate sql.NullTime,
	dueDate sql.NullTime,
	completedAt sql.NullTime,
) {
	if taskKey.Valid {
		response.TaskKey = taskKey.String
	}

	if description.Valid {
		response.Description = &description.String
	}

	if startDate.Valid {
		value := startDate.Time.Format(
			time.DateOnly,
		)
		response.StartDate = &value
	}

	if dueDate.Valid {
		value := dueDate.Time.Format(
			time.DateOnly,
		)
		response.DueDate = &value
	}

	if completedAt.Valid {
		response.CompletedAt = &completedAt.Time
	}
}
