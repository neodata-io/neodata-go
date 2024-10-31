package policy

import (
	_ "embed"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/neodata-io/neodata-go/config"
)

func newModel() model.Model {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act, eft")
	m.AddDef("e", "e", "e = some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && r.obj == p.obj && r.act == p.act")

	return m
}

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

	// Connect to PostgreSQL as the adapter for Casbin
	adapter, err := xormadapter.NewAdapter("postgres", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Casbin adapter: %v", err)
	}

	// Load Casbin model and policy from configuration file
	enforcer, err := casbin.NewEnforcer(newModel(), adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %v", err)
	}

	// Load policies from the database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load Casbin policies: %v", err)
	}

	return enforcer, nil
}
