package company

import (
	"errors"
	"log"
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

func (h *Handler) Create(c *gin.Context) {
	var request CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "会社情報を正しく入力してください。",
		})
		return
	}

	response, err := h.service.Create(
		c.Request.Context(),
		request,
	)
	if err != nil {
		if errors.Is(err, ErrInvalidSlug) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "会社スラッグの形式が正しくありません。",
			})
			return
		}

		var mysqlError *mysql.MySQLError

		if errors.As(err, &mysqlError) &&
			mysqlError.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "この会社スラッグは既に使用されています。",
			})
			return
		}

		log.Printf("failed to create company: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "会社登録に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
