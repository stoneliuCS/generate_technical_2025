package models

import (
	"time"

	"github.com/google/uuid"
)

// Represents a user registering for the technical challenge.
type User struct {
	// Primary key, users must save it if they wish to get their status
	ID uuid.UUID `gorm:"primaryKey"`
	// Northeastern email.
	Email string
	// NUID of the user.
	NUID      string
	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}
