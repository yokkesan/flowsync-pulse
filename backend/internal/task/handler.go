package task

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

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Create godoc
// @Summary タスク登録
// @Description 対象プロジェクトに新しいタスクとブランチ情報を登録します。
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Param request body CreateRequest true "タスク情報"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId}/tasks [post]
func (h *Handler) Create(c *gin.Context) {
	projectID, err := strconv.ParseUint(
		c.Param("projectId"),
		10,
		64,
	)
	if err != nil || projectID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "プロジェクトIDが正しくありません。",
		})
		return
	}

	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "ログイン情報が無効です。",
		})
		return
	}

	var request CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "タスク情報を正しく入力してください。",
		})
		return
	}

	response, err := h.service.Create(
		c.Request.Context(),
		projectID,
		authUserID,
		companyID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"task create authorization denied: user_id=%d company_id=%d project_id=%d",
				authUserID,
				companyID,
				projectID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトへアクセスする権限がありません。",
			})
			return

		case errors.Is(err, ErrAssigneeNotProjectMember):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "担当メンバーが対象プロジェクトに所属していません。",
			})
			return

		case errors.Is(err, ErrBranchAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Message: "このブランチ名は対象プロジェクトで既に使用されています。",
			})
			return

		case errors.Is(err, ErrInvalidDateRange):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "期限は開始日以降の日付を指定してください。",
			})
			return

		default:
			log.Printf(
				"failed to create task: user_id=%d company_id=%d project_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "タスク登録に失敗しました。",
			})
			return
		}
	}

	log.Printf(
		"task created: task_id=%d user_id=%d company_id=%d project_id=%d",
		response.TaskID,
		authUserID,
		companyID,
		projectID,
	)

	c.JSON(http.StatusCreated, response)
}

// List godoc
// @Summary プロジェクトのタスク一覧取得
// @Description 対象プロジェクトに登録されているタスク一覧を取得します。
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Success 200 {object} ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId}/tasks [get]
func (h *Handler) List(c *gin.Context) {
	projectID, err := strconv.ParseUint(
		c.Param("projectId"),
		10,
		64,
	)
	if err != nil || projectID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "プロジェクトIDが正しくありません。",
		})
		return
	}

	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "ログイン情報が無効です。",
		})
		return
	}

	response, err := h.service.List(
		c.Request.Context(),
		projectID,
		authUserID,
		companyID,
	)
	if err != nil {
		if errors.Is(err, ErrProjectAccessDenied) {
			log.Printf(
				"task list authorization denied: user_id=%d company_id=%d project_id=%d",
				authUserID,
				companyID,
				projectID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトへアクセスする権限がありません。",
			})
			return
		}

		log.Printf(
			"failed to list tasks: user_id=%d company_id=%d project_id=%d error=%v",
			authUserID,
			companyID,
			projectID,
			err,
		)

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "タスク一覧の取得に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get godoc
// @Summary タスク詳細取得
// @Description 対象プロジェクトに登録されているタスク詳細を取得します。
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Param taskId path int true "タスクID"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId}/tasks/{taskId} [get]
func (h *Handler) Get(c *gin.Context) {
	projectID, err := strconv.ParseUint(
		c.Param("projectId"),
		10,
		64,
	)
	if err != nil || projectID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "プロジェクトIDが正しくありません。",
		})
		return
	}

	taskID, err := strconv.ParseUint(
		c.Param("taskId"),
		10,
		64,
	)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "タスクIDが正しくありません。",
		})
		return
	}

	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "ログイン情報が無効です。",
		})
		return
	}

	response, err := h.service.Get(
		c.Request.Context(),
		projectID,
		taskID,
		authUserID,
		companyID,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"task detail authorization denied: user_id=%d company_id=%d project_id=%d task_id=%d",
				authUserID,
				companyID,
				projectID,
				taskID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトへアクセスする権限がありません。",
			})
			return

		case errors.Is(err, ErrTaskNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "対象タスクが存在しません。",
			})
			return

		default:
			log.Printf(
				"failed to get task: user_id=%d company_id=%d project_id=%d task_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				taskID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "タスク詳細の取得に失敗しました。",
			})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary タスク編集
// @Description 対象プロジェクトに登録されているタスク情報とブランチ情報を編集します。
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Param taskId path int true "タスクID"
// @Param request body UpdateRequest true "タスク編集情報"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId}/tasks/{taskId} [put]
func (h *Handler) Update(c *gin.Context) {
	projectID, err := strconv.ParseUint(
		c.Param("projectId"),
		10,
		64,
	)
	if err != nil || projectID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "プロジェクトIDが正しくありません。",
		})
		return
	}

	taskID, err := strconv.ParseUint(
		c.Param("taskId"),
		10,
		64,
	)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "タスクIDが正しくありません。",
		})
		return
	}

	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "ログイン情報が無効です。",
		})
		return
	}

	var request UpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "タスク情報を正しく入力してください。",
		})
		return
	}

	response, err := h.service.Update(
		c.Request.Context(),
		projectID,
		taskID,
		authUserID,
		companyID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"task update authorization denied: user_id=%d company_id=%d project_id=%d task_id=%d",
				authUserID,
				companyID,
				projectID,
				taskID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトへアクセスする権限がありません。",
			})
			return

		case errors.Is(err, ErrTaskNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "対象タスクが存在しません。",
			})
			return

		case errors.Is(err, ErrAssigneeNotProjectMember):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "担当メンバーが対象プロジェクトに所属していません。",
			})
			return

		case errors.Is(err, ErrBranchAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Message: "このブランチ名は対象プロジェクトで既に使用されています。",
			})
			return

		case errors.Is(err, ErrInvalidDateRange):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "期限は開始日以降の日付を指定してください。",
			})
			return

		default:
			log.Printf(
				"failed to update task: user_id=%d company_id=%d project_id=%d task_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				taskID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "タスク編集に失敗しました。",
			})
			return
		}
	}

	log.Printf(
		"task updated: task_id=%d user_id=%d company_id=%d project_id=%d",
		taskID,
		authUserID,
		companyID,
		projectID,
	)

	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary タスク削除
// @Description 対象プロジェクトに登録されているタスクを削除します。
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Param taskId path int true "タスクID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId}/tasks/{taskId} [delete]
func (h *Handler) Delete(c *gin.Context) {
	projectID, err := strconv.ParseUint(
		c.Param("projectId"),
		10,
		64,
	)
	if err != nil || projectID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "プロジェクトIDが正しくありません。",
		})
		return
	}

	taskID, err := strconv.ParseUint(
		c.Param("taskId"),
		10,
		64,
	)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "タスクIDが正しくありません。",
		})
		return
	}

	authUserID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if authUserID == 0 || companyID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "ログイン情報が無効です。",
		})
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		projectID,
		taskID,
		authUserID,
		companyID,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"task delete authorization denied: user_id=%d company_id=%d project_id=%d task_id=%d",
				authUserID,
				companyID,
				projectID,
				taskID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトへアクセスする権限がありません。",
			})
			return

		case errors.Is(err, ErrTaskNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "対象タスクが存在しません。",
			})
			return

		default:
			log.Printf(
				"failed to delete task: user_id=%d company_id=%d project_id=%d task_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				taskID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "タスク削除に失敗しました。",
			})
			return
		}
	}

	log.Printf(
		"task deleted: task_id=%d user_id=%d company_id=%d project_id=%d",
		taskID,
		authUserID,
		companyID,
		projectID,
	)

	c.Status(http.StatusNoContent)
}
