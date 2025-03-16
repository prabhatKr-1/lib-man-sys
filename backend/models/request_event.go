package models

import "time"

type RequestEvents struct {
	ReqID          uint      `gorm:"primaryKey"  json:"reqID"`
	BookID         uint      `gorm:"not null" binding:"required" json:"bookID"`
	ReaderID       uint      `gorm:"not null" binding:"required" json:"readerID"`
	RequestDate    time.Time `gorm:"not null"`
	ProcessingDate *time.Time
	AdminID        *uint
	LibID          uint   `gorm:"not null"`
	RequestType    string `gorm:"default:'issue';check:request_type IN ('issue','return')"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Book     Books `gorm:"foreignKey:BookID,LibID;references:ISBN,LibID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reader   User  `gorm:"foreignKey:ReaderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Approver *User `gorm:"foreignKey:AdminID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
