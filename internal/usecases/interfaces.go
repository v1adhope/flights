package usecases

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type Reposer interface {
	// INFO: Tickets
	CreateTicket(ctx context.Context, ticket entities.Ticket) error
	ReplaceTicket(ctx context.Context, ticket entities.Ticket) error
	DeleteTicket(ctx context.Context, id entities.Id) error
	GetTickets(ctx context.Context) ([]entities.Ticket, error)

	// INFO: Passengers
	CreatePassenger(ctx context.Context, passenger entities.Passenger) error
	ReplacePassenger(ctx context.Context, passenger entities.Passenger) error
	DeletePassenger(ctx context.Context, id entities.Id) error
	GetPassengers(ctx context.Context) ([]entities.Passenger, error)

	// INFO: documents
	CreateDocument(ctx context.Context, document entities.Document) error
	ReplaceDocument(ctx context.Context, document entities.Document) error
	DeleteDocument(ctx context.Context, id entities.Id) error
	GetDocumentsByPassengerId(ctx context.Context, id entities.Id) ([]entities.Document, error)
}
