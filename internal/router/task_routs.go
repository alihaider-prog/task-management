package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerTaskRouts(protected *gin.RouterGroup, taskHandler *handler.TaskHandler) {
	tasks := protected.Group("/tasks")
	{
		tasks.GET("", taskHandler.GetTasks)
		tasks.GET("/:id", taskHandler.GetTaskByID)
		tasks.POST("", taskHandler.CreateTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}
}
