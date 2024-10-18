// db/postgres.go
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/config"
)

// NewPool initializes a PostgreSQL connection pool with given parameters.
func NewPool(cfg *config.AppConfig) (*pgxpool.Pool, error) {
	// Validate mandatory environment variables
	if cfg.Database.User == "" {
		return nil, fmt.Errorf("missing database user in configuration")
	}
	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("missing database host in configuration")
	}
	if cfg.Database.Port == 0 {
		return nil, fmt.Errorf("missing database port in configuration")
	}
	if cfg.Database.Name == "" {
		return nil, fmt.Errorf("missing database name in configuration")
	}

	// Construct the database URL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// Create a context with timeout for connecting to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the PostgreSQL database using pgxpool
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to PostgreSQL: %v", err)
	}

	// Set optional connection pool settings
	dbPool.Config().MaxConns = 10                 // Maximum number of connections
	dbPool.Config().MinConns = 2                  // Minimum number of idle connections
	dbPool.Config().MaxConnLifetime = time.Hour   // Max lifetime of a connection
	dbPool.Config().MaxConnIdleTime = time.Minute // Max idle time of a connection

	var result int
	err = dbPool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		log.Fatalf("Failed to validate PostgreSQL connection: %v", err)
	}

	return dbPool, nil
}
