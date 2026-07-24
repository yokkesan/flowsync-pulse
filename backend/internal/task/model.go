package task

import "time"

const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusSuspended  = "suspended"
)

const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=255"`
	Description    *string `json:"description" binding:"omitempty,max=5000"`
	AssigneeUserID uint64  `json:"assignee_user_id" binding:"required"`
	BranchName     string  `json:"branch_name" binding:"required,min=1,max=255"`
	Status         string  `json:"status" binding:"required,oneof=not_started in_progress completed suspended"`
	Priority       string  `json:"priority" binding:"required,oneof=high medium low"`
	StartDate      *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	DueDate        *string `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=255"`
	Description    *string `json:"description" binding:"omitempty,max=5000"`
	AssigneeUserID uint64  `json:"assignee_user_id" binding:"required"`
	BranchName     string  `json:"branch_name" binding:"required,min=1,max=255"`
	Status         string  `json:"status" binding:"required,oneof=not_started in_progress completed suspended"`
	Priority       string  `json:"priority" binding:"required,oneof=high medium low"`
	StartDate      *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	DueDate        *string `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
}

type Response struct {
	TaskID         uint64     `json:"task_id"`
	ProjectID      uint64     `json:"project_id"`
	ProjectName    string     `json:"project_name"`
	Name           string     `json:"name"`
	Description    *string    `json:"description"`
	AssigneeUserID uint64     `json:"assignee_user_id"`
	AssigneeName   string     `json:"assignee_name"`
	BranchName     string     `json:"branch_name"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	StartDate      *string    `json:"start_date"`
	DueDate        *string    `json:"due_date"`
	CompletedAt    *time.Time `json:"completed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ListResponse struct {
	Tasks []Response `json:"tasks"`
}
