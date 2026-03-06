package notification

import (
	"context"
	"fmt"
	"net/smtp"
)

type EmailNotifier struct {
	host string
	port string
	from string
}

func NewEmailNotifier(host, port, from string) *EmailNotifier {
	return &EmailNotifier{
		host: host,
		port: port,
		from: from,
	}
}

func (n *EmailNotifier) BookingCancelled(ctx context.Context, email string, eventTitle string) error {
	addr := fmt.Sprintf("%s:%s", n.host, n.port)

	subject := "Subject: Booking cancelled\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(
		"Hello!\r\n\r\nYour booking for event \"%s\" was cancelled because it was not confirmed in time.\r\n",
		eventTitle,
	)

	msg := []byte(subject + mime + body)

	return smtp.SendMail(
		addr,
		nil,
		n.from,
		[]string{email},
		msg,
	)
}
