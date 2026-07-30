package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerAuthRouts(api *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
