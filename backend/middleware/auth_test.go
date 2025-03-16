package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/utils"
	"github.com/stretchr/testify/assert"
)

// // 🟢 Mock ValidateJWT function
// func mockValidateJWT(token string) (uint, uint, string, string, error) {
// 	if token == "valid-admin-token" {
// 		return 1, 1, "admin@example.com", "Admin", nil
// 	}
// 	if token == "valid-reader-token" {
// 		return 2, 1, "reader@example.com", "Reader", nil
// 	}
// 	return 0, 0, "", "", assert.AnError // 🔹 Simulate failed token validation
// }

// // 🟢 Setup a test router with middleware and mock JWT validator
// func setupRouterWithMiddleware(jwtValidator utils.JWTValidatorFunc, roles ...string) *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	router.Use(AuthMiddleware(jwtValidator, roles...))
// 	router.GET("/protected", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "Success"})
// 	})
// 	return router
// }

// // 🟢 Test case: Invalid token (Expect `401 Unauthorized`)
// func TestAuthMiddleware_InvalidToken(t *testing.T) {
// 	router := setupRouterWithMiddleware(mockValidateJWT, "Admin")

// 	req, _ := http.NewRequest("GET", "/protected", nil)
// 	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid-token"})
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code) // 🔹 Now correctly expecting 401
// 	assert.Contains(t, w.Body.String(), "Invalid token")
// }

func mockValidateJWT(token string) (uint, uint, string, string, error) {
	if token == "valid-admin-token" {
		return 1, 1, "admin@example.com", "Admin", nil
	}
	if token == "valid-owner-token" {
		return 2, 1, "owner@example.com", "Owner", nil
	}
	if token == "valid-reader-token" {
		return 3, 1, "reader@example.com", "Reader", nil
	}
	return 0, 0, "", "", assert.AnError // Simulate failed token validation
}

// 🟢 Setup a test router with middleware and mock JWT validator
func setupRouterWithMiddleware(jwtValidator utils.JWTValidatorFunc, roles ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware(jwtValidator, roles...))
	router.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("id")
		libID, _ := c.Get("libid")
		email, _ := c.Get("email")
		role, _ := c.Get("role")

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"userID":  userID,
			"libID":   libID,
			"email":   email,
			"role":    role,
		})
	})
	return router
}

// 🟢 Test case: No token in request (Expect `401 Unauthorized`)
func TestAuthMiddleware_NoToken(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "Admin")

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

// 🟢 Test case: Invalid token (Expect `401 Unauthorized`)
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "Admin")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid-token"})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

// 🟢 Test case: Valid token but unauthorized role (Expect `403 Forbidden`)
func TestAuthMiddleware_UnauthorizedRole(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "Admin")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid-reader-token"}) // Reader is not Admin
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Forbidden")
}

// 🟢 Test case: Valid token and authorized role (Expect `200 OK`)
func TestAuthMiddleware_Authorized(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "Admin")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid-admin-token"}) // Admin role
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Success")
}

// 🟢 Test case: Middleware correctly sets user data
func TestAuthMiddleware_SetsUserData(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "Admin")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid-admin-token"})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin@example.com")
	assert.Contains(t, w.Body.String(), `"userID":1`)
	assert.Contains(t, w.Body.String(), `"libID":1`)
	assert.Contains(t, w.Body.String(), `"role":"Admin"`)
}

// 🟢 Test case: Multiple roles access check
func TestAuthMiddleware_MultipleRoles(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "Admin", "Owner")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid-owner-token"}) // Owner role
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) // Owner should have access
	assert.Contains(t, w.Body.String(), "Success")
}

// 🟢 Test case: Role validation is case-sensitive
func TestAuthMiddleware_CaseSensitiveRole(t *testing.T) {
	router := setupRouterWithMiddleware(mockValidateJWT, "admin") // Lowercase "admin"

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid-admin-token"}) // Uppercase "Admin"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code) // "Admin" ≠ "admin"
	assert.Contains(t, w.Body.String(), "Forbidden")
}
