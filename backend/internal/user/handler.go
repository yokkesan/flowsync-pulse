package user

import (
	"errors"
	"flowsync-pulse/backend/internal/middleware"
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

// Register godoc
// @Summary ユーザー登録
// @Description 指定した会社に所属するユーザーを登録します。
// @Tags users
// @Accept json
// @Produce json
// @Param companyId path int true "会社ID"
// @Param request body RegisterRequest true "ユーザー情報"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies/{companyId}/users [post]
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

// List godoc
// @Summary 会社所属ユーザー一覧取得
// @Description ログインユーザーと同じ会社に所属する有効ユーザーを取得します。
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users [get]
func (h *Handler) List(c *gin.Context) {
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if companyID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "ログイン情報が無効です。",
		})
		return
	}

	response, err := h.service.ListByCompanyID(
		c.Request.Context(),
		companyID,
	)
	if err != nil {
		log.Printf(
			"failed to list company users: company_id=%d error=%v",
			companyID,
			err,
		)

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "所属ユーザー一覧の取得に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
