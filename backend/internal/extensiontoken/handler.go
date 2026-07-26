package extensiontoken

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

// Create godoc
// @Summary 拡張機能トークン発行
// @Description ログインユーザーに紐づくVS Code拡張機能専用トークンを発行します。平文トークンは発行時のみ返します。
// @Tags extension-tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRequest true "拡張機能トークン情報"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension-tokens [post]
func (h *Handler) Create(
	c *gin.Context,
) {
	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "ログイン情報が無効です。",
			},
		)
		return
	}

	var request CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "端末名を正しく入力してください。",
			},
		)
		return
	}

	response, err := h.service.Create(
		c.Request.Context(),
		authUserID,
		companyID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(
			err,
			ErrInvalidExtensionTokenName,
		):
			c.JSON(
				http.StatusBadRequest,
				ErrorResponse{
					Message: "端末名は1文字以上100文字以下で入力してください。",
				},
			)
			return

		default:
			log.Printf(
				"failed to create extension token: user_id=%d company_id=%d error=%v",
				authUserID,
				companyID,
				err,
			)

			c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{
					Message: "拡張機能トークンの発行に失敗しました。",
				},
			)
			return
		}
	}

	log.Printf(
		"extension token created: token_id=%d user_id=%d company_id=%d",
		response.TokenID,
		authUserID,
		companyID,
	)

	c.JSON(
		http.StatusCreated,
		response,
	)
}

// List godoc
// @Summary 拡張機能トークン一覧取得
// @Description ログインユーザーが発行した拡張機能専用トークンの一覧を取得します。平文トークンとハッシュ値は返しません。
// @Tags extension-tokens
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension-tokens [get]
func (h *Handler) List(
	c *gin.Context,
) {
	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "ログイン情報が無効です。",
			},
		)
		return
	}

	response, err := h.service.List(
		c.Request.Context(),
		authUserID,
		companyID,
	)
	if err != nil {
		log.Printf(
			"failed to list extension tokens: user_id=%d company_id=%d error=%v",
			authUserID,
			companyID,
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Message: "拡張機能トークン一覧の取得に失敗しました。",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Revoke godoc
// @Summary 拡張機能トークン失効
// @Description ログインユーザーが発行した拡張機能専用トークンを失効させます。同じトークンを複数回失効してもエラーになりません。
// @Tags extension-tokens
// @Produce json
// @Security BearerAuth
// @Param tokenId path int true "拡張機能トークンID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension-tokens/{tokenId} [delete]
func (h *Handler) Revoke(
	c *gin.Context,
) {
	tokenID, err := strconv.ParseUint(
		c.Param("tokenId"),
		10,
		64,
	)
	if err != nil || tokenID == 0 {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "トークンIDが正しくありません。",
			},
		)
		return
	}

	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "ログイン情報が無効です。",
			},
		)
		return
	}

	err = h.service.Revoke(
		c.Request.Context(),
		tokenID,
		authUserID,
		companyID,
	)
	if err != nil {
		switch {
		case errors.Is(
			err,
			ErrExtensionTokenNotFound,
		):
			c.JSON(
				http.StatusNotFound,
				ErrorResponse{
					Message: "対象の拡張機能トークンが存在しません。",
				},
			)
			return

		default:
			log.Printf(
				"failed to revoke extension token: token_id=%d user_id=%d company_id=%d error=%v",
				tokenID,
				authUserID,
				companyID,
				err,
			)

			c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{
					Message: "拡張機能トークンの失効に失敗しました。",
				},
			)
			return
		}
	}

	log.Printf(
		"extension token revoked: token_id=%d user_id=%d company_id=%d",
		tokenID,
		authUserID,
		companyID,
	)

	c.Status(http.StatusNoContent)
}
