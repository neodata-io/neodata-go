package neodata

import (
	"context"
	"fmt"

	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/infrastructure/auth/policy"
	"github.com/neodata-io/neodata-go/infrastructure/db/postgres"
	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"github.com/neodata-io/neodata-go/infrastructure/transport/http"
	"github.com/neodata-io/neodata-go/logger"

	"go.uber.org/zap"
)

// App is the main struct that manages all components and services of the application.
//
//	handles the lifecycle and high-level initialization of the app.
type App struct {
	ConfigManager config.ConfigManager // Core config accessed from App
	Logger        *zap.Logger          // Core logger accessed from App
	Context       *NeoCtx              // Scoped services within NeoCtx
}

// Option defines a function that modifies the Neo Context
type Option func(*NeoCtx) error

// New initializes the application with options for services
func New(options ...Option) (*App, error) {
	// Step 1: Load Configuration
	cfgManager, err := config.NewConfigManager("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %w", err)
	}

	// Step 2: Initialize the default logger (ensures a logger is always available)
	log, err := logger.InitServiceLogger(cfgManager.GetAppConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize default logger: %w", err)
	}

	// Step 3: Create a base context with the initialized logger
	neoCtx, err := newContext(context.Background(), log, cfgManager)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		if err := option(neoCtx); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &App{
		Context:       neoCtx,
		Logger:        log,
		ConfigManager: cfgManager,
	}, nil
}

// WithLogger configures the application logger.
func WithLogger() Option {
	return func(ctx *NeoCtx) error {
		log, err := logger.InitServiceLogger(ctx.Config.GetAppConfig())
		if err != nil {
			return fmt.Errorf("failed to initialize logger: %v", err)
		}
		ctx.Logger = log
		return nil
	}
}

// WithPostgres configures a PostgreSQL pool.
func WithPostgres() Option {
	return func(ctx *NeoCtx) error {
		pool, err := postgres.NewPool(ctx.Context, ctx.Config.GetAppConfig())
		if err != nil {
			return fmt.Errorf("failed to initialize PostgreSQL: %v", err)
		}
		ctx.db = pool
		return nil
	}
}

// WithNATS configures a NATS client.
func WithNATS() Option {
	return func(ctx *NeoCtx) error {
		if ctx.messaging != nil {
			return nil
		}
		natsClient, err := messaging.NewNATSClient(ctx.Context, ctx.Config.GetAppConfig().Messaging.PubsubBroker)
		if err != nil {
			return fmt.Errorf("error creating NATS client: %v", err)
		}
		ctx.messaging = messaging.NewPublisher(natsClient, 0, 0)
		return nil
	}
}

// WithPolicyManager configures a Policy Manager.
func WithPolicyManager() Option {
	return func(ctx *NeoCtx) error {
		policyManager, err := policy.NewPolicyManager(ctx.Config.GetAppConfig())
		if err != nil {
			return fmt.Errorf("failed to initialize Policy Manager: %v", err)
		}
		ctx.policyManager = policyManager
		return nil
	}
}

// WithHTTPServer configures an HTTP server.
func WithHTTPServer() Option {
	return func(ctx *NeoCtx) error {
		if ctx.Logger == nil {
			return fmt.Errorf("logger is required but not initialized")
		}
		ctx.httpServer = http.NewHTTPServer(ctx.Config.GetAppConfig(), ctx.Logger)
		ctx.Logger.Info("HTTP server initialized with Fiber")
		return nil
	}
}

// Shutdown gracefully shuts down the app's services
func (a *App) Shutdown(ctx context.Context) error {
	if db, err := a.Context.GetDB(); err == nil {
		db.Close()
	}

	if server, err := a.Context.GetHTTPServer(); err == nil {
		return server.ShutdownWithContext(ctx)
	}

	return nil
}
