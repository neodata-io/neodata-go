package neodata

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/infrastructure/auth/policy"
	"github.com/neodata-io/neodata-go/infrastructure/db/postgres"
	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"github.com/neodata-io/neodata-go/infrastructure/transport/http"
	"github.com/neodata-io/neodata-go/logger"

	"go.uber.org/zap"
)

// App struct will hold the context and allow easier initialization
type App struct {
	Context *NeoCtx
}

// Option defines a function that modifies the Neo Context
type Option func(*NeoCtx) error

// New initializes the application with options for services
func New(opts ...Option) (*App, error) {
	ctx := context.Background()
	// Step 1: Load Configuration
	cfgManager, err := config.NewConfigManager("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %v", err)
	}

	// Create a base context with default components
	neoCtx := NewContext(ctx, nil, nil, nil, nil, cfgManager, nil)

	// Apply options (e.g., Logger, DB, NATS, etc.)
	for _, opt := range opts {
		opt(neoCtx)
	}

	return &App{
		Context: neoCtx,
	}, nil
}

// WithLogger allows the user to inject a logger
func WithLogger(log ...*zap.Logger) Option {
	return func(ctx *NeoCtx) error {
		appConfig := ctx.Config.GetAppConfig()
		log, err := logger.InitServiceLogger(appConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize logger: %v", err)
		}
		ctx.Logger = log
		return nil
	}
}

// WithPostgres allows the user to inject a PostgreSQL pool
func WithPostgres(pool ...*pgxpool.Pool) Option {
	return func(ctx *NeoCtx) error {
		appConfig := ctx.Config.GetAppConfig()
		pool, err := postgres.NewPool(ctx.Context, appConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize PostgreSQL: %v", err)
		}

		ctx.DB = pool
		return nil
	}
}

// WithNATS allows the user to inject a NATS client
func WithNATS(client ...*messaging.NATSClient) Option {
	return func(ctx *NeoCtx) error {
		// Access AppConfig via GetAppConfig method from ConfigManager
		appConfig := ctx.Config.GetAppConfig()
		// Check if the NATS client is already set; if not, create a new one
		if ctx.Messaging == nil {
			natsClient, err := messaging.NewNATSClient(ctx.Context, appConfig.Messaging.PubsubBroker)
			if err != nil {
				return fmt.Errorf("error creating NATS client: %v", err)
			}
			publisher := messaging.NewPublisher(natsClient, 0, 0)

			// Assign the publisher to ctx.Messaging
			ctx.Messaging = publisher
		}
		return nil
	}
}

// WithPolicyManager allows the user to inject a Policy Manager
func WithPolicyManager(policyManger ...*policy.PolicyManager) Option {
	return func(ctx *NeoCtx) error {
		appConfig := ctx.Config.GetAppConfig()
		policyManager, err := policy.NewPolicyManager(appConfig)
		if err != nil {
			if ctx.Logger != nil {
				ctx.Logger.Fatal("Failed to initialize Policy Manager", zap.Error(err))
			}
			return fmt.Errorf("failed to initialize Policy Manager: %v", err)

		}
		ctx.PolicyManager = policyManager
		return nil
	}
}

// WithHTTPServer allows the user to inject an HTTP server
func WithHTTPServer(httpServer ...*fiber.App) Option {
	return func(ctx *NeoCtx) error {
		appConfig := ctx.Config.GetAppConfig()
		fiberApp := http.NewHTTPServer(appConfig)
		if _, err := http.StartServer(fiberApp, appConfig); err != nil {
			ctx.Logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
		ctx.HTTPServer = fiberApp
		return nil
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

// GetPublisher retrieves the publisher from the context
func (a *App) GetPublisher() messaging.Messaging {
	if a.Context.Messaging == nil {
		a.Context.Logger.Warn("Publisher not initialized")
		return nil
	}
	return a.Context.Messaging
}

// GetHTTPServer retrieves the Fiber HTTP server from the context
func (a *App) GetHTTPServer() *fiber.App {
	return a.Context.HTTPServer
}

// GetPolicyManager retrieves the PolicyManager from the context
func (a *App) GetPolicyManager() *policy.PolicyManager {
	return a.Context.PolicyManager
}

// GetLogger retrieves the Logger from the context
func (a *App) GetLogger() *zap.Logger {
	return a.Context.Logger
}

// GetDB retrieves the PostgreSQL pool from the context
func (a *App) GetDB() *pgxpool.Pool {
	return a.Context.DB
}

// GetConfig retrieves the ConfigManager from the context
func (a *App) GetConfig() config.ConfigManager {
	return a.Context.Config
}

// GetServices retrieves the ServiceRegistry from the context
func (a *App) GetServices() *ServiceRegistry {
	return a.Context.Services
}

// GetBaseContext retrieves the base context
func (a *App) GetBaseContext() context.Context {
	return a.Context.Context
}
