package authorization

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
)

// InitializeCasbin creates and returns a new Casbin enforcer with a PostgreSQL adapter.
func InitializeCasbin(dbUser string, dbPassword string, dbHost string, dbPort string, dbName string) (*casbin.Enforcer, error) {
	// Validate mandatory environment variables
	if dbUser == "" {
		return nil, fmt.Errorf("missing environment variable: DB_USER")
	}
	if dbHost == "" {
		return nil, fmt.Errorf("missing environment variable: DB_HOST")
	}
	if dbPort == "" {
		return nil, fmt.Errorf("missing environment variable: DB_PORT")
	}
	if dbName == "" {
		return nil, fmt.Errorf("missing environment variable: DB_NAME")
	}
	// Construct the database URL
	databaseUrl := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser,
		dbPassword,
		dbName,
		dbHost,
		dbPort,
	)

	// Connect to PostgreSQL as the adapter for Casbin
	adapter, err := xormadapter.NewAdapter("postgres", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Casbin adapter: %v", err)
	}

	// Load Casbin model and policy from configuration file
	enforcer, err := casbin.NewEnforcer("/app/configs/casbin/rbac_model.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %v", err)
	}

	// Load policies from the database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load Casbin policies: %v", err)
	}

	return enforcer, nil
}
