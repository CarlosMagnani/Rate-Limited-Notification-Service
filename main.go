package main

import (
	"fmt"
	"time"

	"github.com/seu-usuario/rate-limited-notification-service/notification"
	ratelimit "github.com/seu-usuario/rate-limited-notification-service/ratelimit"
)

func main() {
	fmt.Println("--Initializing Rate-Limited Notification Service--")

	rules := map[string]ratelimit.Rule{
		"Status":    {Limit: 2, Duration: time.Minute},
		"News":      {Limit: 1, Duration: 24 * time.Hour},
		"Marketing": {Limit: 3, Duration: 3 * time.Hour},
	}

	realClock := ratelimit.RealClock{}

	limiter := ratelimit.NewLimiter(rules, realClock)
	gateway := &notification.EmailGateway{}

	service := notification.NewServiceImpl(limiter, gateway)

	fmt.Println("\n--- Test: STATUS (Limit: 2 per minute) ---")

	service.Send("Status", "user-A", "Update 1")
	service.Send("Status", "user-A", "Update 2")
	service.Send("Status", "user-A", "Update 3")

	service.Send("Status", "user-B", "Update 1 B")

	fmt.Println("\n--- Test: NEWS (Limit: 1 per day) ---")

	service.Send("News", "user-C", "Daily News 1")
	service.Send("News", "user-C", "Daily News 2")

	fmt.Println("\n--- Test: Sliding Window (Immediate Re-check) ---")

	service.Send("Status", "user-A", "Update 4 immediately after")
}
