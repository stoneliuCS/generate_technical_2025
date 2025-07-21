package utils

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter manages rate limits per UUID.
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	limit    rate.Limit
	burst    int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		limit:    rate.Every(time.Minute / 10), // 10 requests per minute.
		burst:    10,                           // Allow burst of 10.
	}
}

// GetLimiter returns the rate limiter for a specific UUID.
func (rl *RateLimiter) GetLimiter(uuid string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[uuid]
	if !exists {
		limiter = rate.NewLimiter(rl.limit, rl.burst)
		rl.limiters[uuid] = limiter
	}

	return limiter
}

// Allow checks if the request is allowed for this UUID.
func (rl *RateLimiter) Allow(uuid string) bool {
	limiter := rl.GetLimiter(uuid)
	return limiter.Allow()
}
