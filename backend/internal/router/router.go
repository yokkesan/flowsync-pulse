package router

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "flowsync-pulse/backend/docs"
	"flowsync-pulse/backend/internal/auth"
	"flowsync-pulse/backend/internal/company"
	"flowsync-pulse/backend/internal/extension"
	"flowsync-pulse/backend/internal/extensiontoken"
	"flowsync-pulse/backend/internal/project"
	"flowsync-pulse/backend/internal/realtime"
	"flowsync-pulse/backend/internal/task"
	"flowsync-pulse/backend/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(
	db *sql.DB,
	realtimeHub *realtime.Hub,
) *gin.Engine {
	engine := gin.Default()

	engine.Use(createCORSConfig())

	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		),
	)

	api := engine.Group("/api")

	if realtimeHub == nil {
		realtimeHub = realtime.NewHub()
	}

	realtimeTicketStore :=
		realtime.NewMemoryConnectionTicketStore()

	registerCompanyRoutes(
		api,
		db,
	)
	registerUserRoutes(
		api,
		db,
	)
	registerAuthRoutes(
		api,
		db,
	)
	registerProjectRoutes(
		api,
		db,
	)
	registerTaskRoutes(
		api,
		db,
	)
	registerExtensionTokenRoutes(
		api,
		db,
	)
	registerExtensionRoutes(
		api,
		db,
		realtimeHub,
	)
	registerRealtimeRoutes(
		api,
		realtimeHub,
		realtimeTicketStore,
	)
	registerHealthRoute(
		api,
		db,
	)

	return engine
}

func createCORSConfig() gin.HandlerFunc {
	frontendOrigin := os.Getenv(
		"FRONTEND_ORIGIN",
	)
	if frontendOrigin == "" {
		frontendOrigin =
			"http://localhost:5174"
	}

	return cors.New(
		cors.Config{
			AllowOrigins: []string{
				frontendOrigin,
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowHeaders: []string{
				"Origin",
				"Content-Type",
				"Accept",
				"Authorization",
			},
			ExposeHeaders: []string{
				"Content-Length",
			},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
}

func registerCompanyRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository :=
		company.NewRepository(db)
	service :=
		company.NewService(repository)
	handler :=
		company.NewHandler(service)

	api.POST(
		"/companies",
		handler.Create,
	)
}

func registerUserRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository :=
		user.NewRepository(db)
	service :=
		user.NewService(repository)
	handler :=
		user.NewHandler(service)

	user.RegisterRoutes(
		api,
		handler,
	)
}

func registerAuthRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository :=
		auth.NewRepository(db)
	service :=
		auth.NewService(repository)
	handler :=
		auth.NewHandler(service)

	auth.RegisterRoutes(
		api,
		handler,
	)
}

func registerProjectRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository :=
		project.NewRepository(db)
	service :=
		project.NewService(repository)
	handler :=
		project.NewHandler(service)

	project.RegisterRoutes(
		api,
		handler,
	)
}

func registerTaskRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository :=
		task.NewRepository(db)
	service :=
		task.NewService(repository)
	handler :=
		task.NewHandler(service)

	task.RegisterRoutes(
		api,
		handler,
	)
}

func registerExtensionTokenRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository :=
		extensiontoken.NewRepository(db)
	service :=
		extensiontoken.NewService(
			repository,
		)
	handler :=
		extensiontoken.NewHandler(
			service,
		)

	extensiontoken.RegisterRoutes(
		api,
		handler,
	)
}

func registerExtensionRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
	realtimeHub *realtime.Hub,
) {
	repository :=
		extension.NewRepository(db)

	notifier :=
		realtime.NewExtensionNotifier(
			realtimeHub,
		)

	service :=
		extension.NewService(
			repository,
			notifier,
		)

	handler :=
		extension.NewHandler(
			service,
		)

	tokenRepository :=
		extensiontoken.NewRepository(db)

	extension.RegisterRoutes(
		api,
		handler,
		tokenRepository,
	)
}

func registerRealtimeRoutes(
	api *gin.RouterGroup,
	hub *realtime.Hub,
	ticketStore realtime.ConnectionTicketStore,
) {
	frontendOrigin := os.Getenv(
		"FRONTEND_ORIGIN",
	)
	if frontendOrigin == "" {
		frontendOrigin =
			"http://localhost:5174"
	}

	handler :=
		realtime.NewHandler(
			hub,
			ticketStore,
			frontendOrigin,
		)

	realtime.RegisterRoutes(
		api,
		handler,
	)
}

func registerHealthRoute(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	api.GET(
		"/health",
		healthCheckHandler(db),
	)
}

// healthCheckHandler godoc
// @Summary ヘルスチェック
// @Description APIサーバーとデータベースの接続状態を返します。
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/health [get]
func healthCheckHandler(
	db *sql.DB,
) gin.HandlerFunc {
	return func(
		c *gin.Context,
	) {
		if err := db.Ping(); err != nil {
			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"status":   "error",
					"database": "disconnected",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"status":   "ok",
				"database": "connected",
			},
		)
	}
}
