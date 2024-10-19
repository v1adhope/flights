package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/v1adhope/flights/internal/entities"
)

func (u *Usecases) CreateDocument(ctx context.Context, document entities.Document) (string, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return "", fmt.Errorf("usecases: document: CreateDocument: NewV6: %w", err)
	}

	document.Id = id.String()

	if err := u.repos.CreateDocument(ctx, document); err != nil {
		return "", err
	}

	return document.Id, nil
}

func (u *Usecases) ReplaceDocument(ctx context.Context, document entities.Document) error {
	if err := u.repos.ReplaceDocument(ctx, document); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) DeleteDocument(ctx context.Context, id entities.Id) error {
	if err := u.repos.DeleteDocument(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *Usecases) GetDocumentsByPassengerId(ctx context.Context, id entities.Id) ([]entities.Document, error) {
	documents, err := u.repos.GetDocumentsByPassengerId(ctx, id)
	if err != nil {
		return []entities.Document{}, err
	}

	return documents, nil
}
