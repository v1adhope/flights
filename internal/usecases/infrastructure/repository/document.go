package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/v1adhope/flights/internal/entities"
)

func (r *Repository) CreateDocument(ctx context.Context, document entities.Document) error {
	sql, args, err := r.Builder.Insert("documents").
		Columns(
			"document_id",
			"type",
			"number",
			"passenger_id",
		).
		Values(
			document.Id,
			document.Type,
			document.Number,
			document.PassangerId,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: document: CreateDocument: Insert: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if err := catchExpectedDocumentCreationError(err); err != nil {
			return err
		}

		return fmt.Errorf("repository: document: CreateDocument: Exec: %w", err)
	}

	return nil
}

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

func (r *Repository) ReplaceDocument(ctx context.Context, document entities.Document) error {
	sql, args, err := r.Builder.Update("documents").
		SetMap(squirrel.Eq{
			"type":         document.Type,
			"number":       document.Number,
			"passenger_id": document.PassangerId,
		}).
		Where(squirrel.Eq{
			"document_id": document.Id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: document: ReplaceDocument: Update: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if err := catchExpectedDocumentCreationError(err); err != nil {
			return err
		}

		return fmt.Errorf("repository: document: ReplaceDocument: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: document: ReplaceDocument: RowsAffected: %w", entities.ErrorNothingToChange)
	}

	return nil
}

func (r *Repository) DeleteDocument(ctx context.Context, id entities.Id) error {
	sql, args, err := r.Builder.Delete("documents").
		Where(squirrel.Eq{
			"document_id": id.Value,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository: document: DeleteDocument: Delete: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repository: document: DeleteDocument: Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repository: document: DeleteDocument: RowsAffected: %w", entities.ErrorNothingToDelete)
	}

	return nil
}
