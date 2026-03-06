package notification

import "context"

type NoopNotifier struct{}

func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{}
}

func (n *NoopNotifier) BookingCancelled(ctx context.Context, email string, eventTitle string) error {
	return nil
}