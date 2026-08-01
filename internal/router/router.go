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
	roleMiddleware *middleware.RoleMiddleware,
) *gin.Engine {

	router := gin.Default()
	api := router.Group("api/v1")

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	registerAuthRouts(api, authHandler)
	registerTaskRouts(protected, taskHandler, roleMiddleware)
	registerWorkspaceRouts(protected, workspaceHandler, roleMiddleware)
	registerMemberRouts(protected, memberHandler, roleMiddleware)
	registerProjectRouts(protected, projectHandler, roleMiddleware)

	return router
}
