package token

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	CompanyID uint64 `json:"company_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}
