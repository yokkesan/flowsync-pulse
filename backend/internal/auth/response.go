package auth

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
