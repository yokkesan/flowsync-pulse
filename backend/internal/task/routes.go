package task

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	tasks := api.Group(
		"/projects/:projectId/tasks",
	)

	tasks.Use(middleware.Authenticate())

	tasks.GET(
		"",
		handler.List,
	)

	tasks.GET(
		"/:taskId",
		handler.Get,
	)

	tasks.POST(
		"",
		handler.Create,
	)
}
