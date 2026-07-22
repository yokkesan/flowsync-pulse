package user

type ErrorResponse struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	DisplayName     string `json:"display_name" binding:"required,min=2,max=100"`
	Email           string `json:"email" binding:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	PasswordConfirm string `json:"password_confirmation" binding:"required"`
}

type RegisterResponse struct {
	UserID      uint64 `json:"user_id"`
	CompanyID   uint64 `json:"company_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}
