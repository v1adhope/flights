package usecases

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type Reposer interface {
	Ticket
	Passenger
	Document
	Report
}

type (
	Ticket interface {
		CreateTicket(ctx context.Context, ticket entities.Ticket) error
		ReplaceTicket(ctx context.Context, ticket entities.Ticket) error
		DeleteTicket(ctx context.Context, id entities.Id) error
		GetTickets(ctx context.Context) ([]entities.Ticket, error)
		GetWholeInfoAboutTicket(ctx context.Context, id entities.Id) (entities.TicketWholeInfo, error)
	}

	Passenger interface {
		CreatePassenger(ctx context.Context, passenger entities.Passenger) error
		ReplacePassenger(ctx context.Context, passenger entities.Passenger) error
		DeletePassenger(ctx context.Context, id entities.Id) error
		BoundToTicket(ctx context.Context, id entities.Id, ticketId entities.Id) error
		UnboundToTicket(ctx context.Context, id entities.Id, ticketId entities.Id) error
		GetPassengersByTicketId(ctx context.Context, id entities.Id) ([]entities.Passenger, error)

		GetPassengers(ctx context.Context) ([]entities.Passenger, error)
	}

	Document interface {
		CreateDocument(ctx context.Context, document entities.Document) error
		ReplaceDocument(ctx context.Context, document entities.Document) error
		DeleteDocument(ctx context.Context, id entities.Id) error
		GetDocumentsByPassengerId(ctx context.Context, id entities.Id) ([]entities.Document, error)
	}

	Report interface {
		GetRowsByPassengerIdForPeriod(ctx context.Context, id entities.Id, filter entities.PeriodFilter) ([]entities.ReportRowByPassengerForPeriod, error)
	}
)
