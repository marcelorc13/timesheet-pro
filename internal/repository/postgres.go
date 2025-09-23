// Package repository represents the repository layer, that have direct the access to the database
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connString string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}
	return pool
}
