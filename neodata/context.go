package neodata

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/infrastructure/auth/policy"
	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"go.uber.org/zap"
)

// Context encapsulates shared dependencies for the microservice
type NeoCtx struct {
	HTTPServer    *fiber.App
	PolicyManager *policy.PolicyManager
	NATS          *messaging.NATSClient
	Logger        *zap.Logger
	DB            *pgxpool.Pool
	Services      *ServiceRegistry // Add a dynamic service registry
}

// NewContext initializes a new Neo Context
// Components can be nil if not used by the microservice.
func NewContext(
	httpServer *fiber.App,
	policyManager *policy.PolicyManager,
	natsClient *messaging.NATSClient,
	logger *zap.Logger,
	db *pgxpool.Pool,
) *NeoCtx {
	return &NeoCtx{
		Logger:        logger,
		DB:            db,
		NATS:          natsClient,
		HTTPServer:    httpServer,
		PolicyManager: policyManager,
		Services:      &ServiceRegistry{}, // Initialize the service registry
	}
}
