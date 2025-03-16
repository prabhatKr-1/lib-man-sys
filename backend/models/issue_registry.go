package models

import "time"

type IssueRegistry struct {
	IssueID         uint `gorm:"primaryKey" json:"issueID"`
	ISBN            uint `gorm:"not null" binding:"required" json:"isbn"`
	LibID           uint `gorm:"not null"`
	ReaderID        uint `gorm:"not null" binding:"required" json:"readerID"`
	IssueApproverID uint `gorm:"not null" binding:"required" json:"issueapproverID"`
	Status             string     `gorm:"not null;check:status IN ('issued','returned')" json:"status"`
	IssueDate          time.Time  `gorm:"not null" binding:"required" json:"date"`
	ExpectedReturnDate time.Time  `gorm:"not null" binding:"required" json:"expected_return_date"`
	ReturnDate         *time.Time `json:"return_date"`
	ReturnApproverID   *uint      `json:"returnapproverID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Book           Books `gorm:"foreignKey:ISBN,LibID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reader         User  `gorm:"foreignKey:ReaderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IssueApprover  User  `gorm:"foreignKey:IssueApproverID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	ReturnApprover *User `gorm:"foreignKey:ReturnApproverID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}
