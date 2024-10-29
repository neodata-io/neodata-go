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

// App manages the lifecycle and high-level initialization of the application.
type App struct {
	Config  *config.AppConfig // Application configuration
	Logger  *zap.Logger       // Application logger
	Context *NeoCtx           // Neo context for scoped services
}

// Option defines a function that modifies the Neo Context
type Option func(*NeoCtx) error

// New initializes the application with options for services
func New(options ...Option) (*App, error) {
	// Load Configuration
	cfgManager, err := config.NewConfigManager("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %w", err)
	}
	cfg := cfgManager.GetAppConfig()

	/// Initialize Logger
	log, err := logger.InitServiceLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize default logger: %w", err)
	}

	// Step 3: Create Base Context with Config and Logger references
	neoCtx, err := newContext(context.Background(), log, cfg)
	if err != nil {
		return nil, err
	}

	// Apply Options
	for _, option := range options {
		if err := option(neoCtx); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &App{
		Context: neoCtx,
		Logger:  log,
		Config:  cfgManager.GetAppConfig(),
	}, nil
}

// Run starts the application, including the HTTP server.
func (a *App) Run() error {
	a.Logger.Info("Starting application")
	httpServer, err := a.Context.GetHTTPServer()
	if err != nil {
		a.Logger.Error("HTTP server not configured", zap.Error(err))
		return fmt.Errorf("HTTP server not configured")
	}
	a.Logger.Info("Starting HTTP server")
	if _, err := http.StartServer(httpServer, a.Config); err != nil {
		a.Logger.Error("Failed to start HTTP server", zap.Error(err))
		return err
	}
	return nil
}

// WithPostgres configures a PostgreSQL pool.
func WithPostgres() Option {
	return func(ctx *NeoCtx) error {
		pool, err := postgres.NewPool(ctx.Context, ctx.Config)
		if err != nil {
			ctx.Logger.Error("Failed to initialize PostgreSQL", zap.Error(err))
			return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
		}
		ctx.db = pool
		ctx.Logger.Info("PostgreSQL connection pool initialized")
		return nil
	}
}

// WithNATS configures a NATS client.
func WithNATS() Option {
	return func(ctx *NeoCtx) error {
		if ctx.messaging != nil {
			ctx.Logger.Warn("Messaging client already configured, skipping NATS setup")
			return nil
		}
		natsClient, err := messaging.NewNATSClient(ctx.Context, ctx.Config.Messaging.PubsubBroker)
		if err != nil {
			ctx.Logger.Error("Failed to initialize NATS client", zap.Error(err))
			return fmt.Errorf("failed to initialize NATS client: %w", err)
		}
		ctx.messaging = messaging.NewPublisher(natsClient, 0, 0)
		ctx.Logger.Info("NATS messaging client initialized")
		return nil
	}
}

// WithPolicyManager configures a Policy Manager.
func WithPolicyManager() Option {
	return func(ctx *NeoCtx) error {
		policyManager, err := policy.NewPolicyManager(ctx.Config)
		if err != nil {
			ctx.Logger.Error("Failed to initialize Policy Manager", zap.Error(err))
			return fmt.Errorf("failed to initialize Policy Manager: %w", err)
		}
		ctx.policyManager = policyManager
		ctx.Logger.Info("Policy Manager initialized")
		return nil
	}
}

// WithHTTPServer configures an HTTP server.
func WithHTTPServer() Option {
	return func(ctx *NeoCtx) error {
		ctx.httpServer = http.NewHTTPServer(ctx.Config, ctx.Logger)
		ctx.Logger.Info("HTTP server initialized")
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

/* func (n *neodata.NeoCtx) StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		n.Logger.Info("Starting metrics server on :9090/metrics")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			n.Logger.Error("Metrics server failed", zap.Error(err))
		}
	}()
}
*/
