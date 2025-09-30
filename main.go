package main

import (
	"fmt"
	"time"

	"github.com/seu-usuario/rate-limited-notification-service/notification"
	ratelimit "github.com/seu-usuario/rate-limited-notification-service/rate-limit"
)

func main() {
	fmt.Println("--Initializing service of rate limit notification--")

	rules := map[string]ratelimit.Rule{
		"Status":     {Limit: 2, Duration: time.Minute},
		"News":       {Limit: 1, Duration: 24 * time.Hour},
		"Markenting": {Limit: 3, Duration: 3 * time.Hour},
	}

	limiter := ratelimit.NewLimiter(rules)
	gateway := &notification.EmailGateway{}

	service := &notification.ServiceImpl{
		Limiter: limiter,
		Gateway: gateway,
	}

	fmt.Println("\n--- Teste STATUS (Limite: 2 por minuto) ---")

	service.Send("Status", "user-A", "Update 1")

	service.Send("Status", "user-A", "Update 2")

	service.Send("Status", "user-A", "Update 3")

	service.Send("Status", "user-B", "Update 1 B")

	fmt.Println("\n--- Teste NEWS (Limite: 1 por dia) ---")

	service.Send("News", "user-C", "Daily News 1")

	service.Send("News", "user-C", "Daily News 2")

	fmt.Println("\n--- Teste Janela Deslizante (Simulando 3 segundos depois) ---")

	time.Sleep(3 * time.Second)

	service.Send("Status", "user-A", "Update 4 após 3s")
}
