package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextExtensionTokenIDKey = "extension_token_id"
)

type ExtensionTokenAuthenticator interface {
	FindValidByHash(
		ctx context.Context,
		tokenHash string,
		now time.Time,
	) (
		uint64,
		uint64,
		uint64,
		error,
	)

	UpdateLastUsedAt(
		ctx context.Context,
		tokenID uint64,
		usedAt time.Time,
	) error
}

func AuthenticateExtension(
	authenticator ExtensionTokenAuthenticator,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(
			c.GetHeader("Authorization"),
		)

		plainToken, ok := extractBearerToken(
			authorization,
		)
		if !ok {
			log.Printf(
				"extension authentication failed: reason=token_missing client_ip=%s",
				c.ClientIP(),
			)

			abortInvalidExtensionToken(c)
			return
		}

		tokenHash := hashExtensionToken(
			plainToken,
		)

		now := time.Now().UTC()

		tokenID,
			userID,
			companyID,
			err := authenticator.FindValidByHash(
			c.Request.Context(),
			tokenHash,
			now,
		)
		if err != nil {
			log.Printf(
				"extension authentication failed: reason=token_invalid client_ip=%s error=%v",
				c.ClientIP(),
				err,
			)

			abortInvalidExtensionToken(c)
			return
		}

		if tokenID == 0 ||
			userID == 0 ||
			companyID == 0 {
			log.Printf(
				"extension authentication failed: reason=invalid_identity client_ip=%s",
				c.ClientIP(),
			)

			abortInvalidExtensionToken(c)
			return
		}

		if err := authenticator.UpdateLastUsedAt(
			c.Request.Context(),
			tokenID,
			now,
		); err != nil {
			log.Printf(
				"extension authentication failed: reason=last_used_update_failed token_id=%d user_id=%d company_id=%d error=%v",
				tokenID,
				userID,
				companyID,
				err,
			)

			abortInvalidExtensionToken(c)
			return
		}

		c.Set(
			ContextExtensionTokenIDKey,
			tokenID,
		)
		c.Set(
			ContextUserIDKey,
			userID,
		)
		c.Set(
			ContextCompanyIDKey,
			companyID,
		)

		c.Next()
	}
}

func hashExtensionToken(
	plainToken string,
) string {
	hashedToken := sha256.Sum256(
		[]byte(plainToken),
	)

	return hex.EncodeToString(
		hashedToken[:],
	)
}

func abortInvalidExtensionToken(
	c *gin.Context,
) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{
			"message": "拡張機能トークンが無効です。",
		},
	)
}
