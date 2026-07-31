package handler

import (
	"net/http"
	"strconv"
	"task-management/internal/models"
	"task-management/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	service *service.WorkspaceService
}

func NewWorkspaceHandler(service *service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		service: service,
	}
}

func (h *WorkspaceHandler) Getworkspaces(c *gin.Context) {
	workspaces, err := h.service.GetWorkspaces()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) GetWorkspaceByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workspace id",
		})
		return
	}

	workspace, err := h.service.GetWorkspaceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var workspace models.Workspace
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id type",
		})
		return
	}

	if err := c.ShouldBindBodyWithJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateWorkspace(&workspace, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Workspace Created",
		"workspace id": id,
	})
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workspace id",
		})
		return
	}

	var workspace models.Workspace
	if err := c.ShouldBindBodyWithJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	workspace.ID = id
	if err := h.service.UpdateWorkspace(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workspace updated",
	})
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workspace id",
		})
		return
	}

	if err := h.service.DeleteWorkspace(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workspace deleted",
	})
}
