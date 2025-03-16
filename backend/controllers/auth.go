package controllers

import (
	"net/http"
	"os"

	"github.com/prabhatKr-1/lib-man-sys/backend/config"
	"github.com/prabhatKr-1/lib-man-sys/backend/models"
	"github.com/prabhatKr-1/lib-man-sys/backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var library models.Library

// CREATING OWNER ACCOUNT
func Signup(c *gin.Context) {
	var input struct {
		Name, Email, Password, ContactNumber, LibraryName string `binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CHECKING IF LIBRARY WITH SAME NAME EXISTS OR NOT
	var existing models.Library
	if err := config.DB.Where("name = ?", input.LibraryName).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library already exists"})
		return
	}

	// LIBRARY CREATION
	library := models.Library{Name: input.LibraryName}
	config.DB.Create(&library)

	// PASSWORD HASHING
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:           "Owner",
		LibID:          library.LibID,
	}
	config.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Owner account created successfully"})
}

// USER LOGIN
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CHECKING IF USER EXISTS
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// PASSWORD VERIFICATION
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TOKEN GENERATION
	token, _ := utils.GenerateJWT(user.ID, user.LibID, user.Email, user.Role)

	// FOR SETTING SECURE SITE
	prodMode := os.Getenv("PROD_MODE") == "true"
	c.SetCookie("token", token, 3600*72, "/", "localhost", prodMode, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// CREATING ADMIN USER
func CreateAdminUser(c *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		ContactNumber string `json:"contactNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VALIDATING OWNER USING COOKIES
	ownerEmail, _ := c.Get("email")
	var owner models.User
	if err := config.DB.Where("email = ?", ownerEmail).First(&owner).Error; err != nil || owner.Role != "Owner" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	adminUser := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:           "Admin",
		LibID:          owner.LibID,
	}

	if err := config.DB.Create(&adminUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin user created successfully"})
}

// CREATING READER USER
func CreateReaderUser(c *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		ContactNumber string `json:"contactNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VALIDATING ADMIN USER USING COOKIES
	adminEmail, _ := c.Get("email")
	var admin models.User
	if err := config.DB.Where("email = ?", adminEmail).First(&admin).Error; err != nil || !(admin.Role == "Admin" || admin.Role == "Owner") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	readerUser := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:           "Reader",
		LibID:          admin.LibID,
	}

	if err := config.DB.Create(&readerUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reader user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reader user created successfully"})
}

// UPDATE USER PASSWORD
func UpdatePassword(c *gin.Context) {
	var input struct {
		oldPassword string `binding:"required"`
		newPassword string `binding:"required"`
	}

	// INPUT VALIDATION
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	email, _ := c.Get("email")
	// FETCHING USER INFORMATION
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// EXISTING PASSWORD VERIFICATION
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.oldPassword)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Password"})
		return
	}

	// NEW PASSWORD HASHING AND SAVING
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.newPassword), bcrypt.DefaultCost) //; err != nil {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bcrypt failed to generate password!",
		})
		return
	}

	user.Password = string(hashedPassword)
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Pasword updated successfully!",
	})

}

// LOGOUT FUNCTIONALITY
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
