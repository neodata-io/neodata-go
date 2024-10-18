package policy

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/neodata-io/neodata-go/config"
)

// InitializeCasbin creates and returns a new Casbin enforcer with a PostgreSQL adapter.
func InitializeCasbin(cfg *config.AppConfig) (*casbin.Enforcer, error) {

	databaseUrl := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	// Debugging: Log the database URL (without the password)
	fmt.Printf("Connecting to PostgreSQL at %s:%d (user: %s, dbname: %s)\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name)

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
