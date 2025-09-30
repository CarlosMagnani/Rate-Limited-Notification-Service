package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Rules that define the limit and the time window by type of notification
type Rule struct {
	Limit    int
	Duration time.Duration
}

// Definition of contract to service of notification limit
type RateLimiter interface {
	Allow(notificationType, recipientID string) bool
}

//Implements the RateLimiter using the Sliding Window algorithm

type SlidingWindowLimiter struct {
	rules map[string]Rule
	mu    sync.Mutex
	state map[string][]time.Time
}

func NewLimiter(r map[string]Rule) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		rules: r,
		state: make(map[string][]time.Time),
	}
}

func (l *SlidingWindowLimiter) Allow(notificationType, recipientID string) bool {
	rule, exists := l.rules[notificationType]
	if !exists {
		return true
	}

	key := fmt.Sprintf("%s:%s", notificationType, recipientID)

	//Concurrency protection
	l.mu.Lock()
	defer l.mu.Unlock()

	currentLogs := l.state[key]

	//1. Initial point cutoff window
	cutoff := time.Now().Add(-rule.Duration)

	startIdx := 0

	//Cleanup logs
	for i, t := range currentLogs {
		if t.After(cutoff) {
			startIdx = i
			break
		}
	}

	cleanLogs := currentLogs[startIdx:]

	//Check if the notification is in the limit
	if len(cleanLogs) < rule.Limit {
		l.state[key] = append(cleanLogs, time.Now())
		return true
	}

	// Keeps clean log and return false
	l.state[key] = currentLogs
	return false
}
