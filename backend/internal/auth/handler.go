package auth

import (
	"errors"
	"log"
	"net/http"

	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Login godoc
// @Summary ログイン
// @Description メールアドレスとパスワードで認証し、アクセストークンとユーザー情報を返します。
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "ログイン情報"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "メールアドレスとパスワードを正しく入力してください。",
		})
		return
	}

	response, err := h.service.Login(
		c.Request.Context(),
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "メールアドレスまたはパスワードが正しくありません。",
			})
			return

		case errors.Is(err, ErrJWTSecretNotSet):
			log.Printf(
				"failed to issue access token: JWT_SECRET is not configured",
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "ログイン処理に失敗しました。",
			})
			return
		}

		log.Printf(
			"failed to login user: email=%s error=%v",
			request.Email,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ログイン処理に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CurrentUser godoc
// @Summary ログインユーザー取得
// @Description アクセストークンから現在のログインユーザーと会社情報を取得します。
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CurrentUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/me [get]
func (h *Handler) CurrentUser(c *gin.Context) {
	userIDValue, userIDExists := c.Get(
		middleware.ContextUserIDKey,
	)

	companyIDValue, companyIDExists := c.Get(
		middleware.ContextCompanyIDKey,
	)

	if !userIDExists || !companyIDExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "ログイン情報が取得できません。",
		})
		return
	}

	userID, userIDValid := userIDValue.(uint64)
	companyID, companyIDValid := companyIDValue.(uint64)

	if !userIDValid ||
		!companyIDValid ||
		userID == 0 ||
		companyID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "ログイン情報が無効です。",
		})
		return
	}

	response, err := h.service.CurrentUser(
		c.Request.Context(),
		userID,
		companyID,
	)
	if err != nil {
		if errors.Is(err, ErrCurrentUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "ログイン情報が無効です。",
			})
			return
		}

		log.Printf(
			"failed to get current user: user_id=%d company_id=%d error=%v",
			userID,
			companyID,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ログインユーザー情報の取得に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
