package worker

import (
	"context"
	"log"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/service"
)

type BookingExpirer struct {
	bookingSvc service.BookingService
	interval   time.Duration
}

func NewBookingExpirer(bookingSvc service.BookingService, interval time.Duration) *BookingExpirer {
	return &BookingExpirer{
		bookingSvc: bookingSvc,
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
			n, err := w.bookingSvc.CancelExpired(ctx)
			if err != nil {
				log.Printf("cancel expired error: %v", err)
				continue
			}
			if n > 0 {
				log.Printf("cancelled expired bookings: %d", n)
			}
		}
	}
}