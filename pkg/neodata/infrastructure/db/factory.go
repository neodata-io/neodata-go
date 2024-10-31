package db

import (
	"context"
	"fmt"

	"github.com/neodata-io/neodata-go/pkg/neodata/config"
	"github.com/neodata-io/neodata-go/pkg/neodata/infrastructure/db/postgres"
	"github.com/neodata-io/neodata-go/pkg/neodata/interfaces"
)

// Register creates a new database instance based on the provided configuration.
func Register(ctx context.Context, cfg *config.DatabaseConfig) (interfaces.Database, error) {
	switch cfg.Type {
	case "postgres":
		return postgres.New(ctx, cfg)
	// Add cases for other databases here...
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}
