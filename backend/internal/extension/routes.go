package extension

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
	authenticator middleware.ExtensionTokenAuthenticator,
) {
	extensionGroup := api.Group(
		"/extension",
	)

	extensionGroup.Use(
		middleware.AuthenticateExtension(
			authenticator,
		),
	)

	extensionGroup.POST(
		"/work-context",
		handler.WorkContext,
	)

	extensionGroup.POST(
		"/heartbeat",
		handler.Heartbeat,
	)

	extensionGroup.POST(
		"/disconnect",
		handler.Disconnect,
	)
}
