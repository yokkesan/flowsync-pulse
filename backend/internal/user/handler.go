package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterCompanyOwner(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "入力内容を確認してください。",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.service.RegisterCompanyOwner(
		c.Request.Context(),
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrPasswordMismatch):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "パスワード確認が一致しません。",
			})
			return

		case errors.Is(err, ErrInvalidCompanySlug):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "会社スラッグの形式が正しくありません。",
			})
			return
		}

		var mysqlError *mysql.MySQLError

		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "メールアドレスまたは会社スラッグが既に登録されています。",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ユーザー登録に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
