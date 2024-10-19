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

func (u *Utils) getByOffset(ctx context.Context, offset uint64, table string, cols ...string) string {
	sql, args, err := u.Builder.Select(cols...).
		From(table).
		Offset(offset).
		Limit(1).
		ToSql()
	if err != nil {
		log.Fatalf("testhelpers: utils: getByOffset: Select: %v", err)
	}

	id := ""

	if err := u.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		log.Fatalf("testhelpers: utils: getByOffset: QueryRow: %v", err)
	}

	return id
}

func (u *Utils) GetTicketByOffset(ctx context.Context, offset uint64) string {
	return u.getByOffset(ctx, offset, "tickets", "ticket_id")
}

func (u *Utils) GetPassengerByOffset(ctx context.Context, offset uint64) string {
	return u.getByOffset(ctx, offset, "passengers", "passenger_id")
}

func (u *Utils) GetDocumentByOffset(ctx context.Context, offset uint64) string {
	return u.getByOffset(ctx, offset, "documents", "document_id")
}
