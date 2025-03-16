package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/config"
	"github.com/prabhatKr-1/lib-man-sys/backend/models"
	"github.com/prabhatKr-1/lib-man-sys/backend/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAddBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB() // Ensure DB is set up before running test

	// Ensure DB is not nil before proceeding
	if config.DB == nil {
		t.Fatalf("❌ config.DB is nil! Database not initialized.")
	}

	router := gin.Default()

	// Manually add middleware values required by the controller
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) // Fake library ID
		c.Next()
	})

	router.POST("/books/add", AddBook)

	// Insert a test library
	lib := models.Library{Name: "Test Library"}
	result := config.DB.Create(&lib)

	// Ensure the library creation was successful
	if result.Error != nil {
		t.Fatalf("❌ Failed to insert test library: %v", result.Error)
	}

	// Create a test book payload
	bookPayload := `{
		"title": "Go Programming",
		"authors": "John Doe",
		"publisher": "Tech Press",
		"version": "1st",
		"isbn": 123456,
		"total_copies": 5
	}`

	req, _ := http.NewRequest("POST", "/books/add", bytes.NewBuffer([]byte(bookPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify response status code
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK response")

	// Verify the book exists in the database
	var book models.Books
	err := config.DB.First(&book, "isbn = ?", 123456).Error
	assert.Nil(t, err, "❌ Book should exist in the database")
	assert.Equal(t, "Go Programming", book.Title)
}

func TestSearchBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert a test book
	book := models.Books{
		ISBN:      123456,
		Title:     "Go Programming",
		Authors:   "John Doe",
		Publisher: "Tech Press",
		Version:   "1st",
		LibID:     1,
		Total_copies: 5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	// Middleware to mock authentication
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) // Fake library ID
		c.Next()
	})

	router.POST("/books/search", SearchBook)

	searchPayload := `{
		"title": "Go Programming"
	}`

	req, _ := http.NewRequest("POST", "/books/search", bytes.NewBuffer([]byte(searchPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Go Programming")
}
 
func TestUpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert a test book
	book := models.Books{
		ISBN:      123456,
		Title:     "Go Programming",
		Authors:   "John Doe",
		Publisher: "Tech Press",
		Version:   "1st",
		LibID:     1,
		Total_copies: 5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	// Middleware to mock authentication
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) // Fake library ID
		c.Next()
	})

	router.PUT("/books/update", UpdateBook)

	updatePayload := `{
		"isbn": 123456,
		"title": "Advanced Go Programming",
		"authors": "Jane Doe",
		"publisher": "Tech Press 2nd Edition",
		"version": "2nd",
		"total_copies": 10
	}`

	req, _ := http.NewRequest("PUT", "/books/update", bytes.NewBuffer([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Book updated successfully")
 
	var updatedBook models.Books
	config.DB.Where("isbn = ?", 123456).First(&updatedBook)
	assert.Equal(t, "Advanced Go Programming", updatedBook.Title)
	assert.Equal(t, "Jane Doe", updatedBook.Authors)
	assert.Equal(t, "Tech Press 2nd Edition", updatedBook.Publisher)
	assert.Equal(t, "2nd", updatedBook.Version)
	assert.Equal(t, uint(10), updatedBook.Total_copies)
}
 
func TestDeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
 
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Role:     "Admin",
		LibID:    1,
	}
	config.DB.Create(&admin)
 
	book := models.Books{
		ISBN:      123456,
		Title:     "Go Programming",
		Authors:   "John Doe",
		Publisher: "Tech Press",
		Version:   "1st",
		LibID:     1,
		Total_copies: 5,
		Available_copies: 5,
	}
	config.DB.Create(&book)
 
	router.Use(func(c *gin.Context) {
		c.Set("email", "admin@example.com")  
		c.Next()
	})

	router.DELETE("/books/delete/123456", DeleteBook)

	req, _ := http.NewRequest("DELETE", "/books/delete/123456", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Book deleted successfully")
 
	var deletedBook models.Books
	err := config.DB.Where("isbn = ?", 123456).First(&deletedBook).Error
	assert.NotNil(t, err)   
}
