package realtime

import "time"

const (
	EventWorkContextChanged = "work_context.changed"
	EventPresenceChanged    = "presence.changed"
)

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type WorkContextChangedPayload struct {
	UserID          uint64     `json:"user_id"`
	ProjectID       *uint64    `json:"project_id"`
	ProjectName     *string    `json:"project_name"`
	TaskID          *uint64    `json:"task_id"`
	TaskKey         *string    `json:"task_key"`
	TaskName        *string    `json:"task_name"`
	RepositoryName  *string    `json:"repository_name"`
	BranchName      *string    `json:"branch_name"`
	TicketKey       *string    `json:"ticket_key"`
	WorkspaceName   *string    `json:"workspace_name"`
	MatchStatus     *string    `json:"match_status"`
	SessionStatus   *string    `json:"session_status"`
	ExtensionActive bool       `json:"extension_active"`
	StartedAt       *time.Time `json:"started_at"`
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at"`
	EndedAt         *time.Time `json:"ended_at"`
	EndReason       *string    `json:"end_reason"`
}

type PresenceChangedPayload struct {
	UserID          uint64 `json:"user_id"`
	ExtensionActive bool   `json:"extension_active"`
}
