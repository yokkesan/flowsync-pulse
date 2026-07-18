package user

type RegisterRequest struct {
	CompanyName     string `json:"company_name" binding:"required,min=2,max=150"`
	CompanySlug     string `json:"company_slug" binding:"required,min=2,max=100"`
	DisplayName     string `json:"display_name" binding:"required,min=2,max=100"`
	Email           string `json:"email" binding:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	PasswordConfirm string `json:"password_confirmation" binding:"required"`
	AvatarKey       string `json:"avatar_key" binding:"omitempty,max=100"`
}

type RegisterResponse struct {
	UserID      uint64 `json:"user_id"`
	CompanyID   uint64 `json:"company_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Role        string `json:"role"`
}
