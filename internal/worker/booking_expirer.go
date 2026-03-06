package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/MyNameIsWhaaat/event-booker/internal/notification"
	"github.com/MyNameIsWhaaat/event-booker/internal/service"
)

type BookingExpirer struct {
	bookingSvc service.BookingService
	eventSvc   service.EventService
	notifier   notification.Notifier
	interval   time.Duration
}

func NewBookingExpirer(
	bookingSvc service.BookingService,
	eventSvc service.EventService,
	notifier notification.Notifier,
	interval time.Duration,
) *BookingExpirer {
	return &BookingExpirer{
		bookingSvc: bookingSvc,
		eventSvc:   eventSvc,
		notifier:   notifier,
		interval:   interval,
	}
}

func (w *BookingExpirer) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("booking expirer started, interval=%s", w.interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("booking expirer stopped: %v", ctx.Err())
			return

		case <-ticker.C:
			items, err := w.bookingSvc.CancelExpired(ctx)
			if err != nil {
				log.Printf("cancel expired error: %v", err)
				continue
			}

			if len(items) == 0 {
				continue
			}

			log.Printf("cancelled expired bookings: %d", len(items))

			for _, b := range items {
				eventID := b.EventID
				_ = uuid.Nil

				eventUUID, err := uuid.Parse(eventID)
				if err != nil {
					log.Printf("invalid event id: %s", eventID)
					continue
				}

				event, err := w.eventSvc.GetEvent(ctx, eventUUID)
				if err != nil {
					log.Printf("failed to load event for booking %s: %v", b.ID, err)
					continue
				}

				err = w.notifier.BookingCancelled(ctx, b.UserEmail, event.Title)
				if err != nil {
					log.Printf("failed to send cancellation email to %s: %v", b.UserEmail, err)
					continue
				}

				log.Printf("cancellation email sent to %s for booking %s", b.UserEmail, b.ID)
			}
		}
	}
}
