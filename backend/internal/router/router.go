package router

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "flowsync-pulse/backend/docs"
	"flowsync-pulse/backend/internal/auth"
	"flowsync-pulse/backend/internal/company"
	"flowsync-pulse/backend/internal/task"
	"flowsync-pulse/backend/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(db *sql.DB) *gin.Engine {
	engine := gin.Default()

	engine.Use(createCORSConfig())

	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	api := engine.Group("/api")

	registerCompanyRoutes(api, db)
	registerUserRoutes(api, db)
	registerAuthRoutes(api, db)
	registerTaskRoutes(api, db)
	registerHealthRoute(api, db)

	return engine
}

func createCORSConfig() gin.HandlerFunc {
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5174"
	}

	return cors.New(cors.Config{
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
	})
}

func registerCompanyRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository := company.NewRepository(db)
	service := company.NewService(repository)
	handler := company.NewHandler(service)

	api.POST(
		"/companies",
		handler.Create,
	)
}

func registerUserRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository := user.NewRepository(db)
	service := user.NewService(repository)
	handler := user.NewHandler(service)

	api.POST(
		"/companies/:companyId/users",
		handler.Register,
	)
}

func registerAuthRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository := auth.NewRepository(db)
	service := auth.NewService(repository)
	handler := auth.NewHandler(service)

	auth.RegisterRoutes(api, handler)
}

func registerTaskRoutes(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	repository := task.NewRepository(db)
	service := task.NewService(repository)
	handler := task.NewHandler(service)

	task.RegisterRoutes(api, handler)
}

func registerHealthRoute(
	api *gin.RouterGroup,
	db *sql.DB,
) {
	api.GET("/health", healthCheckHandler(db))
}

// healthCheckHandler godoc
// @Summary ヘルスチェック
// @Description APIサーバーとデータベースの接続状態を返します。
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/health [get]
func healthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"database": "disconnected",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	}
}
