package project

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
// @Summary プロジェクト登録
// @Description ログインユーザーの会社に新しいプロジェクトを登録します。
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRequest true "プロジェクト情報"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects [post]
func (h *Handler) Create(c *gin.Context) {
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
			Message: "プロジェクト情報を正しく入力してください。",
		})
		return
	}

	response, err := h.service.Create(
		c.Request.Context(),
		companyID,
		authUserID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidProjectMember):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "指定されたメンバーに、会社へ所属していないユーザーが含まれています。",
			})
			return

		case errors.Is(err, ErrProjectSlugAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Message: "このスラッグは既に使用されています。",
			})
			return

		case errors.Is(err, ErrInvalidProjectDateRange):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "終了日は開始日以降の日付を指定してください。",
			})
			return

		default:
			log.Printf(
				"failed to create project: user_id=%d company_id=%d error=%v",
				authUserID,
				companyID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "プロジェクト登録に失敗しました。",
			})
			return
		}
	}

	log.Printf(
		"project created: project_id=%d user_id=%d company_id=%d",
		response.ProjectID,
		authUserID,
		companyID,
	)

	c.JSON(http.StatusCreated, response)
}

// List godoc
// @Summary プロジェクト一覧取得
// @Description ログインユーザーが所属しているプロジェクト一覧を取得します。
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects [get]
func (h *Handler) List(c *gin.Context) {
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

	projects, err := h.service.List(
		c.Request.Context(),
		companyID,
		authUserID,
	)
	if err != nil {
		log.Printf(
			"failed to list projects: user_id=%d company_id=%d error=%v",
			authUserID,
			companyID,
			err,
		)

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "プロジェクト一覧の取得に失敗しました。",
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Projects: projects,
	})
}

// Get godoc
// @Summary プロジェクト詳細取得
// @Description ログインユーザーが所属しているプロジェクトの詳細を取得します。
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId} [get]
func (h *Handler) Get(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
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

	response, err := h.service.Get(
		c.Request.Context(),
		projectID,
		companyID,
		authUserID,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"project detail authorization denied: user_id=%d company_id=%d project_id=%d",
				authUserID,
				companyID,
				projectID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトへアクセスする権限がありません。",
			})
			return

		case errors.Is(err, ErrProjectNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "対象プロジェクトが存在しません。",
			})
			return

		default:
			log.Printf(
				"failed to get project: user_id=%d company_id=%d project_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "プロジェクト詳細の取得に失敗しました。",
			})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary プロジェクト編集
// @Description ログインユーザーが所属しているプロジェクトを編集します。
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Param request body UpdateRequest true "プロジェクト編集情報"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId} [put]
func (h *Handler) Update(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
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

	var request UpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "プロジェクト情報を正しく入力してください。",
		})
		return
	}

	response, err := h.service.Update(
		c.Request.Context(),
		projectID,
		companyID,
		authUserID,
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"project update authorization denied: user_id=%d company_id=%d project_id=%d",
				authUserID,
				companyID,
				projectID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトを編集する権限がありません。",
			})
			return

		case errors.Is(err, ErrProjectNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "対象プロジェクトが存在しません。",
			})
			return

		case errors.Is(err, ErrInvalidProjectMember):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "指定されたメンバーに、会社へ所属していないユーザーが含まれています。",
			})
			return

		case errors.Is(err, ErrProjectSlugAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Message: "このスラッグは既に使用されています。",
			})
			return

		case errors.Is(err, ErrInvalidProjectDateRange):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "終了日は開始日以降の日付を指定してください。",
			})
			return

		default:
			log.Printf(
				"failed to update project: user_id=%d company_id=%d project_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "プロジェクト編集に失敗しました。",
			})
			return
		}
	}

	log.Printf(
		"project updated: project_id=%d user_id=%d company_id=%d",
		projectID,
		authUserID,
		companyID,
	)

	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary プロジェクト削除
// @Description ログインユーザーが所属しているプロジェクトを削除します。
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "プロジェクトID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/{projectId} [delete]
func (h *Handler) Delete(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
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

	err = h.service.Delete(
		c.Request.Context(),
		projectID,
		companyID,
		authUserID,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrProjectAccessDenied):
			log.Printf(
				"project delete authorization denied: user_id=%d company_id=%d project_id=%d",
				authUserID,
				companyID,
				projectID,
			)

			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "このプロジェクトを削除する権限がありません。",
			})
			return

		case errors.Is(err, ErrProjectNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "対象プロジェクトが存在しません。",
			})
			return

		default:
			log.Printf(
				"failed to delete project: user_id=%d company_id=%d project_id=%d error=%v",
				authUserID,
				companyID,
				projectID,
				err,
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "プロジェクト削除に失敗しました。",
			})
			return
		}
	}

	log.Printf(
		"project deleted: project_id=%d user_id=%d company_id=%d",
		projectID,
		authUserID,
		companyID,
	)

	c.Status(http.StatusNoContent)
}

func parseProjectID(c *gin.Context) (uint64, error) {
	projectID, err := strconv.ParseUint(
		c.Param("projectId"),
		10,
		64,
	)
	if err != nil || projectID == 0 {
		return 0, errors.New("invalid project id")
	}

	return projectID, nil
}
