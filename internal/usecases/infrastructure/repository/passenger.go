package repository

import (
	"context"
	"fmt"

	"github.com/v1adhope/flights/internal/entities"
)

func (r *Repository) CreatePassenger(ctx context.Context, passenger entities.Passenger) error {
	sql, args, err := r.Builder.Insert("passengers").
		Columns(
			"passenger_id",
			"first_name",
			"last_name",
			"middle_name",
		).
		Values(
			passenger.Id,
			passenger.FirstName,
			passenger.LastName,
			passenger.MiddleName,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: passenger: CreatePassenger: Insert: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("repository: passenger: CreatePassenger: Exec: %w", err)
	}

	return nil
}

func (r *Repository) ReplacePassenger(ctx context.Context, passenger entities.Passenger) error {
	return nil
}

func (r *Repository) DeletePassenger(ctx context.Context, id string) error {
	return nil
}
