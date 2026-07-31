package middleware

import (
	// "errors"
	"net/http"
	"strings"
	"task-management/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("Email", claims.Email)
		c.Next()
	}
}

// func GetUser(c *gin.Context) (*int64, error) {
// 	userIDVal, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "user not authenticated",
// 		})
// 		return nil, errors.New("user not authenticated")
// 	}

// 	userID, ok := userIDVal.(int64)
// 	if !ok {
// 		return nil, errors.New("invalid user id type")
// 	}

// 	return &userID, nil
// }