package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey"`
	TenantID   uuid.UUID  `gorm:"type:char(36)"`
	UserID     *uuid.UUID `gorm:"type:char(36)"`
	Action     string     `gorm:"size:50"`
	EntityType string     `gorm:"size:100"`
	EntityID   uuid.UUID  `gorm:"type:char(36)"`
	OldValues  datatypes.JSON
	NewValues  datatypes.JSON
	IPAddress  string `gorm:"size:45"`
	UserAgent  string

	CreatedAt time.Time
}
