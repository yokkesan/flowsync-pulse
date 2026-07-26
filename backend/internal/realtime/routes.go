package realtime

import (
	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	handler *Handler,
) {
	realtimeGroup := api.Group(
		"/realtime",
	)

	realtimeGroup.POST(
		"/tickets",
		middleware.Authenticate(),
		handler.CreateConnectionTicket,
	)

	realtimeGroup.GET(
		"/ws",
		handler.Connect,
	)
}
