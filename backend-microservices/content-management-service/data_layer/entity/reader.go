package entity

import (
	"time"
)

// Reader represents the readers table
type Reader struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Audiobooks []Audiobook `json:"audiobooks,omitempty" gorm:"foreignKey:ReaderID"`
}

// TableName specifies the table name for the Reader model
func (Reader) TableName() string {
	return "readers"
}
