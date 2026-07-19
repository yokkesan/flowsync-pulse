package user

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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

func (h *Handler) Register(c *gin.Context) {
	companyID, err := strconv.ParseUint(
		c.Param("companyId"),
		10,
		64,
	)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "会社IDが正しくありません。",
		})
		return
	}

	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ユーザー情報を正しく入力してください。",
		})
		return
	}

	response, err := h.service.Register(
		c.Request.Context(),
		companyID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrPasswordMismatch):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "パスワード確認が一致しません。",
			})
			return

		case errors.Is(err, ErrCompanyNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"message": "登録先の会社が見つかりません。",
			})
			return
		}

		var mysqlError *mysql.MySQLError

		if errors.As(err, &mysqlError) &&
			mysqlError.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "このメールアドレスは既に登録されています。",
			})
			return
		}

		log.Printf(
			"failed to register user: company_id=%d error=%v",
			companyID,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ユーザー登録に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
