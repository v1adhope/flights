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

type passengerTicketWholeInfoDto struct {
	Id         *string
	FirstName  *string
	LastName   *string
	MiddleName *string
}

func (d *passengerTicketWholeInfoDto) toEntity() entities.Passenger {
	if d.Id == nil {
		return entities.Passenger{}
	}
	return entities.Passenger{
		Id:         *d.Id,
		FirstName:  *d.FirstName,
		LastName:   *d.LastName,
		MiddleName: *d.MiddleName,
	}
}

type documentTicketWholeInfoDto struct {
	Id     *string
	Type   *string
	Number *string
}

func (d *documentTicketWholeInfoDto) toEntity() entities.DocumentTicketWholeInfo {
	if d.Id == nil {
		return entities.DocumentTicketWholeInfo{}
	}

	return entities.DocumentTicketWholeInfo{
		Id:     *d.Id,
		Type:   *d.Type,
		Number: *d.Number,
	}
}
