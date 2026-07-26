package extension

import "time"

const (
	MatchStatusMatched              = "matched"
	MatchStatusTicketNotFound       = "ticket_not_found"
	MatchStatusBranchNotMatched     = "branch_not_matched"
	MatchStatusTicketBranchMismatch = "ticket_branch_mismatch"
)

const (
	SessionStatusActive    = "active"
	SessionStatusCompleted = "completed"
	SessionStatusTimedOut  = "timed_out"
)

const (
	EndReasonBranchChanged     = "branch_changed"
	EndReasonRepositoryChanged = "repository_changed"
	EndReasonClientClosed      = "client_closed"
	EndReasonTimeout           = "timeout"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type WorkContextRequest struct {
	RepositoryName   string    `json:"repository_name" binding:"required,min=1,max=255"`
	RepositoryURL    string    `json:"repository_url" binding:"required,min=1,max=2048"`
	BranchName       string    `json:"branch_name" binding:"required,min=1,max=255"`
	TicketKey        *string   `json:"ticket_key" binding:"omitempty,max=100"`
	WorkspaceName    *string   `json:"workspace_name" binding:"omitempty,max=255"`
	ExtensionVersion *string   `json:"extension_version" binding:"omitempty,max=50"`
	OccurredAt       time.Time `json:"occurred_at" binding:"required"`
}

type WorkContextResponse struct {
	SessionID        uint64     `json:"session_id"`
	ProjectID        uint64     `json:"project_id"`
	ProjectName      string     `json:"project_name"`
	TaskID           *uint64    `json:"task_id"`
	TaskKey          *string    `json:"task_key"`
	TaskName         *string    `json:"task_name"`
	RepositoryName   string     `json:"repository_name"`
	RepositoryURL    string     `json:"repository_url"`
	BranchName       string     `json:"branch_name"`
	TicketKey        *string    `json:"ticket_key"`
	WorkspaceName    *string    `json:"workspace_name"`
	ExtensionVersion *string    `json:"extension_version"`
	MatchStatus      string     `json:"match_status"`
	Status           string     `json:"status"`
	StartedAt        time.Time  `json:"started_at"`
	LastHeartbeatAt  time.Time  `json:"last_heartbeat_at"`
	EndedAt          *time.Time `json:"ended_at"`
}

type HeartbeatRequest struct {
	OccurredAt time.Time `json:"occurred_at" binding:"required"`
}

type HeartbeatResponse struct {
	SessionID       uint64    `json:"session_id"`
	Status          string    `json:"status"`
	LastHeartbeatAt time.Time `json:"last_heartbeat_at"`
}

type DisconnectRequest struct {
	OccurredAt time.Time `json:"occurred_at" binding:"required"`
}

type DisconnectResponse struct {
	SessionID uint64     `json:"session_id"`
	Status    string     `json:"status"`
	EndReason *string    `json:"end_reason"`
	EndedAt   *time.Time `json:"ended_at"`
}
