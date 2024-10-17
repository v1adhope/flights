package usecases

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

type Reposer interface {
	Create(ctx context.Context, ticket entities.Ticket) error
}
