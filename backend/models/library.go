package models

import "time"

type Library struct {
	LibID uint   `gorm:"primaryKey"  json:"id"`
	Name  string `binding:"required" gorm:"unique;not null" json:"lib_name"`

	CreatedAt time.Time `json:"created_at"`

	Users []User  `gorm:"foreignKey:LibID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Books []Books `gorm:"foreignKey:LibID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}