package models

import "time"

type Books struct {
	ISBN             uint   `gorm:"primaryKey"`
	LibID            uint   `gorm:"primaryKey;autoIncrement:false"`
	Title            string `gorm:"not null" binding:"required" json:"title"`
	Authors          string `gorm:"not null" binding:"required" json:"authors"`
	Publisher        string `gorm:"not null" binding:"required" json:"publisher"`
	Version          string `gorm:"not null" binding:"required" json:"version"`
	Total_copies     uint   `gorm:"not null" binding:"required,min=1" json:"total_copies"`
	Available_copies uint   `gorm:"not null" json:"available_copies"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Library Library `gorm:"foreignKey:LibID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
