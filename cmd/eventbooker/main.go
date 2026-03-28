package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/config"
	handler "github.com/MyNameIsWhaaat/event-booker/internal/handler/http"
	"github.com/MyNameIsWhaaat/event-booker/internal/repository/postgres"
	"github.com/MyNameIsWhaaat/event-booker/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	pgDSN := cfg.PGDSN
	httpAddr := cfg.HTTPAddr

	db, err := postgres.Connect(ctx, pgDSN)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer func() {
		_ = db.Master.Close()
		for _, slave := range db.Slaves {
			_ = slave.Close()
		}
	}()

	tx := postgres.NewTransactor(db)

	srv := postgres.NewEventRepository(db.Master)
	bookingsrv := postgres.NewBookingRepository(db.Master)
	userRepo := postgres.NewUserRepository(db.Master)

	eventSvc := service.NewEventService(srv, bookingsrv)
	bookingSvc := service.NewBookingService(tx, srv, bookingsrv, userRepo)

	handler := handler.New(eventSvc, bookingSvc)

	server := &http.Server{
		Addr:              httpAddr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server started on %s", httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			log.Printf("server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown complete")
	}
}
