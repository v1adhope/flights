package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/v1adhope/flights/internal/entities"
)

type ticketDto struct {
	Id        string
	Provider  string
	FlyFrom   string
	FlyTo     string
	FlyAt     pgtype.Timestamptz
	ArriveAt  pgtype.Timestamptz
	CreatedAt pgtype.Timestamptz
}

func (d *ticketDto) toEntity() entities.Ticket {
	return entities.Ticket{
		Id:        d.Id,
		Provider:  d.Provider,
		FlyFrom:   d.FlyFrom,
		FlyTo:     d.FlyTo,
		FlyAt:     d.FlyAt.Time.Format(time.RFC3339),
		ArriveAt:  d.ArriveAt.Time.Format(time.RFC3339),
		CreatedAt: d.CreatedAt.Time.Format(time.RFC3339),
	}
}
