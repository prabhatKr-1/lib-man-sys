package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/utils"
)

// Auth Middleware with injectable JWT validator
func AuthMiddleware(jwtValidator utils.JWTValidatorFunc, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		id, LibID, email, role, err := jwtValidator(tokenString)
		if err != nil { // 🔹 Return 401 if token validation fails
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !contains(roles, role) { // 🔹 Return 403 if role authorization fails
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

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
