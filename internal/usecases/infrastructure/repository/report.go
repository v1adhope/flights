package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/v1adhope/flights/internal/entities"
)

func (r *Repository) GetRowsByPassengerIdForPeriod(ctx context.Context, id entities.Id, filter entities.PeriodFilter) ([]entities.ReportRowByPassengerForPeriod, error) {
	sql, args, err := r.getBuilderKindOfGetRowsByPassengerForPeriod(id, filter, true).
		Suffix("union").
		SuffixExpr(
			r.getBuilderKindOfGetRowsByPassengerForPeriod(id, filter, false),
		).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return []entities.ReportRowByPassengerForPeriod{}, fmt.Errorf("repository: report: GetRowsByPassengerIdForPeriod: Query: %w", err)
	}

	reportRows := []entities.ReportRowByPassengerForPeriod{}
	reportRowDto := reportRowByPassengerForPeriodDto{}

	tag, err := pgx.ForEachRow(
		rows,
		[]any{
			&reportRowDto.DateOfIssue,
			&reportRowDto.FlyAt,
			&reportRowDto.TicketId,
			&reportRowDto.FlyFrom,
			&reportRowDto.FlyTo,
			&reportRowDto.ServiceProvided,
		},
		func() error {
			reportRows = append(reportRows, reportRowDto.toEntity())
			return nil
		},
	)
	if err != nil {
		return []entities.ReportRowByPassengerForPeriod{}, fmt.Errorf("repository: report: GetRowsByPassengerIdForPeriod: ForEachRow: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return []entities.ReportRowByPassengerForPeriod{}, fmt.Errorf("repository: ticket: GetRowsByPassengerIdForPeriod: RowsAffected: %w", entities.ErrorNothingFound)
	}

	return reportRows, nil
}

func (r *Repository) getBuilderKindOfGetRowsByPassengerForPeriod(id entities.Id, filter entities.PeriodFilter, isProvided bool) squirrel.SelectBuilder {
	var whereArriveStatement squirrel.Sqlizer
	var selectServiceProvided string

	if isProvided {
		selectServiceProvided = "true as service_provided"
		whereArriveStatement = squirrel.And{
			squirrel.LtOrEq{
				"created_at": filter.To,
			},
			squirrel.LtOrEq{
				"arrive_at": filter.To,
			},
		}
	} else {
		selectServiceProvided = "false as service_provided"
		whereArriveStatement = squirrel.And{
			squirrel.Expr(
				"created_at between ? and ?", filter.From, filter.To,
			),
			squirrel.Gt{
				"arrive_at": filter.To,
			},
		}
	}

	return r.Builder.Select(
		"tickets.created_at as date_of_issue",
		"tickets.fly_at",
		"tickets.ticket_id",
		"tickets.fly_from",
		"tickets.fly_to",
		selectServiceProvided,
	).
		From("passenger_ticket").
		LeftJoin("tickets using(ticket_id)").
		Where(squirrel.And{
			squirrel.Eq{
				"passenger_id": id.Value,
			},
			whereArriveStatement,
		},
		)
}
