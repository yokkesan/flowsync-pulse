package user

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	api.POST(
		"/companies/:companyId/users",
		handler.Register,
	)

	users := api.Group("/users")

	users.Use(middleware.Authenticate())

	users.GET(
		"",
		handler.List,
	)
}
