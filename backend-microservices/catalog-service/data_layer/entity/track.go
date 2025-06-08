package entity

import (
	"time"
)

// Track represents the tracks table
type Track struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AudiobookID uint      `json:"audiobook_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	URL         string    `json:"url" gorm:"size:255;not null"`
	Duration    string    `json:"duration" gorm:"size:20"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Audiobook Audiobook `json:"audiobook,omitempty" gorm:"foreignKey:AudiobookID"`
}

// TableName specifies the table name for the Track model
func (Track) TableName() string {
	return "tracks"
}
