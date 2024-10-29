package neodata

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/infrastructure/auth/policy"
	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"go.uber.org/zap"
)

// NeoCtx manages the application's dependencies, request, response, and context.
// encapsulates dependency management and provides a structured way to access shared services.
type NeoCtx struct {
	Context context.Context

	Logger *zap.Logger
	Config config.ConfigManager // Store the interface, not a pointer to the interface

	db            *pgxpool.Pool
	httpServer    *fiber.App
	policyManager *policy.PolicyManager
	messaging     messaging.Messaging
	services      *ServiceRegistry // Add a dynamic service registry
}

// NewContext initializes a new Neo Context
// Components can be nil if not used by the microservice.
func newContext(ctx context.Context, l *zap.Logger, cfgMngr config.ConfigManager) (*NeoCtx, error) {

	return &NeoCtx{
		Context: ctx,
		Logger:  l,
		Config:  cfgMngr,
	}, nil
}

// GetServiceRegistry initializes and returns the ServiceRegistry if it's not already set.
func (n *NeoCtx) getServiceRegistry() *ServiceRegistry {
	if n.services == nil {
		n.services = &ServiceRegistry{}
	}
	return n.services
}

// GetService retrieves a service by name from the ServiceRegistry within NeoCtx.
func (n *NeoCtx) GetService(name string) (interface{}, error) {
	// Ensure the ServiceRegistry is initialized
	serviceRegistry := n.getServiceRegistry()

	// Fetch the service by name
	service, exists := serviceRegistry.Get(name)
	if !exists {
		return nil, fmt.Errorf("service %s not found in registry", name)
	}

	return service, nil
}

// GetLogger returns the logger, defaulting to a no-op logger if none is set.
func (n *NeoCtx) GetLogger() *zap.Logger {
	return n.Logger
}

func (n *NeoCtx) GetDB() (*pgxpool.Pool, error) {
	if n.db == nil {
		return nil, fmt.Errorf("database not configured")
	}
	return n.db, nil
}

func (n *NeoCtx) GetHTTPServer() (*fiber.App, error) {
	if n.httpServer == nil {
		return nil, fmt.Errorf("HTTP server not configured")
	}
	return n.httpServer, nil
}

func (n *NeoCtx) GetPolicyManager() (*policy.PolicyManager, error) {
	if n.policyManager == nil {
		return nil, fmt.Errorf("policy manager not configured")
	}
	return n.policyManager, nil
}

func (n *NeoCtx) GetPublisher() (messaging.Messaging, error) {
	if n.messaging == nil {
		return nil, fmt.Errorf("messaging client not configured")
	}
	return n.messaging, nil
}

func (n *NeoCtx) GetSubscriber() (messaging.Messaging, error) {
	if n.messaging == nil {
		return nil, fmt.Errorf("messaging client not configured")
	}
	return n.messaging, nil
}
