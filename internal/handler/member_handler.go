package handler

import (
	"net/http"
	"strconv"
	"task-management/internal/models"
	"task-management/internal/service"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	service *service.MemberService
}

func NewMemberHandler(service *service.MemberService) *MemberHandler {
	return &MemberHandler{service: service}
}

func (h *MemberHandler) AddWorkspaceMember(c *gin.Context) {
	var member models.WorkspaceMembers

	workspaceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workspaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workspace id",
		})
		return
	}
	member.WorkspaceID = workspaceID

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := h.service.AddWorkspaceMember(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Workspace member added",
		"member_id": id,
	})
}

func (h *MemberHandler) ListWorkspaceMembers(c *gin.Context) {
	workspaceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workspaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workspace id",
		})
		return
	}

	members, err := h.service.ListWorkspaceMembers(workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) RemoveWorkspaceMember(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid member id",
		})
		return
	}

	if err := h.service.RemoveWorkspaceMember(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workspace member removed",
	})
}

func (h *MemberHandler) AddProjectMember(c *gin.Context) {
	var member models.ProjectMembers

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}
	member.ProjectID = projectID

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := h.service.AddProjectMember(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Project member added",
		"member_id": id,
	})
}

func (h *MemberHandler) ListProjectMembers(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	members, err := h.service.ListProjectMembers(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) RemoveProjectMember(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid member id",
		})
		return
	}

	if err := h.service.RemoveProjectMember(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project member removed",
	})
}
