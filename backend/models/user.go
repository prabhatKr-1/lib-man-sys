package models

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"not null" binding:"required" json:"name"`
	Email          string `gorm:"unique;not null" binding:"required" json:"email"`
	Password       string `gorm:"not null" binding:"required" json:"password"`
	Contact_number string `gorm:"not null" binding:"required" json:"contact_number"`
	LibID          uint   `gorm:"not null" json:"lib_id"`
	Role           string `gorm:"not null;check:role IN ('Owner','Admin','Reader')"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
