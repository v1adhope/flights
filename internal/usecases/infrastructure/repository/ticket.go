package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/v1adhope/flights/internal/entities"
)

func (r *Repository) CreateTicket(ctx context.Context, ticket entities.Ticket) error {
	sql, args, err := r.Builder.Insert("tickets").
		Columns(
			"ticket_id",
			"provider",
			"fly_from",
			"fly_to",
			"fly_at",
			"arrive_at",
			"created_at",
		).
		Values(
			ticket.Id,
			ticket.Provider,
			ticket.FlyFrom,
			ticket.FlyTo,
			ticket.FlyAt,
			ticket.ArriveAt,
			ticket.CreatedAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: ticket: Replace: Insert: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("repository: ticket: Replace: Exec: %w", err)
	}

	return nil
}

func (r *Repository) ReplaceTicket(ctx context.Context, ticket entities.Ticket) error {
	sql, args, err := r.Builder.Update("tickets").
		SetMap(squirrel.Eq{
			"provider":  ticket.Provider,
			"fly_from":  ticket.FlyFrom,
			"fly_to":    ticket.FlyTo,
			"fly_at":    ticket.FlyAt,
			"arrive_at": ticket.ArriveAt,
		}).
		Where(squirrel.Eq{
			"ticket_id": ticket.Id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: ticket: Replace: Update: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repository: ticket: Replace: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: ticket: Replace: RowsAffected: %w", entities.ErrorNothingToChange)
	}

	return nil
}

func (r *Repository) DeleteTicket(ctx context.Context, id string) error {
	sql, args, err := r.Builder.Delete("tickets").
		Where(squirrel.Eq{
			"ticket_id": id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: ticket: Delete: Delete: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repository: ticket: Delete: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: ticket: Delete: RowsAffected: %w", entities.ErrorNothingToDelete)
	}

	return nil
}

func (r *Repository) GetAllTickets(ctx context.Context) ([]entities.Ticket, error) {
	sql, args, err := r.Builder.Select(
		"ticket_id",
		"provider",
		"fly_from",
		"fly_to",
		"fly_at",
		"arrive_at",
		"created_at",
	).
		From("tickets").
		ToSql()
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("repository: ticket: GetAll: Select: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)

	tickets := []entities.Ticket{}
	ticket := ticketDto{}

	_, err = pgx.ForEachRow(rows, []any{
		&ticket.Id,
		&ticket.Provider,
		&ticket.FlyFrom,
		&ticket.FlyTo,
		&ticket.FlyAt,
		&ticket.ArriveAt,
		&ticket.CreatedAt,
	}, func() error {
		tickets = append(tickets, ticket.toEntity())
		return nil
	})
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("repository: ticket: GetAll: ForEachRow: %w", err)
	}

	if len(tickets) == 0 {
		return []entities.Ticket{}, fmt.Errorf("repository: ticket: GetAll: len: %w", entities.ErrorNothingFound)
	}

	return tickets, nil
}
