package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/babakgh/ticket_allocation_coding_test/internal/app"
	"github.com/babakgh/ticket_allocation_coding_test/internal/config"
)

var (
	timeout = flag.Duration("timeout", 1*time.Second, "timeout for graceful shutdown")
)

func main() {
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run the application
	run(cfg)

	// If it's here then the server has been stopped gracefully
	log.Println("Server gracefully stopped")
}

// TODO: Too many lines, we can refactor this function.
func run(cfg *config.Config) {
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.Close()

	log.Println("Starting the server...")
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down app: %v", err)
	}
}
