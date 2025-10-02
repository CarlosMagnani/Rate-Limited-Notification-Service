package ratelimit

import "time"

// Clock permite mockar tempo nos testes.
type Clock interface {
	Now() time.Time
}

// RealClock usa time.Now.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }
