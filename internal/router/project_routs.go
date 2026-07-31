package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerProjectRouts(protected *gin.RouterGroup, projectHandler *handler.ProjectHandler) {
	projects := protected.Group("/projects")
	{
		projects.GET("", projectHandler.GetProjects)
		projects.GET("/:id", projectHandler.GetProjectByID)
		projects.POST("", projectHandler.CreateProject)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)
	}
}
