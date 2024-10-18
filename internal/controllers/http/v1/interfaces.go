package v1

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type TicketUsecaser interface {
	CreateTicket(ctx context.Context, ticket entities.Ticket) (string, error)
	ReplaceTicket(ctx context.Context, ticket entities.Ticket) error
	DeleteTicket(ctx context.Context, id string) error
	GetAllTickets(ctx context.Context) ([]entities.Ticket, error)
}

type PassengerUsecaser interface {
	CreatePassenger(ctx context.Context, passenger entities.Passenger) error
	ReplacePassenger(ctx context.Context, passenger entities.Passenger) error
	DeletePassenger(ctx context.Context, id string) error
}

type Logger interface {
	Debug(err error, format string, msg ...any)
	Error(err error, format string, msg ...any)
}
