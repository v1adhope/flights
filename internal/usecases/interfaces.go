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

	CreatePassenger(ctx context.Context, passenger entities.Passenger) error
	ReplacePassenger(ctx context.Context, passenger entities.Passenger) error
	DeletePassenger(ctx context.Context, id string) error
	GetAllPassengers(ctx context.Context) ([]entities.Passenger, error)
}
