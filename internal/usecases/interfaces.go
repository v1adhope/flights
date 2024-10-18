package usecases

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type Reposer interface {
	CreateTicket(ctx context.Context, ticket entities.Ticket) error
	ReplaceTicket(ctx context.Context, ticket entities.Ticket) error
	DeleteTicket(ctx context.Context, id string) error
	GetAllTickets(ctx context.Context) ([]entities.Ticket, error)
}
