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
// @Summary 拡張機能の作業コンテキスト送信
// @Description VS Code拡張機能からリポジトリ・ブランチ・チケット情報を受信し、現在の作業セッションを開始または更新します。
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
				Message: "拡張機能トークンが無効です。",
			},
		)
		return
	}

	var request WorkContextRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "作業情報を正しく入力してください。",
			},
		)
		return
	}

	response, err := h.service.WorkContext(
		c.Request.Context(),
		authUserID,
		companyID,
		request,
	)
	if err != nil {
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
			return

		case errors.Is(
			err,
			ErrInvalidOccurredAt,
		):
			c.JSON(
				http.StatusBadRequest,
				ErrorResponse{
					Message: "情報取得日時が正しくありません。",
				},
			)
			return

		case errors.Is(
			err,
			ErrRepositoryNotFound,
		):
			c.JSON(
				http.StatusNotFound,
				ErrorResponse{
					Message: "対象のリポジトリが登録されていません。",
				},
			)
			return

		case errors.Is(
			err,
			ErrExtensionProjectAccessDenied,
		):
			log.Printf(
				"extension work context authorization denied: user_id=%d company_id=%d",
				authUserID,
				companyID,
			)

			c.JSON(
				http.StatusForbidden,
				ErrorResponse{
					Message: "対象プロジェクトへアクセスする権限がありません。",
				},
			)
			return

		default:
			log.Printf(
				"failed to save extension work context: user_id=%d company_id=%d error=%v",
				authUserID,
				companyID,
				err,
			)

			c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{
					Message: "作業セッションの保存に失敗しました。",
				},
			)
			return
		}
	}

	log.Printf(
		"extension work context saved: session_id=%d user_id=%d company_id=%d project_id=%d match_status=%s",
		response.SessionID,
		authUserID,
		companyID,
		response.ProjectID,
		response.MatchStatus,
	)

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Heartbeat godoc
// @Summary 拡張機能Heartbeat送信
// @Description 現在進行中の作業セッションの最終Heartbeat日時を更新します。
// @Tags extension
// @Accept json
// @Produce json
// @Security ExtensionTokenAuth
// @Param request body HeartbeatRequest true "Heartbeat情報"
// @Success 200 {object} HeartbeatResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension/heartbeat [post]
func (h *Handler) Heartbeat(
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
				Message: "拡張機能トークンが無効です。",
			},
		)
		return
	}

	var request HeartbeatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "Heartbeat情報を正しく入力してください。",
			},
		)
		return
	}

	response, err := h.service.Heartbeat(
		c.Request.Context(),
		authUserID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(
			err,
			ErrInvalidOccurredAt,
		):
			c.JSON(
				http.StatusBadRequest,
				ErrorResponse{
					Message: "Heartbeat日時が正しくありません。",
				},
			)
			return

		case errors.Is(
			err,
			ErrActiveSessionNotFound,
		):
			c.JSON(
				http.StatusNotFound,
				ErrorResponse{
					Message: "進行中の作業セッションが存在しません。",
				},
			)
			return

		default:
			log.Printf(
				"failed to update extension heartbeat: user_id=%d company_id=%d error=%v",
				authUserID,
				companyID,
				err,
			)

			c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{
					Message: "Heartbeatの更新に失敗しました。",
				},
			)
			return
		}
	}

	log.Printf(
		"extension heartbeat updated: session_id=%d user_id=%d company_id=%d",
		response.SessionID,
		authUserID,
		companyID,
	)

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Disconnect godoc
// @Summary 拡張機能終了通知
// @Description VS Code終了時などに現在進行中の作業セッションを終了します。既に終了済みの場合も正常終了とします。
// @Tags extension
// @Accept json
// @Produce json
// @Security ExtensionTokenAuth
// @Param request body DisconnectRequest true "終了通知情報"
// @Success 200 {object} DisconnectResponse
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/extension/disconnect [post]
func (h *Handler) Disconnect(
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
				Message: "拡張機能トークンが無効です。",
			},
		)
		return
	}

	var request DisconnectRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "終了通知情報を正しく入力してください。",
			},
		)
		return
	}

	response, err := h.service.Disconnect(
		c.Request.Context(),
		authUserID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(
			err,
			ErrInvalidOccurredAt,
		):
			c.JSON(
				http.StatusBadRequest,
				ErrorResponse{
					Message: "終了日時が正しくありません。",
				},
			)
			return

		default:
			log.Printf(
				"failed to disconnect extension session: user_id=%d company_id=%d error=%v",
				authUserID,
				companyID,
				err,
			)

			c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{
					Message: "作業セッションの終了に失敗しました。",
				},
			)
			return
		}
	}

	// 進行中のセッションが存在しない場合も冪等な成功とする。
	if response == nil {
		c.Status(
			http.StatusNoContent,
		)
		return
	}

	log.Printf(
		"extension session disconnected: session_id=%d user_id=%d company_id=%d",
		response.SessionID,
		authUserID,
		companyID,
	)

	c.JSON(
		http.StatusOK,
		response,
	)
}
