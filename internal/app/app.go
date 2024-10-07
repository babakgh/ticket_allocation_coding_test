package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/babakgh/ticket_allocation_coding_test/internal/config"
	"github.com/babakgh/ticket_allocation_coding_test/internal/handler"
	"github.com/babakgh/ticket_allocation_coding_test/internal/repository"
	"github.com/babakgh/ticket_allocation_coding_test/internal/service"
	"github.com/babakgh/ticket_allocation_coding_test/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	config *config.Config
	db     *database.Database
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	repo := repository.NewTicketRepository(db)
	svc := service.NewTicketService(repo)
	h := handler.NewHandler(svc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	h.Register(e)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: e,
	}

	// Print all routes
	for _, route := range e.Routes() {
		log.Printf("%s %s\n", route.Method, route.Path)
	}

	return &App{
		config: cfg,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Start() error {
	listener, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return fmt.Errorf("Failed to create listener: %w", err)
	}

	log.Printf("Server is started at %s ...\n", a.server.Addr)

	// TODO: tls support here
	go func() {
		if err := a.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			// TODO: this is a bit lazy, will leave the db open, should use channels to communicate with the main goroutine
			log.Fatalf("Failed to start server %s", err)
		}
	}()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %w", err)
	}

	return nil
}

func (a *App) Close() {
	log.Println("Closing connections...")
	if err := a.db.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}
