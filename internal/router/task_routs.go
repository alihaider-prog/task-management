package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerTaskRouts(protected *gin.RouterGroup, taskHandler *handler.TaskHandler) {
	projectTasks := protected.Group("/projects/:id/tasks")
	tasks := protected.Group("/tasks")
	{
		projectTasks.GET("", taskHandler.GetTasks)
		tasks.GET("/:id", taskHandler.GetTaskByID)
		tasks.POST("", taskHandler.CreateTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}
}
