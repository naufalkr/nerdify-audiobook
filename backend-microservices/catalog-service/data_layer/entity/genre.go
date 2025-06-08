package entity

import (
	"time"
)

// Genre represents the genres table
type Genre struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100;not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Many-to-Many relationship with Audiobooks
	Audiobooks []Audiobook `json:"audiobooks,omitempty" gorm:"many2many:audiobook_genres;"`
}

// TableName specifies the table name for the Genre model
func (Genre) TableName() string {
	return "genres"
}
