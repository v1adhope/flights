package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/v1adhope/flights/internal/entities"
)

func (u *Usecases) Create(ctx context.Context, ticket entities.Ticket) (string, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return "", err
	}

	ticket.Id = id.String()

	ticket.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := u.repos.Create(ctx, ticket); err != nil {
		return "", err
	}

	return ticket.Id, nil
}

func (u *Usecases) Replace(ctx context.Context, ticket entities.Ticket) error {
	if err := u.repos.Replace(ctx, ticket); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) Delete(ctx context.Context, id string) error {
	if err := u.repos.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) GetAll(ctx context.Context) ([]entities.Ticket, error) {
	tickets, err := u.repos.GetAll(ctx)
	if err != nil {
		return []entities.Ticket{}, err
	}

	return tickets, nil
}
