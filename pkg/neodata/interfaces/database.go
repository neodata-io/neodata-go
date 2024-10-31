// interfaces/database.go
package interfaces

import (
	"context"

	"github.com/neodata-io/neodata-go/pkg/neodata/config"
)

// Database defines the methods required for any database implementation.
type Database interface {
	// Add methods common to all databases, e.g., for querying or managing connections
	New(ctx context.Context, cfg *config.DatabaseConfig)
}
