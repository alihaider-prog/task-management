package middleware

import (
	"net/http"
	"strconv"
	"task-management/internal/repository"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
	workspaceRepo *repository.WorkspaceRepository
}

func NewRoleMiddleware(repo *repository.WorkspaceRepository) *RoleMiddleware {
	return &RoleMiddleware{
		workspaceRepo: repo,
	}
}

// func (m *RoleMiddleware) AdminOnly() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.GetInt64("userID")
// 		workspaceID, err := strconv.ParseInt(
// 			c.Param("workspaceID"),
// 			10,
// 			64,
// 		)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 				"error": "invalid workspace is",
// 			})
// 			return
// 		}

// 		role, err := m.workspaceRepo.GetUserRole(workspaceID, userID)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
// 				"error": "access denied",
// 			})
// 			return
// 		}

// 		if role != "Admin" {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
// 				"error": "this feature for admin only",
// 			})
// 			return
// 		}

// 		c.Next()
// 	}
// }

func (m *RoleMiddleware) Authorize(roles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := c.GetInt64("userID")

		workspaceID, err := strconv.ParseInt(c.Param("workspaceID"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid workspace is",
			})
			return
		}

		role, err := m.workspaceRepo.GetUserRole(workspaceID, userID)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "access denied",
			})
			return
		}

		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error": "permission denied",
		})
	}
}
