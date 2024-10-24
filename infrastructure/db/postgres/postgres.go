// db/postgres.go
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/config"
)

// NewPool initializes a PostgreSQL connection pool with given parameters.
func NewPool(ctx context.Context, cfg *config.AppConfig) (*pgxpool.Pool, error) {
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

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Connect to the PostgreSQL database using pgxpool
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to PostgreSQL: %v", err)
	}

	// Set optional connection pool settings
	dbPool.Config().MaxConns = 10                 // Maximum number of connections
	dbPool.Config().MinConns = 2                  // Minimum number of idle connections
	dbPool.Config().MaxConnLifetime = time.Hour   // Max lifetime of a connection
	dbPool.Config().MaxConnIdleTime = time.Minute // Max idle time of a connection

	return dbPool, nil
}
