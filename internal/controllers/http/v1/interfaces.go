package v1

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type TicketUsecaser interface {
	CreateTicket(ctx context.Context, ticket entities.Ticket) (entities.Id, error)
	ReplaceTicket(ctx context.Context, ticket entities.Ticket) error
	DeleteTicket(ctx context.Context, id entities.Id) error
	GetTickets(ctx context.Context) ([]entities.Ticket, error)
	GetWholeInfoAboutTicket(ctx context.Context, id entities.Id) (entities.TicketWholeInfo, error)
}

type PassengerUsecaser interface {
	CreatePassenger(ctx context.Context, passenger entities.Passenger) (entities.Id, error)
	ReplacePassenger(ctx context.Context, passenger entities.Passenger) error
	DeletePassenger(ctx context.Context, id entities.Id) error
	BoundToTicket(ctx context.Context, id entities.Id, ticketId entities.Id) error
	UnboundToTicket(ctx context.Context, id entities.Id, ticketId entities.Id) error
	GetPassengersByTicketId(ctx context.Context, id entities.Id) ([]entities.Passenger, error)

	GetPassengers(ctx context.Context) ([]entities.Passenger, error)
}

type DocumentUsecaser interface {
	CreateDocument(ctx context.Context, document entities.Document) (string, error)
	ReplaceDocument(ctx context.Context, document entities.Document) error
	DeleteDocument(ctx context.Context, id entities.Id) error
	GetDocumentsByPassengerId(ctx context.Context, id entities.Id) ([]entities.Document, error)
}

type ReportUsecaser interface {
	GetRowsByPassengerIdForPeriod(ctx context.Context, id entities.Id, filter entities.PeriodFilter) ([]entities.ReportRowByPassengerForPeriod, error)
}

type Logger interface {
	Debug(err error, format string, msg ...any)
	Error(err error, format string, msg ...any)
}
