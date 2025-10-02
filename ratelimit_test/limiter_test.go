package ratelimit_test

import (
	"testing"
	"time"

	ratelimit "github.com/seu-usuario/rate-limited-notification-service/ratelimit"
)

type MockClock struct {
	CurrentTime time.Time
}

func (c *MockClock) Now() time.Time {
	return c.CurrentTime
}

func (c *MockClock) Advance(d time.Duration) {
	c.CurrentTime = c.CurrentTime.Add(d)
}

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	startTime := time.Date(2025, time.January, 1, 10, 0, 0, 0, time.UTC)
	mockClock := &MockClock{CurrentTime: startTime}

	rules := map[string]ratelimit.Rule{
		"Status": {Limit: 2, Duration: time.Minute},
		"News":   {Limit: 1, Duration: 24 * time.Hour},
	}

	limiter := ratelimit.NewLimiter(rules, mockClock)

	if !limiter.Allow("Status", "user-A") {
		t.Fatal("Expected: Allowed. Got: Rejected on 1st 'Status' send.")
	}
	if !limiter.Allow("News", "user-B") {
		t.Fatal("Expected: Allowed. Got: Rejected on 1st 'News' send.")
	}
	if !limiter.Allow("Status", "user-A") {
		t.Fatal("Expected: Allowed on 2nd 'Status' send. Got: Rejected. (Limit=2).")
	}
	if limiter.Allow("News", "user-B") {
		t.Fatal("Expected: Rejected. Got: Allowed on 2nd 'News' send (Limit=1).")
	}

	if limiter.Allow("Status", "user-A") {
		t.Fatal("Expected: Rejected. Got: Allowed on 3rd 'Status' send (Limit=2).")
	}

	mockClock.Advance(2*time.Minute + 1*time.Second)

	if !limiter.Allow("Status", "user-A") {
		t.Fatal("Expected: Allowed after window expired. Got: Rejected.")
	}
}
