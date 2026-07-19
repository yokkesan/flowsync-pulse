package auth

import (
	"errors"
	"log"
	"net/http"

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
