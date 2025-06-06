package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserTenant represents the many-to-many relationship between users and tenants
type UserTenant struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID"`
}
