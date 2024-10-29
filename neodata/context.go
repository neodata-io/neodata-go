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
	Config  *config.AppConfig
	Logger  *zap.Logger // Injected from the main application to enable structured logging

	db            *pgxpool.Pool
	httpServer    *fiber.App
	policyManager *policy.PolicyManager
	messaging     messaging.Messaging
	Services      *ServiceRegistry // Add a dynamic service registry
}

// NewContext initializes a new Neo Context
// Components can be nil if not used by the microservice.
func newContext(ctx context.Context, l *zap.Logger, cfg *config.AppConfig) (*NeoCtx, error) {
	return &NeoCtx{
		Context: ctx,
		Logger:  l,
		Config:  cfg,
	}, nil
}

// GetDB retrieves the database pool, logging an error if it is not configured.
func (n *NeoCtx) GetDB() (*pgxpool.Pool, error) {
	if n.db == nil {
		n.Logger.Error("Database not configured")
		return nil, fmt.Errorf("database not configured")
	}
	n.Logger.Info("Database retrieved successfully")
	return n.db, nil
}

// GetHTTPServer retrieves the HTTP server instance, logging an error if it is not configured.
func (n *NeoCtx) GetHTTPServer() (*fiber.App, error) {
	if n.httpServer == nil {
		n.Logger.Error("HTTP server not configured")
		return nil, fmt.Errorf("HTTP server not configured")
	}
	n.Logger.Info("HTTP server retrieved successfully")
	return n.httpServer, nil
}

// GetPolicyManager retrieves the policy manager, logging an error if it is not configured.
func (n *NeoCtx) GetPolicyManager() (*policy.PolicyManager, error) {
	if n.policyManager == nil {
		n.Logger.Error("Policy manager not configured")
		return nil, fmt.Errorf("policy manager not configured")
	}
	n.Logger.Info("Policy manager retrieved successfully")
	return n.policyManager, nil
}

// GetPublisher retrieves the messaging publisher, logging an error if it is not configured.
func (n *NeoCtx) GetPublisher() (messaging.Messaging, error) {
	if n.messaging == nil {
		n.Logger.Error("Messaging client not configured")
		return nil, fmt.Errorf("messaging client not configured")
	}
	n.Logger.Info("Messaging publisher retrieved successfully")
	return n.messaging, nil
}

// GetSubscriber retrieves the messaging subscriber, logging an error if it is not configured.
func (n *NeoCtx) GetSubscriber() (messaging.Messaging, error) {
	if n.messaging == nil {
		n.Logger.Error("Messaging client not configured")
		return nil, fmt.Errorf("messaging client not configured")
	}
	n.Logger.Info("Messaging subscriber retrieved successfully")
	return n.messaging, nil
}
