package notification

import "context"

type Notifier interface {
	BookingCancelled(ctx context.Context, email string, eventTitle string) error
}