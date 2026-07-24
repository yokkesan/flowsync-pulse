package task

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	accessibleTasks := api.Group(
		"/tasks",
	)

	accessibleTasks.Use(
		middleware.Authenticate(),
	)

	accessibleTasks.GET(
		"",
		handler.ListAccessible,
	)

	projectTasks := api.Group(
		"/projects/:projectId/tasks",
	)

	projectTasks.Use(
		middleware.Authenticate(),
	)

	projectTasks.GET(
		"",
		handler.List,
	)

	projectTasks.GET(
		"/:taskId",
		handler.Get,
	)

	projectTasks.PUT(
		"/:taskId",
		handler.Update,
	)

	projectTasks.DELETE(
		"/:taskId",
		handler.Delete,
	)

	projectTasks.POST(
		"",
		handler.Create,
	)
}
