package v1

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type TicketUsecaser interface {
	Create(ctx context.Context, ticket entities.Ticket) (string, error)
}

type Logger interface {
	Debug(err error, format string, msg ...any)
	Error(err error, format string, msg ...any)
}
