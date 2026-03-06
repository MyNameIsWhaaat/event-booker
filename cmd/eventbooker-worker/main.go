package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/notification"
	"github.com/MyNameIsWhaaat/event-booker/internal/repository/postgres"
	"github.com/MyNameIsWhaaat/event-booker/internal/service"
	"github.com/MyNameIsWhaaat/event-booker/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := postgres.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	tx := postgres.NewTransactor(db)
	eventRepo := postgres.NewEventRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	userRepo := postgres.NewUserRepository(db)

	eventSvc := service.NewEventService(eventRepo, bookingRepo)
	bookingSvc := service.NewBookingService(tx, eventRepo, bookingRepo, userRepo)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpFrom := os.Getenv("SMTP_FROM")

	var notifier notification.Notifier
	if smtpHost != "" && smtpPort != "" && smtpFrom != "" {
		notifier = notification.NewEmailNotifier(smtpHost, smtpPort, smtpFrom)
	} else {
		notifier = notification.NewNoopNotifier()
	}

	expirer := worker.NewBookingExpirer(
		bookingSvc,
		eventSvc,
		notifier,
		5*time.Second,
	)

	expirer.Run(ctx)
}
