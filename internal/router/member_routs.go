package router

import (
	"task-management/internal/handler"
	"task-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerMemberRouts(protected *gin.RouterGroup, memberHandler *handler.MemberHandler, roleMiddleware *middleware.RoleMiddleware) {
	workspaceMembers := protected.Group("/workspaces/:id/members")
	{
		workspaceMembers.POST("", roleMiddleware.AuthorizeWorkspaceAdmin(), memberHandler.AddWorkspaceMember)
		workspaceMembers.GET("", roleMiddleware.AuthorizeWorkspaceMembers(), memberHandler.ListWorkspaceMembers)
		workspaceMembers.DELETE("/:id", roleMiddleware.AuthorizeWorkspaceAdmin(), memberHandler.RemoveWorkspaceMember)
	}

	projectMembers := protected.Group("/projects/:id/members")
	{
		projectMembers.POST("", roleMiddleware.AuthorizeProjectWorkspaceAdmin(), memberHandler.AddProjectMember)
		projectMembers.GET("", roleMiddleware.AuthorizeProjectRead(), memberHandler.ListProjectMembers)
		projectMembers.DELETE("/:id", roleMiddleware.AuthorizeProjectWorkspaceAdmin(), memberHandler.RemoveProjectMember)
	}
}
