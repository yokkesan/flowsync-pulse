package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"flowsync-pulse/backend/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserIDKey    = "auth_user_id"
	ContextCompanyIDKey = "auth_company_id"
	ContextRoleKey      = "auth_role"
)

func Authenticate() gin.HandlerFunc {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	return func(c *gin.Context) {
		authorization := strings.TrimSpace(
			c.GetHeader("Authorization"),
		)

		tokenString, ok := extractBearerToken(authorization)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "ログインが必要です。",
				},
			)
			return
		}

		claims := &token.AccessTokenClaims{}

		parsedToken, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(parsedToken *jwt.Token) (any, error) {
				if parsedToken.Method != jwt.SigningMethodHS256 {
					return nil, errors.New(
						"unexpected signing method",
					)
				}

				return jwtSecret, nil
			},
			jwt.WithIssuer("flowsync-pulse"),
			jwt.WithExpirationRequired(),
			jwt.WithValidMethods([]string{
				jwt.SigningMethodHS256.Alg(),
			}),
		)
		if err != nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "ログイン情報が無効です。",
				},
			)
			return
		}

		userID, err := strconv.ParseUint(
			claims.Subject,
			10,
			64,
		)
		if err != nil ||
			userID == 0 ||
			claims.CompanyID == 0 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "ログイン情報が無効です。",
				},
			)
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Set(ContextCompanyIDKey, claims.CompanyID)
		c.Set(ContextRoleKey, claims.Role)

		c.Next()
	}
}

func extractBearerToken(
	authorization string,
) (string, bool) {
	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(authorization, bearerPrefix) {
		return "", false
	}

	tokenString := strings.TrimSpace(
		strings.TrimPrefix(
			authorization,
			bearerPrefix,
		),
	)

	if tokenString == "" {
		return "", false
	}

	return tokenString, true
}
