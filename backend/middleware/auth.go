package middleware

import (
    "github.com/prabhatKr-1/lib-man-sys/backend/utils"
    "net/http" 

    "github.com/gin-gonic/gin"
)
func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		id, LibID, email, role, err := utils.ValidateJWT(tokenString)
		if err != nil || !contains(roles, role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		id = uint(id)
		LibID = uint(LibID)
		c.Set("id", id)
		c.Set("libid", LibID)
		c.Set("email", email)
		c.Set("role", role)
		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
