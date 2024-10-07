package repository

import (
	"context"
	"errors"
	"time"

	"github.com/babakgh/ticket_allocation_coding_test/pkg/database"
	"github.com/google/uuid"
)

var ErrInsufficientTickets = errors.New("insufficient tickets available")

type TicketOption struct {
    ID         uuid.UUID
    Name       string
    Desc       string
    Allocation int
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type Purchase struct {
    ID             uuid.UUID
    Quantity       int
    UserID         uuid.UUID
    TicketOptionID uuid.UUID
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type TicketRepository interface {
    CreateTicketOption(ctx context.Context, option TicketOption) (TicketOption, error)
    GetTicketOption(ctx context.Context, id uuid.UUID) (TicketOption, error)
    PurchaseTickets(ctx context.Context, id uuid.UUID, purchase Purchase) error
}

type ticketRepository struct {
    db *database.Database
}

func NewTicketRepository(db *database.Database) TicketRepository {
    return &ticketRepository{db: db}
}

func (r *ticketRepository) CreateTicketOption(ctx context.Context, option TicketOption) (TicketOption, error) {
    err := r.db.GetDB().QueryRowContext(ctx, `
        INSERT INTO ticket_options (name, "desc", allocation)
        VALUES ($1, $2, $3)
        RETURNING id, name, "desc", allocation, created_at, updated_at
    `, option.Name, option.Desc, option.Allocation).Scan(
        &option.ID, &option.Name, &option.Desc, &option.Allocation, &option.CreatedAt, &option.UpdatedAt,
    )

    return option, err
}

func (r *ticketRepository) GetTicketOption(ctx context.Context, id uuid.UUID) (TicketOption, error) {
    var option TicketOption
    err := r.db.GetDB().QueryRowContext(ctx, `
        SELECT id, name, "desc", allocation, created_at, updated_at
        FROM ticket_options
        WHERE id = $1
    `, id).Scan(&option.ID, &option.Name, &option.Desc, &option.Allocation, &option.CreatedAt, &option.UpdatedAt)

    return option, err
}

func (r *ticketRepository) PurchaseTickets(ctx context.Context, id uuid.UUID, purchase Purchase) error {
    tx, err := r.db.GetDB().BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    var currentAllocation int
    err = tx.QueryRowContext(ctx, `
        SELECT allocation
        FROM ticket_options
        WHERE id = $1
        FOR UPDATE
    `, id).Scan(&currentAllocation)

    if err != nil {
        return err
    }

    if currentAllocation < purchase.Quantity {
        return ErrInsufficientTickets
    }

    _, err = tx.ExecContext(ctx, `
        UPDATE ticket_options
        SET allocation = allocation - $1
        WHERE id = $2
    `, purchase.Quantity, id)

    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx, `
        INSERT INTO purchases (id, quantity, user_id, ticket_option_id)
        VALUES ($1, $2, $3, $4)
    `, purchase.ID, purchase.Quantity, purchase.UserID, id)

    if err != nil {
        return err
    }

    for i := 0; i < purchase.Quantity; i++ {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO tickets (ticket_option_id, purchase_id)
            VALUES ($1, $2)
        `, id, purchase.ID)

        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
