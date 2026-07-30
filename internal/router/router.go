package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	taskHandler *handler.TaskHandler,
	authHandler *handler.AuthHandler,
) *gin.Engine {

	router := gin.Default()
	api := router.Group("api/v1")

	protected := api.Group("/")

	registerAuthRouts(api, authHandler)
	registerTaskRouts(protected, taskHandler)

	return router
}
