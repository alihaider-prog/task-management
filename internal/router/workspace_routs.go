package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerWorkspaceRouts(protected *gin.RouterGroup, workspaceHandler *handler.WorkspaceHandler) {
	workspace := protected.Group("/workspaces")
	{
		workspace.GET("", workspaceHandler.Getworkspaces)
		workspace.GET("/:id", workspaceHandler.GetWorkspaceByID)
		workspace.POST("", workspaceHandler.CreateWorkspace)
		workspace.PUT("/:id", workspaceHandler.UpdateWorkspace)
		workspace.DELETE("/:id", workspaceHandler.DeleteWorkspace)
	}
}
