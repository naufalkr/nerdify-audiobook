package entity

import (
	"time"
)

// User represents the users table
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;size:255"` // UUID from auth service
	Role      string    `json:"role" gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Analytics []Analytics `json:"analytics,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
