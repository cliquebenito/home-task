package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustNewPool(ctx context.Context, connStr string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}
	return pool
}
