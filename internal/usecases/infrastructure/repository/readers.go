package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/v1adhope/flights/internal/entities"
)

func passengersRowReader(rows pgx.Rows) ([]entities.Passenger, error) {
	passengers := []entities.Passenger{}
	passenger := entities.Passenger{}

	_, err := pgx.ForEachRow(
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
		return []entities.Passenger{}, fmt.Errorf("repository: readers: passengersRowReader: ForEachRow: %w", err)
	}

	if len(passengers) == 0 {
		return []entities.Passenger{}, fmt.Errorf("repository: readers: passengersRowReader: len: %w", entities.ErrorNothingFound)
	}

	return passengers, nil
}
