package auth

type ErrorResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        LoginUser    `json:"user"`
	Company     LoginCompany `json:"company"`
}

type LoginUser struct {
	ID          uint64 `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

type LoginCompany struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CurrentUserResponse struct {
	User    CurrentUserCompanyMember `json:"user"`
	Company CurrentUserCompany       `json:"company"`
}

type CurrentUserCompanyMember struct {
	ID          uint64 `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

type CurrentUserCompany struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
