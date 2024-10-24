package neodata

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/infrastructure/auth/policy"
	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"go.uber.org/zap"
)

// Context encapsulates shared dependencies for the microservice
type NeoCtx struct {
	Context       context.Context // Add the base context here
	HTTPServer    *fiber.App
	PolicyManager *policy.PolicyManager
	Messaging     messaging.Messaging
	Logger        *zap.Logger
	DB            *pgxpool.Pool
	Config        config.ConfigManager // Store the interface, not a pointer to the interface
	Services      *ServiceRegistry     // Add a dynamic service registry
}

// NewContext initializes a new Neo Context
// Components can be nil if not used by the microservice.
func NewContext(
	context context.Context,
	httpServer *fiber.App,
	policyManager *policy.PolicyManager,
	messaging messaging.Messaging,
	logger *zap.Logger,
	configManager config.ConfigManager,
	db *pgxpool.Pool,
) *NeoCtx {
	return &NeoCtx{
		Context:       context,
		Logger:        logger,
		DB:            db,
		Messaging:     messaging,
		HTTPServer:    httpServer,
		PolicyManager: policyManager,
		Config:        configManager,
		Services:      &ServiceRegistry{}, // Initialize the service registry
	}
}
