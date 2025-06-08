package entity

import (
	"time"
)

// Audiobook represents the audiobooks table
type Audiobook struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title            string    `gorm:"not null" json:"title"`
	AuthorID         uint      `gorm:"not null" json:"author_id"`
	ReaderID         uint      `gorm:"not null" json:"reader_id"`
	Description      string    `gorm:"type:text" json:"description"`
	ImageURL         string    `json:"image_url"`
	Language         string    `json:"language"`
	YearOfPublishing int       `json:"year_of_publishing"`
	TotalDuration    string    `json:"total_duration"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relationships
	Author    *Author      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Reader    *Reader      `gorm:"foreignKey:ReaderID" json:"reader,omitempty"`
	Genres    []Genre     `gorm:"many2many:audiobook_genres" json:"genres,omitempty"`
	Tracks    []Track     `gorm:"foreignKey:AudiobookID" json:"tracks,omitempty"`
}

// TableName specifies the table name for the Audiobook model
func (Audiobook) TableName() string {
	return "audiobooks"
}
