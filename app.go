package app

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/infrastructure/auth/policy"
	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"github.com/neodata-io/neodata-go/neodata"
	"go.uber.org/zap"
)

// App struct will hold the context and allow easier initialization
type App struct {
	Context *neodata.Context
}

// Option defines a function that modifies the Neo Context
type Option func(*neodata.Context)

// New initializes the application with options for services
func New(ctx context.Context, cfg *config.AppConfig, opts ...Option) (*App, error) {
	// Create a base context with default components
	neoCtx := neodata.NewContext(nil, nil, nil, nil, nil)

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
	return func(ctx *neodata.Context) {
		ctx.Logger = logger
	}
}

// WithPostgres allows the user to inject a PostgreSQL pool
func WithPostgres(dbPool *pgxpool.Pool) Option {
	return func(ctx *neodata.Context) {
		ctx.DB = dbPool
	}
}

// WithNATS allows the user to inject a NATS client
func WithNATS(natsClient *messaging.NATSClient) Option {
	return func(ctx *neodata.Context) {
		ctx.NATS = natsClient
	}
}

// WithPolicyManager allows the user to inject a Policy Manager
func WithPolicyManager(policyManager *policy.PolicyManager) Option {
	return func(ctx *neodata.Context) {
		ctx.PolicyManager = policyManager
	}
}

// WithHTTPServer allows the user to inject an HTTP server
func WithHTTPServer(httpServer *fiber.App) Option {
	return func(ctx *neodata.Context) {
		ctx.HTTPServer = httpServer
	}
}

// Shutdown gracefully shuts down the app's services
func (a *App) Shutdown(ctx context.Context) error {
	if a.Context.DB != nil {
		a.Context.DB.Close()
	}
	if a.Context.NATS != nil {
		a.Context.NATS.Close()
	}
	if a.Context.HTTPServer != nil {
		if err := a.Context.HTTPServer.ShutdownWithContext(ctx); err != nil {
			return err
		}
	}
	return nil
}
