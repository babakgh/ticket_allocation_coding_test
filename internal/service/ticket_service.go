package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/babakgh/ticket_allocation_coding_test/internal/repository"
	"github.com/google/uuid"
)

type TicketOptionRequest struct {
    Name       string `json:"name"`
    Desc       string `json:"desc"`
    Allocation int    `json:"allocation"`
}

type TicketOptionResponse struct {
    ID         uuid.UUID `json:"id"`
    Name       string    `json:"name"`
    Desc       string    `json:"desc"`
    Allocation int       `json:"allocation"`
}

type PurchaseRequest struct {
    Quantity int       `json:"quantity"`
    UserID   uuid.UUID `json:"user_id"`
}

type TicketService interface {
    CreateTicketOption(ctx context.Context, req TicketOptionRequest) (TicketOptionResponse, error)
    GetTicketOption(ctx context.Context, id uuid.UUID) (TicketOptionResponse, error)
    PurchaseTickets(ctx context.Context, id uuid.UUID, req PurchaseRequest) error
}

type ticketService struct {
    repo repository.TicketRepository
}

// Ensure ticketService implements TicketService
var _ TicketService = (*ticketService)(nil)

func NewTicketService(repo repository.TicketRepository) TicketService {
    return &ticketService{repo: repo}
}

func (s *ticketService) CreateTicketOption(ctx context.Context, req TicketOptionRequest) (TicketOptionResponse, error) {
    if req.Name == "" {
        return TicketOptionResponse{}, errors.New("ticket option name cannot be empty")
    }
    if req.Allocation <= 0 {
        return TicketOptionResponse{}, errors.New("ticket allocation must be positive")
    }

    repoOption := repository.TicketOption{
        Name:       req.Name,
        Desc:       req.Desc,
        Allocation: req.Allocation,
    }

    createdOption, err := s.repo.CreateTicketOption(ctx, repoOption)
    if err != nil {
        log.Printf("Error creating ticket option: %v", err)
        return TicketOptionResponse{}, fmt.Errorf("failed to create ticket option: %w", err)
    }

    response := TicketOptionResponse{
        ID:         createdOption.ID,
        Name:       createdOption.Name,
        Desc:       createdOption.Desc,
        Allocation: createdOption.Allocation,
    }

    log.Printf("Ticket option created: %s", response.ID)
    return response, nil
}

func (s *ticketService) GetTicketOption(ctx context.Context, id uuid.UUID) (TicketOptionResponse, error) {
    option, err := s.repo.GetTicketOption(ctx, id)
    if err != nil {
        log.Printf("Error retrieving ticket option: %v", err)
        return TicketOptionResponse{}, fmt.Errorf("failed to retrieve ticket option: %w", err)
    }

    return TicketOptionResponse{
        ID:         option.ID,
        Name:       option.Name,
        Desc:       option.Desc,
        Allocation: option.Allocation,
    }, nil
}

func (s *ticketService) PurchaseTickets(ctx context.Context, id uuid.UUID, req PurchaseRequest) error {
    if req.Quantity <= 0 {
        return errors.New("purchase quantity must be positive")
    }

    purchase := repository.Purchase{
        ID:             uuid.New(),
        Quantity:       req.Quantity,
        UserID:         req.UserID,
        TicketOptionID: id,
    }

    err := s.repo.PurchaseTickets(ctx, id, purchase)
    if err != nil {
        if errors.Is(err, repository.ErrInsufficientTickets) {
            log.Printf("Insufficient tickets for purchase: %s, quantity: %d", id, req.Quantity)
            return fmt.Errorf("insufficient tickets available: %w", err)
        }
        log.Printf("Error purchasing tickets: %v", err)
        return fmt.Errorf("failed to purchase tickets: %w", err)
    }

    log.Printf("Tickets purchased: %s, quantity: %d, user: %s", id, req.Quantity, req.UserID)
    return nil
}
