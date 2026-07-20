package auth

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	authGroup := api.Group("/auth")

	authGroup.POST(
		"/login",
		handler.Login,
	)

	api.GET(
		"/me",
		middleware.Authenticate(),
		handler.CurrentUser,
	)
}
