package entity

import (
	"time"
)

// Author represents the authors table
type Author struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Audiobooks []Audiobook `json:"audiobooks,omitempty" gorm:"foreignKey:AuthorID"`
}

// TableName specifies the table name for the Author model
func (Author) TableName() string {
	return "authors"
}
