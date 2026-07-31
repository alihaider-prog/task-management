package router

import (
	"task-management/internal/handler"
	"task-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	taskHandler *handler.TaskHandler,
	authHandler *handler.AuthHandler,
) *gin.Engine {

	router := gin.Default()
	api := router.Group("api/v1")

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	registerAuthRouts(api, authHandler)
	registerTaskRouts(protected, taskHandler)

	return router
}
