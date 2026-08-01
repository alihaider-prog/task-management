package router

import (
	"task-management/internal/handler"
	"task-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	taskHandler *handler.TaskHandler,
	workspaceHandler *handler.WorkspaceHandler,
	projectHandler *handler.ProjectHandler,
	memberHandler *handler.MemberHandler,
) *gin.Engine {

	router := gin.Default()
	api := router.Group("api/v1")

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	registerAuthRouts(api, authHandler)
	registerTaskRouts(protected, taskHandler)
	registerWorkspaceRouts(protected, workspaceHandler)
	registerMemberRouts(protected, memberHandler)
	registerProjectRouts(protected, projectHandler)

	return router
}
