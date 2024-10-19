package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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
			document.PassengerId,
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

func (r *Repository) ReplaceDocument(ctx context.Context, document entities.Document) error {
	sql, args, err := r.Builder.Update("documents").
		SetMap(squirrel.Eq{
			"type":         document.Type,
			"number":       document.Number,
			"passenger_id": document.PassengerId,
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

// TODO: exclude PassengerId from set
func (r *Repository) GetDocumentsByPassengerId(ctx context.Context, id entities.Id) ([]entities.Document, error) {
	sql, args, err := r.Builder.Select(
		"document_id",
		"type",
		"number",
		"passenger_id",
	).
		From("documents").
		Where(squirrel.Eq{
			"passenger_id": id.Value,
		}).
		ToSql()
	if err != nil {
		return []entities.Document{}, fmt.Errorf("repository: document: GetDocumentsByPassengerId: Select: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return []entities.Document{}, fmt.Errorf("repository: document: GetDocumentsByPassengerId: Query: %w", err)
	}

	documents := []entities.Document{}
	document := entities.Document{}

	_, err = pgx.ForEachRow(
		rows,
		[]any{
			&document.Id,
			&document.Type,
			&document.Number,
			&document.PassengerId,
		},
		func() error {
			documents = append(documents, document)
			return nil
		},
	)
	if err != nil {
		return []entities.Document{}, fmt.Errorf("repository: document: GetDocumentsByPassengerId: ForEachRow: %w", err)
	}

	if len(documents) == 0 {
		return []entities.Document{}, fmt.Errorf("repository: document: GetDocumentsByPassengerId: len: %w", entities.ErrorNothingFound)
	}

	return documents, nil
}
