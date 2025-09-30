package notification

import "fmt"

type Gateway interface {
	Send(userId, message string)
}

type EmailGateway struct{}

func (g *EmailGateway) Send(userId, message string) {
	fmt.Printf("[OK] SENDED: '%s' to user %s.\n", message, userId)
}
