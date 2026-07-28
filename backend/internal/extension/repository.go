package extension

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrRepositoryNotFound = errors.New(
		"repository not found",
	)
	ErrActiveSessionNotFound = errors.New(
		"active work session not found",
	)
)

type Repository struct {
	db *sql.DB
}

type RepositoryCandidate struct {
	RepositoryID uint64
	ProjectID    uint64
	ProjectName  string
	RemoteURL    string
}

type TaskMatch struct {
	TaskID   uint64
	TaskKey  string
	TaskName string
}

type StartSessionParams struct {
	UserID        uint64
	ProjectID     uint64
	RepositoryID  uint64
	TaskID        *uint64
	BranchName    string
	TicketKey     *string
	WorkspaceName *string
	MatchStatus   string
	OccurredAt    time.Time
}

type SessionUpdateResult struct {
	SessionID       uint64
	SessionCreated  bool
	ProjectID       uint64
	RepositoryID    uint64
	TaskID          *uint64
	BranchName      string
	TicketKey       *string
	WorkspaceName   *string
	MatchStatus     string
	Status          string
	StartedAt       time.Time
	LastHeartbeatAt time.Time
	EndedAt         *time.Time
}

type HeartbeatResult struct {
	SessionID       uint64
	Status          string
	LastHeartbeatAt time.Time
}

type DisconnectResult struct {
	SessionID uint64
	Status    string
	EndReason string
	EndedAt   time.Time
}

func NewRepository(
	db *sql.DB,
) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindRepositoryCandidates(
	ctx context.Context,
	companyID uint64,
) ([]RepositoryCandidate, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				r.id,
				r.project_id,
				p.name,
				r.remote_url
			FROM repositories r
			INNER JOIN projects p
				ON p.id = r.project_id
			WHERE p.company_id = ?
			  AND r.status = 'active'
			ORDER BY r.id
		`,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to find repository candidates: %w",
			err,
		)
	}
	defer rows.Close()

	repositories := make(
		[]RepositoryCandidate,
		0,
	)

	for rows.Next() {
		var repository RepositoryCandidate

		if err := rows.Scan(
			&repository.RepositoryID,
			&repository.ProjectID,
			&repository.ProjectName,
			&repository.RemoteURL,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan repository candidate: %w",
				err,
			)
		}

		repositories = append(
			repositories,
			repository,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while iterating repository candidates: %w",
			err,
		)
	}

	return repositories, nil
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
			"failed to check extension project access: %w",
			err,
		)
	}

	return exists, nil
}

func (r *Repository) FindTaskByKey(
	ctx context.Context,
	projectID uint64,
	taskKey string,
) (*TaskMatch, error) {
	var task TaskMatch

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				t.id,
				t.task_key,
				t.name
			FROM tasks t
			WHERE t.project_id = ?
			  AND t.task_key = ?
			LIMIT 1
		`,
		projectID,
		taskKey,
	).Scan(
		&task.TaskID,
		&task.TaskKey,
		&task.TaskName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf(
			"failed to find task by key: %w",
			err,
		)
	}

	return &task, nil
}

func (r *Repository) SaveWorkContext(
	ctx context.Context,
	params StartSessionParams,
) (SessionUpdateResult, error) {
	tx, err := r.db.BeginTx(
		ctx,
		nil,
	)
	if err != nil {
		return SessionUpdateResult{}, fmt.Errorf(
			"failed to begin work session transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := lockUser(
		ctx,
		tx,
		params.UserID,
	); err != nil {
		return SessionUpdateResult{}, err
	}

	currentSession, err := findActiveSessionForUpdate(
		ctx,
		tx,
		params.UserID,
	)
	if err != nil {
		return SessionUpdateResult{}, err
	}

	if currentSession != nil &&
		currentSession.RepositoryID == params.RepositoryID &&
		currentSession.BranchName == params.BranchName {
		if err := updateExistingSessionContext(
			ctx,
			tx,
			currentSession.SessionID,
			params,
		); err != nil {
			return SessionUpdateResult{}, err
		}

		if err := tx.Commit(); err != nil {
			return SessionUpdateResult{}, fmt.Errorf(
				"failed to commit work session update: %w",
				err,
			)
		}

		currentSession.ProjectID = params.ProjectID
		currentSession.RepositoryID = params.RepositoryID
		currentSession.TaskID = params.TaskID
		currentSession.TicketKey = params.TicketKey
		currentSession.WorkspaceName = params.WorkspaceName
		currentSession.MatchStatus = params.MatchStatus
		currentSession.LastHeartbeatAt = params.OccurredAt

		return *currentSession, nil
	}

	if currentSession != nil {
		endReason := EndReasonBranchChanged

		if currentSession.RepositoryID != params.RepositoryID {
			endReason = EndReasonRepositoryChanged
		}

		if err := endActiveSession(
			ctx,
			tx,
			currentSession.SessionID,
			SessionStatusCompleted,
			endReason,
			params.OccurredAt,
		); err != nil {
			return SessionUpdateResult{}, err
		}
	}

	result, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO work_sessions (
				user_id,
				project_id,
				repository_id,
				task_id,
				branch_name,
				ticket_key,
				match_status,
				workspace_name,
				status,
				started_at,
				last_heartbeat_at,
				ended_at,
				end_reason
			)
			VALUES (
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				'active',
				?,
				?,
				NULL,
				NULL
			)
		`,
		params.UserID,
		params.ProjectID,
		params.RepositoryID,
		params.TaskID,
		params.BranchName,
		params.TicketKey,
		params.MatchStatus,
		params.WorkspaceName,
		params.OccurredAt,
		params.OccurredAt,
	)
	if err != nil {
		return SessionUpdateResult{}, fmt.Errorf(
			"failed to create work session: %w",
			err,
		)
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return SessionUpdateResult{}, fmt.Errorf(
			"failed to get work session id: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return SessionUpdateResult{}, fmt.Errorf(
			"failed to commit work session creation: %w",
			err,
		)
	}

	return SessionUpdateResult{
		SessionID:       uint64(sessionID),
		SessionCreated:  true,
		ProjectID:       params.ProjectID,
		RepositoryID:    params.RepositoryID,
		TaskID:          params.TaskID,
		BranchName:      params.BranchName,
		TicketKey:       params.TicketKey,
		WorkspaceName:   params.WorkspaceName,
		MatchStatus:     params.MatchStatus,
		Status:          SessionStatusActive,
		StartedAt:       params.OccurredAt,
		LastHeartbeatAt: params.OccurredAt,
		EndedAt:         nil,
	}, nil
}

func (r *Repository) Heartbeat(
	ctx context.Context,
	userID uint64,
	occurredAt time.Time,
) (HeartbeatResult, error) {
	tx, err := r.db.BeginTx(
		ctx,
		nil,
	)
	if err != nil {
		return HeartbeatResult{}, fmt.Errorf(
			"failed to begin heartbeat transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := lockUser(
		ctx,
		tx,
		userID,
	); err != nil {
		return HeartbeatResult{}, err
	}

	currentSession, err := findActiveSessionForUpdate(
		ctx,
		tx,
		userID,
	)
	if err != nil {
		return HeartbeatResult{}, err
	}

	if currentSession == nil {
		return HeartbeatResult{},
			ErrActiveSessionNotFound
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			UPDATE work_sessions
			SET last_heartbeat_at = ?
			WHERE id = ?
			  AND status = 'active'
		`,
		occurredAt,
		currentSession.SessionID,
	); err != nil {
		return HeartbeatResult{}, fmt.Errorf(
			"failed to update work session heartbeat: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return HeartbeatResult{}, fmt.Errorf(
			"failed to commit heartbeat update: %w",
			err,
		)
	}

	return HeartbeatResult{
		SessionID:       currentSession.SessionID,
		Status:          SessionStatusActive,
		LastHeartbeatAt: occurredAt,
	}, nil
}

func (r *Repository) Disconnect(
	ctx context.Context,
	userID uint64,
	occurredAt time.Time,
) (*DisconnectResult, error) {
	tx, err := r.db.BeginTx(
		ctx,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to begin disconnect transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := lockUser(
		ctx,
		tx,
		userID,
	); err != nil {
		return nil, err
	}

	currentSession, err := findActiveSessionForUpdate(
		ctx,
		tx,
		userID,
	)
	if err != nil {
		return nil, err
	}

	if currentSession == nil {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf(
				"failed to commit empty disconnect: %w",
				err,
			)
		}

		return nil, nil
	}

	if err := endActiveSession(
		ctx,
		tx,
		currentSession.SessionID,
		SessionStatusCompleted,
		EndReasonClientClosed,
		occurredAt,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(
			"failed to commit disconnect: %w",
			err,
		)
	}

	return &DisconnectResult{
		SessionID: currentSession.SessionID,
		Status:    SessionStatusCompleted,
		EndReason: EndReasonClientClosed,
		EndedAt:   occurredAt,
	}, nil
}

func lockUser(
	ctx context.Context,
	tx *sql.Tx,
	userID uint64,
) error {
	var lockedUserID uint64

	err := tx.QueryRowContext(
		ctx,
		`
			SELECT id
			FROM users
			WHERE id = ?
			FOR UPDATE
		`,
		userID,
	).Scan(&lockedUserID)
	if err != nil {
		return fmt.Errorf(
			"failed to lock work session user: %w",
			err,
		)
	}

	return nil
}

func findActiveSessionForUpdate(
	ctx context.Context,
	tx *sql.Tx,
	userID uint64,
) (*SessionUpdateResult, error) {
	var session SessionUpdateResult
	var taskID sql.NullInt64
	var ticketKey sql.NullString
	var workspaceName sql.NullString
	var endedAt sql.NullTime

	err := tx.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				project_id,
				repository_id,
				task_id,
				branch_name,
				ticket_key,
				workspace_name,
				match_status,
				status,
				started_at,
				last_heartbeat_at,
				ended_at
			FROM work_sessions
			WHERE user_id = ?
			  AND status = 'active'
			ORDER BY id DESC
			LIMIT 1
			FOR UPDATE
		`,
		userID,
	).Scan(
		&session.SessionID,
		&session.ProjectID,
		&session.RepositoryID,
		&taskID,
		&session.BranchName,
		&ticketKey,
		&workspaceName,
		&session.MatchStatus,
		&session.Status,
		&session.StartedAt,
		&session.LastHeartbeatAt,
		&endedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf(
			"failed to find active work session: %w",
			err,
		)
	}

	if taskID.Valid {
		value := uint64(taskID.Int64)
		session.TaskID = &value
	}

	if ticketKey.Valid {
		value := ticketKey.String
		session.TicketKey = &value
	}

	if workspaceName.Valid {
		value := workspaceName.String
		session.WorkspaceName = &value
	}

	if endedAt.Valid {
		value := endedAt.Time
		session.EndedAt = &value
	}

	return &session, nil
}

func updateExistingSessionContext(
	ctx context.Context,
	tx *sql.Tx,
	sessionID uint64,
	params StartSessionParams,
) error {
	_, err := tx.ExecContext(
		ctx,
		`
			UPDATE work_sessions
			SET
				project_id = ?,
				repository_id = ?,
				task_id = ?,
				ticket_key = ?,
				match_status = ?,
				workspace_name = ?,
				last_heartbeat_at = ?
			WHERE id = ?
			  AND status = 'active'
		`,
		params.ProjectID,
		params.RepositoryID,
		params.TaskID,
		params.TicketKey,
		params.MatchStatus,
		params.WorkspaceName,
		params.OccurredAt,
		sessionID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update existing work session: %w",
			err,
		)
	}

	return nil
}

func endActiveSession(
	ctx context.Context,
	tx *sql.Tx,
	sessionID uint64,
	status string,
	endReason string,
	endedAt time.Time,
) error {
	_, err := tx.ExecContext(
		ctx,
		`
			UPDATE work_sessions
			SET
				status = ?,
				ended_at = ?,
				end_reason = ?
			WHERE id = ?
			  AND status = 'active'
		`,
		status,
		endedAt,
		endReason,
		sessionID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to end active work session: %w",
			err,
		)
	}

	return nil
}
