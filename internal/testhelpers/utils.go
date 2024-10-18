package testhelpers

import (
	"context"
	"log"

	"github.com/v1adhope/flights/pkg/postgresql"
)

type Utils struct {
	*postgresql.Driver
}

func NewUtils(d *postgresql.Driver) *Utils {
	return &Utils{d}
}

func (u *Utils) GetTicketByOffset(ctx context.Context, offset uint64) string {
	sql, args, err := u.Builder.Select("ticket_id").
		From("tickets").
		Offset(offset).
		Limit(1).
		ToSql()
	if err != nil {
		log.Fatalf("testhelpers: utils: GetTicketByOffset: Select: %v", err)
	}

	id := ""

	if err := u.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		log.Fatalf("testhelpers: utils: GetTicketByOffset: QueryRow: %v", err)
	}

	return id
}

func (u *Utils) GetPassengerByOffset(ctx context.Context, offset uint64) string {
	sql, args, err := u.Builder.Select("passenger_id").
		From("passengers").
		Offset(offset).
		Limit(1).
		ToSql()
	if err != nil {
		log.Fatalf("testhelpers: utils: GetPassengerByOffset: Select: %v", err)
	}

	id := ""

	if err := u.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		log.Fatalf("testhelpers: utils: GetPassengerByOffset: QueryRow: %v", err)
	}

	return id
}
