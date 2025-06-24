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
	NUID string
	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateUser(email string, nuid string) *User {
	user := &User{}
	user.ID = uuid.New()
	user.Email = email
	user.NUID = nuid
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}
