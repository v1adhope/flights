package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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
	sql, args, err := r.Builder.Update("passengers").
		SetMap(squirrel.Eq{
			"first_name":  passenger.FirstName,
			"last_name":   passenger.LastName,
			"middle_name": passenger.MiddleName,
		}).
		Where(squirrel.Eq{
			"passenger_id": passenger.Id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: passenger: ReplacePassenger: Update: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repository: passenger: ReplacePassenger: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: passenger: ReplacePassenger: RowsAffected: %w", entities.ErrorNothingToChange)
	}

	return nil
}

func (r *Repository) DeletePassenger(ctx context.Context, id string) error {
	sql, args, err := r.Builder.Delete("passengers").
		Where(squirrel.Eq{
			"passenger_id": id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: passenger: DeletePassenger: Delete: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repository: passenger: DeletePassenger: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: ticket: DeletePassenger: RowsAffected: %w", entities.ErrorNothingToDelete)
	}

	return nil
}

func (r *Repository) GetAllPassengers(ctx context.Context) ([]entities.Passenger, error) {
	sql, args, err := r.Builder.Select(
		"passenger_id",
		"first_name",
		"last_name",
		"middle_name",
	).
		From("passengers").
		ToSql()
	if err != nil {
		return []entities.Passenger{}, fmt.Errorf("repository: passenger: GetAllPassengers: Select: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return []entities.Passenger{}, fmt.Errorf("repository: passenger: GetAllPassengers: Query: %w", err)
	}

	passengers := []entities.Passenger{}
	passenger := entities.Passenger{}

	_, err = pgx.ForEachRow(
		rows,
		[]any{
			&passenger.Id,
			&passenger.FirstName,
			&passenger.LastName,
			&passenger.MiddleName,
		},
		func() error {
			passengers = append(passengers, passenger)
			return nil
		},
	)
	if err != nil {
		return []entities.Passenger{}, fmt.Errorf("repository: passenger: GetAllPassengers: ForEachRow: %w", err)
	}

	if len(passengers) == 0 {
		return []entities.Passenger{}, fmt.Errorf("repository: passenger: GetAllPassengers: len: %w", entities.ErrorNothingFound)
	}

	return passengers, nil
}
