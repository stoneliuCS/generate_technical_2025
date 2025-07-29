package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ALGORITHM_CHALLENGE_TYPE = "algorithm"
	NGROK_CHALLENGE_TYPE     = "ngrok"
	INVALID_SCORE            = -1
)

// Represents a score for a user's takehome challenge submission.
type Score struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	UserID        uuid.UUID `gorm:"not null;index"`
	ChallengeType string    `gorm:"not null"`
	Score         int       `gorm:"not null"`
	IsValid       bool      `gorm:"not null;default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Member Member `gorm:"foreignKey:UserID;references:ID"`
}

func CreateScore(userID uuid.UUID, challengeType string, score int, isValid bool) *Score {
	scoreRecord := &Score{}
	scoreRecord.ID = uuid.New()
	scoreRecord.UserID = userID
	scoreRecord.ChallengeType = challengeType
	scoreRecord.Score = score
	scoreRecord.IsValid = isValid
	scoreRecord.CreatedAt = time.Now()
	scoreRecord.UpdatedAt = time.Now()
	return scoreRecord
}
