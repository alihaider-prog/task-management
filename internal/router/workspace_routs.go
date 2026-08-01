package router

import (
	"task-management/internal/handler"
	"task-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerWorkspaceRouts(protected *gin.RouterGroup, workspaceHandler *handler.WorkspaceHandler, roleMiddleware *middleware.RoleMiddleware) {
	workspace := protected.Group("/workspaces")
	{
		workspace.GET("", workspaceHandler.Getworkspaces)
		workspace.GET("/:id", roleMiddleware.AuthorizeWorkspaceMembers(), workspaceHandler.GetWorkspaceByID)
		workspace.POST("", workspaceHandler.CreateWorkspace)
		workspace.PUT("/:id", roleMiddleware.AuthorizeWorkspaceAdmin(), workspaceHandler.UpdateWorkspace)
		workspace.DELETE("/:id", roleMiddleware.AuthorizeWorkspaceAdmin(), workspaceHandler.DeleteWorkspace)
	}
}
