package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/config" 
	"github.com/prabhatKr-1/lib-man-sys/backend/models"
	"github.com/prabhatKr-1/lib-man-sys/backend/testutils"
	"github.com/stretchr/testify/assert"
)

// 🟢 Test RaiseBookRequest (Issue Request)
func TestRaiseBookRequest_Issue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert a test book with available copies
	book := models.Books{
		ISBN:            123456,
		Title:           "Go Programming",
		Authors:         "John Doe",
		Publisher:       "Tech Press",
		Version:         "1st",
		LibID:           1,
		Total_copies:    5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	// Middleware to mock the reader's login
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(2))     // Fake reader ID
		c.Set("libid", uint(1)) // Fake library ID
		c.Next()
	})

	router.POST("/requests/raise", RaiseBookRequest)

	issuePayload := `{
		"isbn": 123456,
		"requestType": "issue"
	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(issuePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Issue request raised successfully")
}

// 🟢 Test RaiseBookRequest (Return Request)
func TestRaiseBookRequest_Return(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert a test book and issued record
	book := models.Books{
		ISBN:            123456,
		Title:           "Go Programming",
		Authors:         "John Doe",
		Publisher:       "Tech Press",
		Version:         "1st",
		LibID:           1,
		Total_copies:    5,
		Available_copies: 4,
	}
	config.DB.Create(&book)

	issueRegistry := models.IssueRegistry{
		ISBN:      book.ISBN,
		LibID:     book.LibID,
		ReaderID:  2,
		Status:    "issued",
		IssueDate: time.Now(),
	}
	config.DB.Create(&issueRegistry)

	// Middleware to mock the reader's login
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(2)) // Fake reader ID
		c.Set("libid", uint(1)) // Fake library ID
		c.Next()
	})

	router.POST("/requests/raise",  RaiseBookRequest)

	returnPayload := `{
		"isbn": 123456,
		"requestType": "return"
	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(returnPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Return request raised successfully")
}
 
func TestListRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.GET("/requests/list",  ListRequests)

	// Insert sample requests
	request1 := models.RequestEvents{BookID: 123456, ReaderID: 2, LibID: 1, RequestType: "issue", RequestDate: time.Now()}
	request2 := models.RequestEvents{BookID: 123457, ReaderID: 3, LibID: 1, RequestType: "return", RequestDate: time.Now()}
	config.DB.Create(&request1)
	config.DB.Create(&request2)

	req, _ := http.NewRequest("GET", "/requests/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "requests")
}
 
func TestProcessRequest_ApproveIssue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert test book and request
	book := models.Books{
		ISBN:            123456,
		Title:           "Go Programming",
		Authors:         "John Doe",
		Publisher:       "Tech Press",
		Version:         "1st",
		LibID:           1,
		Total_copies:    5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	request := models.RequestEvents{
		BookID:      book.ISBN,
		ReaderID:    2,
		LibID:       book.LibID,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)
 
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1)) // Fake admin ID
		c.Next()
	})

	router.POST("/requests/process",  ProcessRequest)

	approvePayload := `{
		"action": "approve",
		"reqtype": "issue",
		"reqid": 1
	}`

	req, _ := http.NewRequest("POST", "/requests/process", bytes.NewBuffer([]byte(approvePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Issue request approved successfully")
}
 
func TestProcessRequest_Reject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	// Insert test request
	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	// Middleware to mock admin rejection
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1)) // Fake admin ID
		c.Next()
	})

	router.POST("/requests/process", ProcessRequest)

	rejectPayload := `{
		"action": "reject",
		"reqtype": "issue",
		"reqid": 1
	}`

	req, _ := http.NewRequest("POST", "/requests/process", bytes.NewBuffer([]byte(rejectPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Request processed succesfully! Issue req rejected!")
}
