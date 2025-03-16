package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prabhatKr-1/lib-man-sys/backend/config"
	"github.com/prabhatKr-1/lib-man-sys/backend/models"

	"github.com/gin-gonic/gin"
)

func AddBook(c *gin.Context) {
	var book models.Books
	libId, _ := c.Get("libid")

	var lib models.Library
	if err := config.DB.First(&lib, libId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	book.LibID = libId.(uint)
	book.Library = lib
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Books
	if err := config.DB.Where("isbn = ? AND lib_id = ?", book.ISBN, book.LibID).First(&existing).Error; err == nil {
		existing.Total_copies += book.Total_copies
		existing.Available_copies += book.Total_copies
		config.DB.Save(&existing)
		c.JSON(http.StatusOK, gin.H{"message": "Book copies updated"})
		return
	}

	book.Available_copies = book.Total_copies
	config.DB.Create(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})
}

// SEARCHING A BOOK
func SearchBook(c *gin.Context) {
	var input struct {
		Title   string `json:"title"`
		ISBN    uint   `json:"isbn"`
		Authors string `json:"authors"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title == "" && input.ISBN == 0 && input.Authors == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search field is required"})
		return
	}

	libId, _ := c.Get("libid")

	var book models.Books
	query := config.DB.Where("lib_id = ?", libId)
	// if input.Title != "" {
	// 	query = query.Where("title ILIKE ?", "%"+input.Title+"%")
	// }
	if input.Title != "" {
		if config.DB.Dialector.Name() == "sqlite" {
			query = query.Where("title LIKE ?", "%"+input.Title+"%")
		} else {
			query = query.Where("title ILIKE ?", "%"+input.Title+"%")
		}
	}
	
	
	if input.ISBN != 0 {
		query = query.Where("isbn = ?", input.ISBN)
	}
	if input.Authors != "" {
		query = query.Where("authors ILIKE ?", "%"+input.Authors+"%")
	}

	// EXPECTED DATE A BOOK BECOMES AVAILABLE
	if err := query.First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var bookDetails struct {
		ISBN, Available_copies             uint
		Title, Authors, Publisher, Version string
	}
	bookDetails.ISBN = book.ISBN
	bookDetails.Title = book.Title
	bookDetails.Authors = book.Authors
	bookDetails.Publisher = book.Publisher
	bookDetails.Version = book.Version
	bookDetails.Available_copies = book.Available_copies

	if book.Available_copies == 0 {
		var nextAvailable time.Time
		if e := config.DB.Raw(`
			SELECT expected_return_date FROM issue_registries 
			WHERE isbn = ? AND lib_id = ? AND issue_status ILIKE 'approve' 
			ORDER BY expected_return_date ASC LIMIT 1`, book.ISBN, libId).Scan(&nextAvailable).Error; e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": e.Error(),
			})
			return
		}

		if !nextAvailable.IsZero() {
			c.JSON(http.StatusOK, gin.H{
				"Message":                    "Book is currently unavailable",
				"expected_availability_date": nextAvailable.Format("2006-01-02"),
				"book":                       bookDetails,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"book": bookDetails})
}


// UPDATING THE DETAILS OF A BOOK
func UpdateBook(c *gin.Context) {
	var input struct {
		ISBN             uint   `json:"isbn" binding:"required"`
		LibID            uint   `json:"lib_id"`
		Title            string `json:"title"`
		Authors          string `json:"authors"`
		Publisher        string `json:"publisher"`
		Version          string `json:"version"`
		TotalCopies      uint   `json:"total_copies"`
		Available_copies uint   `json:"available_copies"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	libId, _ := c.Get("libid")
	input.LibID = libId.(uint)

	// SEARCHING BOOK
	var book models.Books
	if err := config.DB.Where("isbn = ? AND lib_id = ?", input.ISBN, input.LibID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	flag := true

	// UPDATIONS
	if input.Title != "" {
		book.Title = input.Title
		flag = false
	}
	if input.Authors != "" {
		book.Authors = input.Authors
		flag = false
	}
	if input.Publisher != "" {
		book.Publisher = input.Publisher
		flag = false
	}
	if input.Version != "" {
		book.Version = input.Version
		flag = false
	}

	if input.TotalCopies != 0 {
		avaialbe_copies := book.Available_copies + input.TotalCopies - book.Total_copies
		if avaialbe_copies < 0 {
			book.Available_copies = input.TotalCopies
		}
		book.Total_copies = input.TotalCopies

		flag = false
	}
	if input.Available_copies != 0 {
		book.Available_copies = input.Available_copies
		if book.Available_copies > book.Total_copies {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Available copies can not be more than toatl copies",
			})
			return
		}
		flag = false
	}

	// TODO If nothing to update, return error
	if flag {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nothing to update",
		})
		return
	}

	config.DB.Save(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// DELETING A BOOK
func DeleteBook(c *gin.Context) {
	isbnStr := c.Param("isbn")

	var isbn uint
	_, err := fmt.Sscanf(isbnStr, "%d", &isbn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ISBN format"})
		return
	}

	// GET ADMIN ID
	email, _ := c.Get("email")
	var admin models.User
	if err := config.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// SEARCHING THE BOOK
	var book models.Books
	if err := config.DB.Where("isbn = ? AND lib_id = ?", isbn, admin.LibID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// CHECKING IF BOOK HAS BEEN ISSUED
	var issuedCount int64
	config.DB.Raw(`
		SELECT COUNT(*) FROM issue_requests 
		WHERE book_id = ? AND status = 'Issued'`, book.ISBN).Scan(&issuedCount)

	if issuedCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete the book. It is currently issued."})
		return
	}

	config.DB.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
