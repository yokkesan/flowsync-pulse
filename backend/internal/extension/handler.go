package extension

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

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

// WorkContext godoc
// @Summary 作業コンテキスト送信
// @Description VSCode拡張機能から、現在のリポジトリ・ブランチ・タスク情報を送信します。
// @Tags extension
// @Accept json
// @Produce json
// @Security ExtensionTokenAuth
// @Param request body WorkContextRequest true "作業コンテキスト"
// @Success 200 {object} WorkContextResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension/work-context [post]
func (h *Handler) WorkContext(
	c *gin.Context,
) {
	var request WorkContextRequest

	if err := c.ShouldBindJSON(
		&request,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "入力内容が正しくありません。",
			},
		)
		return
	}

	userID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if userID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "拡張機能の認証情報が無効です。",
			},
		)
		return
	}

	response, err := h.service.WorkContext(
		c.Request.Context(),
		userID,
		companyID,
		request,
	)
	if err != nil {
		h.handleError(
			c,
			err,
			"failed to save extension work context",
			userID,
			companyID,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Heartbeat godoc
// @Summary 作業セッションのハートビート
// @Description 現在の作業セッションが継続中であることを通知します。
// @Tags extension
// @Accept json
// @Produce json
// @Security ExtensionTokenAuth
// @Param request body HeartbeatRequest true "ハートビート"
// @Success 200 {object} HeartbeatResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension/heartbeat [post]
func (h *Handler) Heartbeat(
	c *gin.Context,
) {
	var request HeartbeatRequest

	if err := c.ShouldBindJSON(
		&request,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "入力内容が正しくありません。",
			},
		)
		return
	}

	userID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if userID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "拡張機能の認証情報が無効です。",
			},
		)
		return
	}

	response, err := h.service.Heartbeat(
		c.Request.Context(),
		userID,
		companyID,
		request,
	)
	if err != nil {
		h.handleError(
			c,
			err,
			"failed to update extension heartbeat",
			userID,
			companyID,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Disconnect godoc
// @Summary 作業セッション終了
// @Description VSCode拡張機能の終了時に、現在の作業セッションを終了します。
// @Tags extension
// @Accept json
// @Produce json
// @Security ExtensionTokenAuth
// @Param request body DisconnectRequest true "切断情報"
// @Success 200 {object} DisconnectResponse
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension/disconnect [post]
func (h *Handler) Disconnect(
	c *gin.Context,
) {
	var request DisconnectRequest

	if err := c.ShouldBindJSON(
		&request,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "入力内容が正しくありません。",
			},
		)
		return
	}

	userID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if userID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "拡張機能の認証情報が無効です。",
			},
		)
		return
	}

	response, err := h.service.Disconnect(
		c.Request.Context(),
		userID,
		companyID,
		request,
	)
	if err != nil {
		h.handleError(
			c,
			err,
			"failed to disconnect extension session",
			userID,
			companyID,
		)
		return
	}

	if response == nil {
		c.Status(
			http.StatusNoContent,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

func (h *Handler) handleError(
	c *gin.Context,
	err error,
	logMessage string,
	userID uint64,
	companyID uint64,
) {
	switch {
	case errors.Is(
		err,
		ErrInvalidRepositoryURL,
	):
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "リポジトリURLが正しくありません。",
			},
		)

	case errors.Is(
		err,
		ErrInvalidOccurredAt,
	):
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "発生日時が正しくありません。",
			},
		)

	case errors.Is(
		err,
		ErrRepositoryNotFound,
	):
		c.JSON(
			http.StatusNotFound,
			ErrorResponse{
				Message: "登録済みリポジトリが見つかりません。",
			},
		)

	case errors.Is(
		err,
		ErrExtensionProjectAccessDenied,
	):
		c.JSON(
			http.StatusForbidden,
			ErrorResponse{
				Message: "対象プロジェクトへアクセスする権限がありません。",
			},
		)

	case errors.Is(
		err,
		ErrActiveSessionNotFound,
	):
		c.JSON(
			http.StatusNotFound,
			ErrorResponse{
				Message: "有効な作業セッションが見つかりません。",
			},
		)

	default:
		log.Printf(
			"%s: user_id=%d company_id=%d error=%v",
			logMessage,
			userID,
			companyID,
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Message: "サーバー内部でエラーが発生しました。",
			},
		)
	}
}
