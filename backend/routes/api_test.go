package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/testutils"
	"github.com/stretchr/testify/assert"
)


func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	
	testutils.SetupTestDB()
	SetupRoutes(r)

	return r
}


func TestAuthRoutes(t *testing.T) {
	router := setupTestRouter()

	
	signupPayload := `{"name":"Test User","email":"test@example.com","password":"testpassword","contactNumber":"1234567890","libraryName":"Test Library"}`
	req, _ := http.NewRequest("POST", "/v1/auth/signup", bytes.NewBuffer([]byte(signupPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	
	loginPayload := `{"email":"test@example.com","password":"testpassword"}`
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer([]byte(loginPayload)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
