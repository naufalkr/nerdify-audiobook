package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"size:50;not null;unique"`
	Description string
	IsSystem    bool `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
