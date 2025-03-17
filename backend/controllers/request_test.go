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

type CustomResponseRecorder struct {
	*httptest.ResponseRecorder
}

func (c *CustomResponseRecorder) CloseNotify() <-chan bool {
	return make(chan bool) 
}


func TestRaiseBookRequest_Issue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	
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

	
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(2))     
		c.Set("libid", uint(1)) 
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


func TestRaiseBookRequest_Return(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	
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

	
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(2)) 
		c.Set("libid", uint(1)) 
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
		c.Set("id", uint(1)) 
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

	
	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	
	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1)) 
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
func TestHandleReturnRequest_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    testutils.SetupTestDB()

    
    book := models.Books{
        ISBN:             123456,
        Title:            "Go Programming",
        Authors:          "John Doe",
        Publisher:        "Tech Press",
        Version:          "1st",
        LibID:            1,
        Total_copies:     5,
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

    
    returnRequest := models.RequestEvents{
        BookID:      book.ISBN,
        ReaderID:    2,
        LibID:       book.LibID,
        RequestType: "return",
        RequestDate: time.Now(),
    }
    config.DB.Create(&returnRequest)

    
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    c.Set("id", uint(1)) 
    c.Params = []gin.Param{{Key: "reqId", Value: "1"}}

    
    handleReturnRequest(c, 1, returnRequest.ReqID)

    
    var updatedIssueRegistry models.IssueRegistry
    err := config.DB.Where("isbn = ? AND reader_id = ?", book.ISBN, 2).First(&updatedIssueRegistry).Error
    assert.Nil(t, err)
    assert.Equal(t, "returned", updatedIssueRegistry.Status)

    var updatedBook models.Books
    err = config.DB.Where("isbn = ?", book.ISBN).First(&updatedBook).Error
    assert.Nil(t, err)
    assert.Equal(t, uint(5), updatedBook.Available_copies) 

    var deletedRequest models.RequestEvents
    err = config.DB.Where("req_id = ?", returnRequest.ReqID).First(&deletedRequest).Error
    assert.NotNil(t, err) 
}

func TestHandleReturnRequest_RequestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	
	w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
	c, _ := gin.CreateTestContext(w)

	
	c.Set("id", uint(1)) 
	c.Params = []gin.Param{{Key: "reqId", Value: "999"}} 

	
	handleReturnRequest(c, 1, 999)

	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}


func TestHandleReturnRequest_IssueRegistryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	
	returnRequest := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "return",
		RequestDate: time.Now(),
	}
	config.DB.Create(&returnRequest)

	
	w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
	c, _ := gin.CreateTestContext(w)
	c.Set("id", uint(1)) 

	
	handleReturnRequest(c, 1, returnRequest.ReqID)

	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
