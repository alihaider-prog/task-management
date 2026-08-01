package router

import (
	"task-management/internal/handler"
	"task-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerTaskRouts(protected *gin.RouterGroup, taskHandler *handler.TaskHandler, roleMiddleware *middleware.RoleMiddleware) {
	projectTasks := protected.Group("/projects/:id/tasks")
	tasks := protected.Group("/tasks")
	{
		projectTasks.GET("", roleMiddleware.AuthorizeProjectRead(), taskHandler.GetTasks)
		tasks.GET("/:id", roleMiddleware.AuthorizeTaskRead(), taskHandler.GetTaskByID)
		tasks.POST("", roleMiddleware.AuthorizeTaskCreate(), taskHandler.CreateTask)
		tasks.PUT("/:id", roleMiddleware.AuthorizeTaskWrite(), taskHandler.UpdateTask)
		tasks.DELETE("/:id", roleMiddleware.AuthorizeTaskWrite(), taskHandler.DeleteTask)
	}
}
