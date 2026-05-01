package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nats/internal/config"
	"nats/internal/handler"
	"nats/internal/nats"
)

func main() {
	cfg := config.Load()

	jsClient, err := nats.New(cfg.NatsURL, cfg.Stream, cfg.Subject)
	if err != nil {
		log.Fatalf("failed to init nats: %v", err)
	}

	handler := handler.New(jsClient)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP server started on :%s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
		os.Exit(1)
	}
}
