package repository

import "github.com/v1adhope/flights/pkg/postgresql"

type Repository struct {
	*postgresql.Driver
}

func New(d *postgresql.Driver) *Repository {
	return &Repository{d}
}
