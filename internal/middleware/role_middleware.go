package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"task-management/internal/models"
	"task-management/internal/repository"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
	workspaceRepo *repository.WorkspaceRepository
	projectRepo   *repository.ProjectRepository
	taskRepo      *repository.TaskRepository
	memberRepo    *repository.MemberRepository
}

func NewRoleMiddleware(
	workspaceRepo *repository.WorkspaceRepository,
	projectRepo *repository.ProjectRepository,
	taskRepo *repository.TaskRepository,
	memberRepo *repository.MemberRepository,
) *RoleMiddleware {
	return &RoleMiddleware{
		workspaceRepo: workspaceRepo,
		projectRepo:   projectRepo,
		taskRepo:      taskRepo,
		memberRepo:    memberRepo,
	}
}

func (m *RoleMiddleware) AuthorizeProjectRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || projectID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
			return
		}
		m.authorizeProjectResource(c, projectID, models.RoleMember, models.RoleViewer)
	}
}

func (m *RoleMiddleware) AuthorizeProjectWrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || projectID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
			return
		}
		m.authorizeProjectResource(c, projectID, models.RoleMember)
	}
}

func (m *RoleMiddleware) AuthorizeTaskRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || taskID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
			return
		}

		projectID, err := m.taskRepo.GetProjectIDByTaskID(taskID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		m.authorizeProjectResource(c, *projectID, models.RoleMember, models.RoleViewer)
	}
}

func (m *RoleMiddleware) AuthorizeTaskWrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || taskID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
			return
		}

		projectID, err := m.taskRepo.GetProjectIDByTaskID(taskID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		m.authorizeProjectResource(c, *projectID, models.RoleMember)
	}
}

func (m *RoleMiddleware) AuthorizeTaskCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			ProjectID int64 `json:"project_id"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		if payload.ProjectID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
			return
		}

		m.authorizeProjectResource(c, payload.ProjectID, models.RoleMember)
	}
}

func (m *RoleMiddleware) AuthorizeTaskAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || taskID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
			return
		}

		projectID, err := m.taskRepo.GetProjectIDByTaskID(taskID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		m.authorizeProjectResource(c, *projectID, models.RoleAdmin)
	}
}

func (m *RoleMiddleware) AuthorizeWorkspaceMembers() gin.HandlerFunc {
	return m.AuthorizeWorkspace(models.RoleAdmin, models.RoleMember, models.RoleViewer)
}

func (m *RoleMiddleware) AuthorizeWorkspaceAdmin() gin.HandlerFunc {
	return m.AuthorizeWorkspace(models.RoleAdmin)
}

func (m *RoleMiddleware) AuthorizeProjectWorkspaceAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || projectID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
			return
		}
		m.authorizeProjectResource(c, projectID, models.RoleAdmin)
	}
}

func (m *RoleMiddleware) CheckProjectAccess(projectID, userID int64, allowed ...models.WorkspaceMemberRole) error {
	workspaceID, err := m.projectRepo.GetWorkspaceIDByProjectID(projectID)
	if err != nil {
		return err
	}

	role, err := m.workspaceRepo.GetUserRole(*workspaceID, userID)
	if err != nil {
		return err
	}

	if role == string(models.RoleAdmin) {
		return nil
	}

	if !m.roleAllowed(role, allowed) {
		return errors.New("permission denied")
	}

	isMember, err := m.memberRepo.IsProjectMember(projectID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("permission denied")
	}

	return nil
}

func (m *RoleMiddleware) authorizeProjectResource(c *gin.Context, projectID int64, allowed ...models.WorkspaceMemberRole) {
	userID := c.GetInt64("userID")
	workspaceID, err := m.projectRepo.GetWorkspaceIDByProjectID(projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	role, err := m.workspaceRepo.GetUserRole(*workspaceID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
		return
	}

	if role == string(models.RoleAdmin) {
		c.Next()
		return
	}

	if !m.roleAllowed(role, allowed) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
		return
	}

	isMember, err := m.memberRepo.IsProjectMember(projectID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
		return
	}

	if !isMember {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
		return
	}

	c.Next()
}

func (m *RoleMiddleware) AuthorizeWorkspace(allowed ...models.WorkspaceMemberRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("userID")
		workspaceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || workspaceID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid workspace id",
			})
			return
		}

		role, err := m.workspaceRepo.GetUserRole(workspaceID, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		if m.roleAllowed(role, allowed) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
	}
}

func (m *RoleMiddleware) roleAllowed(role string, allowed []models.WorkspaceMemberRole) bool {
	for _, allowedRole := range allowed {
		if role == string(allowedRole) {
			return true
		}
	}
	return false
}