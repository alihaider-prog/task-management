package router

import (
	"task-management/internal/handler"
	"task-management/internal/middleware"
	"task-management/internal/models"

	"github.com/gin-gonic/gin"
)

func registerProjectRouts(protected *gin.RouterGroup, projectHandler *handler.ProjectHandler, roleMiddleware *middleware.RoleMiddleware) {
	projects := protected.Group("/projects")
	workspaceProjects := protected.Group("/workspace/:id/projects")
	{
		workspaceProjects.GET("", roleMiddleware.AuthorizeWorkspaceMembers(), projectHandler.GetProjects)
		workspaceProjects.POST("", roleMiddleware.AuthorizeWorkspace(models.RoleAdmin, models.RoleMember), projectHandler.CreateProject)
		projects.GET("/:id", roleMiddleware.AuthorizeProjectRead(), projectHandler.GetProjectByID)
		projects.PUT("/:id", roleMiddleware.AuthorizeProjectWrite(), projectHandler.UpdateProject)
		projects.DELETE("/:id", roleMiddleware.AuthorizeProjectWrite(), projectHandler.DeleteProject)
	}
}
