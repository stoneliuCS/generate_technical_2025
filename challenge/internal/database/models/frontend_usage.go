package models

import (
	"time"

	"github.com/google/uuid"
)

type FrontendUsage struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Timestamp time.Time `gorm:"not null;default:now()" json:"timestamp"`
}
