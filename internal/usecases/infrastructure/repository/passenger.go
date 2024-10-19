package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *Repository) DeletePassenger(ctx context.Context, id entities.Id) error {
	sql, args, err := r.Builder.Delete("passengers").
		Where(squirrel.Eq{
			"passenger_id": id.Value,
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
		return fmt.Errorf("repository: passenger: DeletePassenger: RowsAffected: %w", entities.ErrorNothingToDelete)
	}

	return nil
}

func (r *Repository) GetPassengers(ctx context.Context) ([]entities.Passenger, error) {
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

	return passengersRowReader(rows)
}

func (r *Repository) BoundToTicket(ctx context.Context, id entities.Id, ticketId entities.Id) error {
	sql, args, err := r.Builder.Insert("passenger_ticket").
		Columns(
			"passenger_id",
			"ticket_id",
		).
		Values(
			id.Value,
			ticketId.Value,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: passenger: BoundToTicket: Insert: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "pk_ticket_passenger_ticket_id_passenger_id" {
				return fmt.Errorf("repository: document: BoundToTicket: Exec: %w", entities.ErrorHasAlreadyExists)
			}

			if pgErr.ConstraintName == "fk_ticket_passenger_passenger_passenger_id" {
				return fmt.Errorf("repository: document: BoundToTicket: Exec: %w", entities.ErrorPassengerDoesNotExists)
			}

			if pgErr.ConstraintName == "fk_ticket_passenger_tickets_ticket_id" {
				return fmt.Errorf("repository: document: BoundToTicket: Exec: %w", entities.ErrorTicketDoesNotExists)
			}
		}

		return fmt.Errorf("repository: passenger: BoundToTicket: Exec: %w", err)
	}

	return nil
}

func (r *Repository) UnboundToTicket(ctx context.Context, id entities.Id, ticketId entities.Id) error {
	sql, args, err := r.Builder.Delete("passenger_ticket").
		Where(squirrel.Eq{
			"passenger_id": id.Value,
			"ticket_id":    ticketId.Value,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: passenger: UnboundToTicket: Delete: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repository: passenger: UnboundToTicket: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: passenger: UnboundToTicket: RowsAffected: %w", entities.ErrorNothingToDelete)
	}

	return nil
}

func (r *Repository) GetPassengersByTicketId(ctx context.Context, id entities.Id) ([]entities.Passenger, error) {
	sql, args, err := r.Builder.Select(
		"passengers.passenger_id",
		"passengers.first_name",
		"passengers.last_name",
		"passengers.middle_name",
	).
		From("passenger_ticket").
		LeftJoin("passengers using(passenger_id)").
		Where(squirrel.Eq{
			"ticket_id": id.Value,
		}).
		ToSql()
	if err != nil {
		return []entities.Passenger{}, fmt.Errorf("repository: passenger: GetPassengersByTicketId: Select: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return []entities.Passenger{}, fmt.Errorf("repository: passenger: GetPassengersByTicketId: Query: %w", err)
	}

	return passengersRowReader(rows)
}
