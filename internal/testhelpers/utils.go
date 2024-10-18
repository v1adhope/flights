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

func (u *Utils) GetFirstTicketID(ctx context.Context) string {
	sql, args, err := u.Builder.Select("ticket_id").
		From("tickets").
		Limit(1).
		ToSql()
	if err != nil {
		log.Fatalf("testhelpers: utils: GetFirstTicketID: Select: %v", err)
	}

	id := ""

	if err := u.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		log.Fatalf("testhelpers: utils: GetFirstTicketID: QueryRow: %v", err)
	}

	return id
}
