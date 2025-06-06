package entity

import (
	"time"
)

// Audiobook represents the audiobooks table
type Audiobook struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title            string    `json:"title" gorm:"size:255;not null"`
	AuthorID         uint      `json:"author_id" gorm:"not null"`
	ReaderID         uint      `json:"reader_id" gorm:"not null"`
	Description      string    `json:"description" gorm:"type:text"`
	ImageURL         string    `json:"image_url" gorm:"size:255"`
	Language         string    `json:"language" gorm:"size:50"`
	YearOfPublishing int       `json:"year_of_publishing"`
	TotalDuration    string    `json:"total_duration" gorm:"size:50"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Author    Author      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Reader    Reader      `json:"reader,omitempty" gorm:"foreignKey:ReaderID"`
	Tracks    []Track     `json:"tracks,omitempty" gorm:"foreignKey:AudiobookID"`
	Genres    []Genre     `json:"genres,omitempty" gorm:"many2many:audiobook_genres;"`
	Analytics []Analytics `json:"analytics,omitempty" gorm:"foreignKey:AudiobookID"`
}

// TableName specifies the table name for the Audiobook model
func (Audiobook) TableName() string {
	return "audiobooks"
}
