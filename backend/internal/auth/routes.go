package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	authGroup := api.Group("/auth")

	authGroup.POST(
		"/login",
		handler.Login,
	)
}
