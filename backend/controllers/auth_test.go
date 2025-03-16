package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/config"

	// "github.com/prabhatKr-1/lib-man-sys/backend/controllers"
	"github.com/prabhatKr-1/lib-man-sys/backend/models"
	"github.com/prabhatKr-1/lib-man-sys/backend/testutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.POST("/auth/signup", Signup)

	signupPayload := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"contactNumber": "1234567890",
		"libraryName": "Unique Library Name"
	}`

	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer([]byte(signupPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Ensure signup response is 200 OK
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK response")
}
 
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/login", Login)

	// Insert a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "Owner",
		LibID:    1,
	}
	config.DB.Create(&user)

	loginPayload := `{"email": "test@example.com", "password": "password123"}`

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(loginPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateAdminUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Manually inject an Owner into the DB
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("ownerpassword"), bcrypt.DefaultCost)
	owner := models.User{
		Name:     "Owner User",
		Email:    "owner@example.com",
		Password: string(hashedPassword),
		Role:     "Owner",
		LibID:    1,
	}
	config.DB.Create(&owner)

	// Middleware to mock the owner’s login
	router.Use(func(c *gin.Context) {
		c.Set("email", "owner@example.com")
		c.Next()
	})

	router.POST("/admin/create", CreateAdminUser)

	adminPayload := `{
		"name": "Admin User",
		"email": "admin@example.com",
		"password": "adminpassword",
		"contactNumber": "9876543210"
	}`

	req, _ := http.NewRequest("POST", "/admin/create", bytes.NewBuffer([]byte(adminPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Admin user created successfully")

	// Verify admin exists in DB
	var admin models.User
	err := config.DB.Where("email = ?", "admin@example.com").First(&admin).Error
	assert.Nil(t, err)
	assert.Equal(t, "Admin", admin.Role)
}
 
func TestCreateReaderUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert an Admin User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Role:     "Admin",
		LibID:    1,
	}
	config.DB.Create(&admin)

	// Middleware to mock the admin’s login
	router.Use(func(c *gin.Context) {
		c.Set("email", "admin@example.com")
		c.Next()
	})

	router.POST("/reader/create", CreateReaderUser)

	readerPayload := `{
		"name": "Reader User",
		"email": "reader@example.com",
		"password": "readerpassword",
		"contactNumber": "1234567890"
	}`

	req, _ := http.NewRequest("POST", "/reader/create", bytes.NewBuffer([]byte(readerPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Reader user created successfully")

	// Verify reader exists in DB
	var reader models.User
	err := config.DB.Where("email = ?", "reader@example.com").First(&reader).Error
	assert.Nil(t, err)
	assert.Equal(t, "Reader", reader.Role)
}
 
func TestUpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()   
	router := gin.Default()
 
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "Test User",
		Email:    "user@example.com",
		Password: string(hashedPassword),  
		Role:     "Reader",
		LibID:    1,
	}
 
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
 
	router.Use(func(c *gin.Context) {
		c.Set("email", "user@example.com")
		c.Next()
	})

	router.POST("/user/update-password", UpdatePassword)
 
	passwordPayload := `{
		"oldPassword": "oldpassword",
		"newPassword": "newsecurepassword"
	}`
	req, _ := http.NewRequest("POST", "/user/update-password", bytes.NewBuffer([]byte(passwordPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
 
	assert.Equal(t, http.StatusOK, w.Code, "Password update request should succeed")
	assert.Contains(t, w.Body.String(), "Password updated successfully!")
 
	var updatedUser models.User
	config.DB.Where("email = ?", "user@example.com").First(&updatedUser)
 
	err := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newsecurepassword"))
	assert.Nil(t, err, "Password should be updated correctly")
 
	incorrectPayload := `{
		"oldPassword": "wrongpassword",
		"newPassword": "anothernewpassword"
	}`
	req, _ = http.NewRequest("POST", "/user/update-password", bytes.NewBuffer([]byte(incorrectPayload)))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should fail with incorrect old password")
	assert.Contains(t, w.Body.String(), "Wrong Password")
}
 
func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.GET("/logout", Logout)

	req, _ := http.NewRequest("GET", "/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Logged out successfully")
}
