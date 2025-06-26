package models

import (
	"time"

	"github.com/google/uuid"
)

// Represents a user registering for the technical challenge.
type Member struct {
	// Primary key, users must save it if they wish to get their status
	ID uuid.UUID `gorm:"primaryKey"`
	// Northeastern email.
	Email string
	// NUID of the user.
	Nuid string
	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateMember(email string, nuid string) *Member {
	user := &Member{}
	user.ID = uuid.New()
	user.Email = email
	user.Nuid = nuid
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}
