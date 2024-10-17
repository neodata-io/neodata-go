// db/postgres.go
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to create PostgreSQL connection pool: %v", err)
		return nil, err
	}

	return pool, nil
}
