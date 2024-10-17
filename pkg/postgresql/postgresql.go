package postgresql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Driver struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func Build(ctx context.Context, opts ...Option) (*Driver, error) {
	cfg := config(opts...)

	pool, err := pgxpool.New(ctx, cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("postgresql: postgresql: Build: New: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgresql: postgresql: Build: Ping: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Driver{pool, builder}, nil
}

func (p *Driver) Close() {
	p.Pool.Close()
}
