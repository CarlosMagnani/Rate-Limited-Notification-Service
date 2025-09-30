package notification

import (
	"fmt"

	ratelimit "github.com/seu-usuario/rate-limited-notification-service/rate-limit"
)

type NotificationService interface {
	Send(notifType, userId, message string)
}

type ServiceImpl struct {
	Limiter ratelimit.RateLimiter
	Gateway Gateway
}

func (s *ServiceImpl) Send(notifyType, userId, message string) {
	if !s.Limiter.Allow(notifyType, userId) {
		fmt.Printf("[REJECT] Limit exceed to user %s notification type '%s'.\n", userId, notifyType)
		return
	}

	s.Gateway.Send(userId, message)
}
