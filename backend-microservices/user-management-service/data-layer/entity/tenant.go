package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key"`
	Name                  string    `gorm:"size:255;not null"`
	Description           string    `gorm:"size:1000"`
	LogoURL               string    `gorm:"size:255"`
	ContactEmail          string    `gorm:"size:255"`
	ContactPhone          string    `gorm:"size:20"`
	MaxUsers              int       `gorm:"not null;default:0"`
	SubscriptionPlan      string    `gorm:"size:20"`
	SubscriptionStartDate time.Time
	SubscriptionEndDate   time.Time
	IsActive              bool           `gorm:"not null;default:true"`
	CreatedAt             time.Time      `gorm:"not null"`
	UpdatedAt             time.Time      `gorm:"not null"`
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}
