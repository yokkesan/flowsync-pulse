package extensiontoken

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	extensionTokens := api.Group(
		"/extension-tokens",
	)

	extensionTokens.Use(
		middleware.Authenticate(),
	)

	extensionTokens.GET(
		"",
		handler.List,
	)

	extensionTokens.POST(
		"",
		handler.Create,
	)

	extensionTokens.DELETE(
		"/:tokenId",
		handler.Revoke,
	)
}
