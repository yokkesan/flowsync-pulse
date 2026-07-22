package project

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	projects := api.Group("/projects")

	projects.Use(middleware.Authenticate())

	projects.GET(
		"",
		handler.List,
	)

	projects.GET(
		"/:projectId",
		handler.Get,
	)

	projects.POST(
		"",
		handler.Create,
	)

	projects.PUT(
		"/:projectId",
		handler.Update,
	)

	projects.DELETE(
		"/:projectId",
		handler.Delete,
	)
}
