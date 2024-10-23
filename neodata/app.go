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

// App struct will hold the context and allow easier initialization
type App struct {
	Context *NeoCtx
}

// Option defines a function that modifies the Neo Context
type Option func(*NeoCtx)

// New initializes the application with options for services
func New(ctx context.Context, configPath string, opts ...Option) (*App, error) {
	// Step 1: Load Configuration
	cfgManager, err := config.NewConfigManager(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %v", err)
	}

	// Create a base context with default components
	neoCtx := NewContext(nil, nil, nil, nil, cfgManager, nil)

	// Apply options (e.g., Logger, DB, NATS, etc.)
	for _, opt := range opts {
		opt(neoCtx)
	}

	return &App{
		Context: neoCtx,
	}, nil
}

// WithLogger allows the user to inject a logger
func WithLogger(logger *zap.Logger) Option {
	return func(ctx *NeoCtx) {
		ctx.Logger = logger
	}
}

// WithPostgres allows the user to inject a PostgreSQL pool
func WithPostgres(dbPool *pgxpool.Pool) Option {
	return func(ctx *NeoCtx) {
		ctx.DB = dbPool
	}
}

// WithNATS allows the user to inject a NATS client
func WithNATS() Option {
	return func(ctx *NeoCtx) {
		// Access AppConfig via GetAppConfig method from ConfigManager
		appConfig := ctx.Config.GetAppConfig()
		// Check if the NATS client is already set; if not, create a new one
		if ctx.Messaging == nil {
			natsClient, err := messaging.NewNATSClient(context.Background(), appConfig.Messaging.PubsubBroker)
			if err != nil {
				fmt.Printf("Error creating NATS client: %v\n", err)
				return
			}
			publisher := messaging.NewPublisher(natsClient, 0, 0)

			// Assign the publisher to ctx.Messaging
			ctx.Messaging = publisher
		}
	}
}

// WithPolicyManager allows the user to inject a Policy Manager
func WithPolicyManager(policyManager *policy.PolicyManager) Option {
	return func(ctx *NeoCtx) {
		ctx.PolicyManager = policyManager
	}
}

// WithHTTPServer allows the user to inject an HTTP server
func WithHTTPServer(httpServer *fiber.App) Option {
	return func(ctx *NeoCtx) {
		ctx.HTTPServer = httpServer
	}
}

// Shutdown gracefully shuts down the app's services
func (a *App) Shutdown(ctx context.Context) error {
	if a.Context.DB != nil {
		a.Context.DB.Close()
	}

	if a.Context.HTTPServer != nil {
		if err := a.Context.HTTPServer.ShutdownWithContext(ctx); err != nil {
			return err
		}
	}
	return nil
}
