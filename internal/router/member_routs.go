package router

import (
	"task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerMemberRouts(protected *gin.RouterGroup, memberHandler *handler.MemberHandler) {
	workspaceMembers := protected.Group("/workspaces/:id/members")
	{
		workspaceMembers.POST("", memberHandler.AddWorkspaceMember)
		workspaceMembers.GET("", memberHandler.ListWorkspaceMembers)
		workspaceMembers.DELETE("/:id", memberHandler.RemoveWorkspaceMember)
	}

	projectMembers := protected.Group("/projects/:id/members")
	{
		projectMembers.POST("", memberHandler.AddProjectMember)
		projectMembers.GET("", memberHandler.ListProjectMembers)
		projectMembers.DELETE("/:id", memberHandler.RemoveProjectMember)
	}
}
