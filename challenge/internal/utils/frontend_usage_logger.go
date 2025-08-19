package utils

import (
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/transactions"
	"log"
	"time"

	"github.com/google/uuid"
)

type UsageLogger struct {
	usageChan   chan models.FrontendUsage
	memberTrans transactions.MemberTransactions
}

func NewUsageLogger(memberTrans transactions.MemberTransactions) *UsageLogger {
	ul := &UsageLogger{
		usageChan:   make(chan models.FrontendUsage, 1000),
		memberTrans: memberTrans,
	}

	go ul.processUsage()
	return ul
}

func (ul *UsageLogger) LogUsage(userID uuid.UUID) {
	usage := models.FrontendUsage{
		ID:        uuid.New(),
		UserID:    userID,
		Timestamp: time.Now(),
	}

	select {
	case ul.usageChan <- usage:
		// Success, loop!
	default:
		// Drop if full, DON'T block.
		// Don't bother erroring or panicking.
	}
}

// Goroutine running to constantly either:
// - batch insert each time 50 usages are logged.
// OR (whichever comes first every N seconds)
//   - batch insert each every N seconds.
//     note: the overhead of attempting to insert when the batch is empty
//     is minimal, empty queue -> it never opens a connection.
func (ul *UsageLogger) processUsage() {
	batch := make([]models.FrontendUsage, 0, 50)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop() // for cleanliness, but this won't ever happen.

	for {
		select {
		case usage := <-ul.usageChan:
			batch = append(batch, usage)

			if len(batch) >= 50 {
				ul.batchInsert(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				ul.batchInsert(batch)
				batch = batch[:0]
			}
		}
	}
}

func (ul *UsageLogger) batchInsert(batch []models.FrontendUsage) {
	if len(batch) == 0 {
		return
	}

	if err := ul.memberTrans.BatchInsertFrontendUsage(batch); err != nil {
		// Whateva, don't sweat it.
		log.Printf("Failed to batch insert frontend usage: %v", err)
	}
}
