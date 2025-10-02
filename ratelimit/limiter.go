package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type Rule struct {
	Limit    int
	Duration time.Duration
}

type RateLimiter interface {
	Allow(notificationType, recipientID string) bool
}

type SlidingWindowLimiter struct {
	rules map[string]Rule
	mu    sync.Mutex
	state map[string][]time.Time
	Clock Clock
}

func NewLimiter(r map[string]Rule, clock Clock) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		rules: r,
		state: make(map[string][]time.Time),
		Clock: clock,
	}
}

func (l *SlidingWindowLimiter) Allow(notificationType, recipientID string) bool {
	rule, exists := l.rules[notificationType]
	if !exists {
		return true
	}

	key := fmt.Sprintf("%s:%s", notificationType, recipientID)

	l.mu.Lock()
	defer l.mu.Unlock()

	currentLogs := l.state[key]

	cutoff := l.Clock.Now().Add(-rule.Duration)

	startIdx := len(currentLogs)

	for i, t := range currentLogs {
		if t.After(cutoff) {
			startIdx = i
			break
		}
	}

	cleanLogs := currentLogs[startIdx:]

	if len(cleanLogs) < rule.Limit {
		l.state[key] = append(cleanLogs, l.Clock.Now())
		return true
	}

	l.state[key] = currentLogs
	return false
}
