package models

import (
	"time"

	"github.com/google/uuid"
)

type AlienChallengeSolution struct {
	ID                uuid.UUID `gorm:"primaryKey"`
	MemberID          uuid.UUID `gorm:"not null;index"`
	IdealCommandsUsed int
	IdealHpRemaining  int
	IdealAliensLeft   int
	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateAlienChallengeSolutionEntry(memberID uuid.UUID, commandsUsed int, hpRemaining int, aliensLeft int) *AlienChallengeSolution {
	sol := &AlienChallengeSolution{}
	sol.ID = uuid.New()
	sol.MemberID = memberID
	sol.IdealCommandsUsed = commandsUsed
	sol.IdealHpRemaining = hpRemaining
	sol.IdealAliensLeft = aliensLeft
	sol.CreatedAt = time.Now()
	sol.UpdatedAt = time.Now()
	return sol
}
