package project

import "time"

const (
	StatusPlanned   = "planned"
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusArchived  = "archived"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=150"`
	Slug          string   `json:"slug" binding:"required,min=2,max=100"`
	RepositoryURL string   `json:"repository_url" binding:"required,max=500"`
	Description   *string  `json:"description" binding:"omitempty,max=5000"`
	Status        string   `json:"status" binding:"required,oneof=planned active completed archived"`
	StartDate     *string  `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate       *string  `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	MemberIDs     []uint64 `json:"member_ids" binding:"required,min=1,dive,gt=0"`
}

type UpdateRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=150"`
	Slug          string   `json:"slug" binding:"required,min=2,max=100"`
	ProjectKey    *string  `json:"project_key" binding:"omitempty,min=2,max=10"`
	RepositoryURL string   `json:"repository_url" binding:"required,max=500"`
	Description   *string  `json:"description" binding:"omitempty,max=5000"`
	Status        string   `json:"status" binding:"required,oneof=planned active completed archived"`
	StartDate     *string  `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate       *string  `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	MemberIDs     []uint64 `json:"member_ids" binding:"required,min=1,dive,gt=0"`
}

type MemberResponse struct {
	UserID      uint64 `json:"user_id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}

type Response struct {
	ProjectID     uint64           `json:"project_id"`
	CompanyID     uint64           `json:"company_id"`
	Name          string           `json:"name"`
	Slug          string           `json:"slug"`
	ProjectKey    *string          `json:"project_key"`
	RepositoryURL *string          `json:"repository_url"`
	Description   *string          `json:"description"`
	Status        string           `json:"status"`
	StartDate     *string          `json:"start_date"`
	EndDate       *string          `json:"end_date"`
	Members       []MemberResponse `json:"members"`
	TaskCount     uint64           `json:"task_count"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type ListResponse struct {
	Projects []Response `json:"projects"`
}
