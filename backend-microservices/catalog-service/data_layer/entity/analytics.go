package entity

import (
	"time"
)

// Analytics represents the analytics table
type Analytics struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AudiobookID    uint      `json:"audiobook_id" gorm:"not null"`
	UserID         string    `json:"user_id" gorm:"size:255;not null"`
	EventType      string    `json:"event_type" gorm:"size:50;not null"`
	EventTimestamp time.Time `json:"event_timestamp" gorm:"not null"`

	// Relationships
	Audiobook Audiobook `json:"audiobook,omitempty" gorm:"foreignKey:AudiobookID"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the Analytics model
func (Analytics) TableName() string {
	return "analytics"
}
