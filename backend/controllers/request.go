package controllers

import (
  "errors"
  "github.com/prabhatKr-1/lib-man-sys/backend/config"
  "github.com/prabhatKr-1/lib-man-sys/backend/models"
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
)

// RAISING ISSUE/RETURN REQUESTS
func RaiseBookRequest(c *gin.Context) {
	var input struct {
		ISBN        uint   `binding:"required"`
		RequestType string `binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !(input.RequestType == "issue" || input.RequestType == "return") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request types can be issue OR return only!",
		})
		return
	}

	var book models.Books
	if err := config.DB.Where("isbn = ?", input.ISBN).First(&book).Error; err != nil || book.Available_copies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book unavailable"})
		return
	}

	floatId, _ := c.Get("id")
	id := floatId.(uint)

	libId, _ := c.Get("libid")

	// CHECKING IF THE REQUEST EXISTS ALREADY
	var bookReq models.RequestEvents
	if input.RequestType == "issue" {
		if err := config.DB.Where("libid = ? AND reader_Id= ? AND book_ID = ? ", libId, id, input.ISBN).Take(&bookReq).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Duplicate request!",
			})
			return
		}

		bookReq.BookID = input.ISBN
		bookReq.ReaderID = id
		bookReq.LibID = libId.(uint)
		bookReq.RequestType = input.RequestType
		bookReq.RequestDate = time.Now()

		// CREATING BOOK REQUEST
		if err := config.DB.Create(&bookReq).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Issue request raised successfully"})
		return
	}

	// BOOK RETURN REQUEST
	var issueReg models.IssueRegistry

	if err := config.DB.Where("isbn = ?", input.ISBN).Where("reader_ID = ?", id).Where("status = ?", "issued").First(&issueReg).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No issue request exists corresponding to this return request",
		})
		return
	}

	bookReq.BookID = input.ISBN
	bookReq.ReaderID = id
	bookReq.LibID = libId.(uint)
	bookReq.RequestType = input.RequestType
	bookReq.RequestDate = time.Now()

	if err := config.DB.Create(&bookReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Return request raised successfully",
	})

}

func ListRequests(c *gin.Context) {
	var requests []models.RequestEvents

	if err := config.DB.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
	})
}

// HANDLING RETURN REQUEST
func handleReturnRequest(c *gin.Context, returnapproverID, reqId uint) {
	var req models.RequestEvents
	// FETCHING REQUEST DETAILS
	if err := config.DB.First(&req, reqId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// UPDATING THE ISSUE REGISTRY
	var retRegistry models.IssueRegistry
	if err := config.DB.Where("isbn = ? AND reader_id = ? AND lib_id = ?", req.BookID, req.ReaderID, req.LibID).First(&retRegistry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	retRegistry.ReturnApproverID = &returnapproverID
	now := time.Now()
	retRegistry.ReturnDate = &now
	retRegistry.Status = "returned"

	if err := config.DB.Save(&retRegistry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.Delete(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	// UPDATING BOOK COUNT
	var book models.Books
	if err := config.DB.Where("isbn = ? AND lib_id = ?", retRegistry.ISBN, retRegistry.LibID).First(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	book.Available_copies += 1
	if err := config.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while updating book available copies",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Return request processed successfully!",
	})
}

func ProcessRequest(c *gin.Context) {
	var input struct {
		Action  string `binding:"required" json:"action"`
		Reqtype string `binding:"required" json:"reqtype"`
		ReqID   uint   `binding:"required" json:"reqid"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// GETTING ADMIN ID from JWT
	id, _ := c.Get("id")
	ApproverID := id.(uint)

	// HANDLING BOOK RETURNS
	if input.Reqtype == "return" {
		handleReturnRequest(c, ApproverID, input.ReqID)
		return
	}

	if !(input.Action == "approve" || input.Action == "reject" || input.Action == "Approve" || input.Action == "Reject") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Allowed actions are approve or reject",
		})
		return
	}

	if input.Action == "reject" {
		// REMOVE ENTRY FROM REQUESTS TABLE
		if err := config.DB.Model(&models.RequestEvents{}).Delete(input.ReqID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed succesfully! Issue req rejected!",
		})

		return
	}

	txErr := config.DB.Transaction(func(tx *gorm.DB) error {
		var req models.RequestEvents
		if tx.Where("req_id = ?", input.ReqID).First(&req).Error != nil {
			return errors.New("request not found")
		}

		var book models.Books
		if tx.Where("isbn = ?", req.BookID).First(&book).Error != nil || book.Available_copies <= 0 {
			return errors.New("no copies available")
		}

		if tx.Delete(&req).Error != nil {
			return errors.New("request not found")
		}

		// ISSUE REQUEST APPROVED, ADD THIS RECORD INTO ISSUE REGISTRY
		now := time.Now()
		issueReg := models.IssueRegistry{
			ISBN:               book.ISBN,
			LibID:              req.LibID,
			ReaderID:           req.ReaderID,
			IssueApproverID:    ApproverID,
			Status:             "issued",
			IssueDate:          now,
			ExpectedReturnDate: now.AddDate(0, 0, 14),
		}
		if err := tx.Create(&issueReg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		book.Available_copies -= 1
		tx.Save(&book)

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": txErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue request approved successfully"})
}
