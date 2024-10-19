package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/v1adhope/flights/internal/entities"
)

func (u *Usecases) CreateTicket(ctx context.Context, ticket entities.Ticket) (entities.Id, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return entities.Id{}, fmt.Errorf("usecases: ticket: CreateTicket: NewV6: %w", err)
	}

	ticket.Id = id.String()

	ticket.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := u.repos.CreateTicket(ctx, ticket); err != nil {
		return entities.Id{}, err
	}

	return entities.Id{ticket.Id}, nil
}

func (u *Usecases) ReplaceTicket(ctx context.Context, ticket entities.Ticket) error {
	if err := u.repos.ReplaceTicket(ctx, ticket); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) DeleteTicket(ctx context.Context, id entities.Id) error {
	if err := u.repos.DeleteTicket(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) GetTickets(ctx context.Context) ([]entities.Ticket, error) {
	tickets, err := u.repos.GetTickets(ctx)
	if err != nil {
		return []entities.Ticket{}, err
	}

	return tickets, nil
}
