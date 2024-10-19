package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/v1adhope/flights/internal/entities"
)

func catchExpectedDocumentCreationError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == "uq_documents_type_number" {
			return fmt.Errorf("repository: document: catchExpectedErrorDocumentCreation: %w", entities.ErrorHasAlreadyExists)
		}

		if pgErr.ConstraintName == "fk_document_passenger_passenger_id" {
			return fmt.Errorf("repository: document: catchExpectedErrorDocumentCreation: %w", entities.ErrorPassengerDoesNotExists)
		}
	}

	return nil
}
