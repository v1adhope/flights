package usecases

import (
	"context"

	"github.com/v1adhope/flights/internal/entities"
)

func (u *Usecases) GetRowsByPassengerIdForPeriod(ctx context.Context, id entities.Id, filter entities.PeriodFilter) ([]entities.ReportRowByPassengerForPeriod, error) {
	rows, err := u.repos.GetRowsByPassengerIdForPeriod(ctx, id, filter)
	if err != nil {
		return []entities.ReportRowByPassengerForPeriod{}, err
	}

	return rows, nil
}
