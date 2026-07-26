package extensiontoken

import "time"

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type CreateResponse struct {
	TokenID   uint64    `json:"token_id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	TokenID    uint64     `json:"token_id"`
	Name       string     `json:"name"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ListResponse struct {
	Tokens []Response `json:"tokens"`
}
