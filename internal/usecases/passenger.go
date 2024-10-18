package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/v1adhope/flights/internal/entities"
)

func (u *Usecases) CreatePassenger(ctx context.Context, passenger entities.Passenger) error {
	id, err := uuid.NewV6()
	if err != nil {
		return fmt.Errorf("usecases: passenger: CreatePassenger: NewV6: %w", err)
	}

	passenger.Id = id.String()

	if err := u.repos.CreatePassenger(ctx, passenger); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) ReplacePassenger(ctx context.Context, passenger entities.Passenger) error {
	if err := u.repos.ReplacePassenger(ctx, passenger); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) DeletePassenger(ctx context.Context, id string) error {
	if err := u.repos.DeletePassenger(ctx, id); err != nil {
		return err
	}

	return nil
}
