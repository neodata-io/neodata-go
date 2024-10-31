package neodata

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/neodata-io/neodata-go/infrastructure/messaging"
	"github.com/neodata-io/neodata-go/infrastructure/transport/http"
	"github.com/neodata-io/neodata-go/pkg/neodata/config"

	"github.com/neodata-io/neodata-go/neodata/container"
	"github.com/neodata-io/neodata-go/neodata/interfaces"

	"go.uber.org/zap"
)

const (
	DatabaseConfigKey   = "database"
	MessagingConfigKey  = "messaging"
	AuthPolicyConfigKey = "auth.policy"
	AppPortConfigKey    = "app.port"
)

// / App manages the app lifecycle and optional services.
type App struct {
	Config    config.ConfigProvider // Application configuration
	container *container.Container  // container is unexported because this is an internal implementation and applications are provided access to it via Context

	db            *pgxpool.Pool
	httpServer    interfaces.HTTPServer
	policyManager interfaces.PolicyManager
	messaging     messaging.Messaging //TODO: Natsclient

}

// Option defines a function that configures optional services in App.
type Option func(*App) error

// New initializes the application with the provided configuration and options.
func New(options ...Option) (*App, error) {
	// Load Configuration
	cfgManager, err := config.NewConfigManager("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %w", err)
	}

	// Create shared Container (Logger, httpServer)
	container, err := container.NewContainer(cfgManager)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config:    cfgManager,
		container: container,
	}

	// Conditionally add WithPostgres if database config is enabled
	if cfgManager.IsEnabled(DatabaseConfigKey) {
		options = append(options, WithPostgres())
	}
	if cfgManager.IsEnabled(MessagingConfigKey) {
		options = append(options, WithNATS())
	}
	if cfgManager.IsEnabled(AuthPolicyConfigKey) {
		options = append(options, WithPolicyManager())
	}
	if cfgManager.IsEnabled(AppPortConfigKey) {
		options = append(options, WithHTTPServer())
	}

	// Apply Options
	for _, option := range options {
		if err := option(app); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return app, nil
}

// Run starts the application, including the HTTP server.
func (a *App) Run() error {
	a.Logger().Info(`
	███╗   ██╗███████╗ ██████╗ ██████╗  █████╗ ████████╗ █████╗ 
	████╗  ██║██╔════╝██╔════╝ ██╔══██╗██╔══██╗╚══██╔══╝██╔══██╗
	██╔██╗ ██║█████╗  ██║  ███╗██████╔╝███████║   ██║   ███████║
	██║╚██╗██║██╔══╝  ██║   ██║██╔═══╝ ██╔══██║   ██║   ██╔══██║
	██║ ╚████║███████╗╚██████╔╝██║     ██║  ██║   ██║   ██║  ██║
	╚═╝  ╚═══╝╚══════╝ ╚═════╝ ╚═╝     ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝
	
		Version: 1.0.0
		Author: Ruben Claes
	`)

	// Start the HTTP Server
	if a.Config.IsEnabled(AppPortConfigKey) {
		httpServer, _ := a.HTTPServer()
		if err := httpServer.Start(a.Config.AppPort()); err != nil {
			a.Logger().Error("Failed to start HTTP server", zap.Error(err))
			return err
		}
	}
	return nil
}

// WithDatabase configures a PostgreSQL pool.
func WithDatabase() Option {
	return func(a *App) error {
		dbConfig := a.Config.GetDatabaseConfig() // Get the database configuration

		// Create the database instance using the factory
		database, err := db.Reg(context.Background(), dbConfig)
		if err != nil {
			a.Logger().Error("Failed to initialize database", zap.Error(err))
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		a.db = database // Store the database instance
		a.Logger().Info("Database connection initialized")
		return nil

	}
}

// WithNATS configures a NATS client.
func WithNATS() Option {
	return func(a *App) error {
		if a.messaging != nil {
			a.Logger().Warn("Messaging client already configured, skipping NATS setup")
			return nil
		}
		natsClient, err := messaging.NewNATSClient(ctx.Context, a.Config.GetMessagingConfig().PubsubBroker)
		if err != nil {
			a.Logger().Error("Failed to initialize NATS client", zap.Error(err))
			return fmt.Errorf("failed to initialize NATS client: %w", err)
		}
		a.messaging = messaging.NewPublisher(natsClient, 0, 0)
		a.Logger().Info("NATS messaging client initialized")
		return nil
	}
}

// WithPolicyManager configures a Policy Manager.
func WithPolicyManager() Option {
	return func(a *App) error {
		policyManager, err := policy.CasbinAdapter(a.Config)
		if err != nil {
			a.Logger().Error("Failed to initialize Policy Manager", zap.Error(err))
			return fmt.Errorf("failed to initialize Policy Manager: %w", err)
		}
		a.policyManager = policyManager
		a.Logger().Info("Policy Manager initialized")
		return nil
	}
}

// WithHTTPServer configures an HTTP server.
func WithHTTPServer() Option {
	return func(a *App) error {
		a.httpServer = http.NewFiberServer(a.Config, a.container)
		a.Logger().Info("HTTP server initialized")
		return nil
	}
}

// Logger returns the logger instance associated with the App.
func (a *App) Logger() *zap.Logger {
	return a.container.Logger
}

// Shutdown gracefully shuts down the app's services
func (a *App) Shutdown(ctx context.Context) error {
	if db, err := a.DB(); err == nil {
		a.DB.Close()
	}

	if server, err := a.HTTPServer(); err == nil {
		return server.ShutdownWithContext(ctx)
	}

	if nats, err := a.Nats(); err == nil {
		db.Close()
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

// GetDB retrieves the database pool, logging an error if it is not configured.
func (a *App) DB() (*pgxpool.Pool, error) {
	if a.db == nil {
		a.Logger().Error("Database not configured")
		return nil, fmt.Errorf("database not configured")
	}
	return a.db, nil
}

// GetHTTPServer retrieves the HTTP server instance, logging an error if it is not configured.
func (a *App) HTTPServer() (HTTPServer, error) {
	if a.httpServer == nil {
		a.Logger().Error("HTTP server not configured")
		return nil, fmt.Errorf("HTTP server not configured")
	}
	a.Logger().Info("HTTP server retrieved successfully")
	return a.httpServer, nil
}

// GetPolicyManager retrieves the policy manager, logging an error if it is not configured.
func (a *App) PolicyManager() (*policy.PolicyManager, error) {
	if a.policyManager == nil {
		a.Logger().Error("Policy manager not configured")
		return nil, fmt.Errorf("policy manager not configured")
	}
	a.Logger().Info("Policy manager retrieved successfully")
	return a.policyManager, nil
}

// GetPublisher retrieves the messaging publisher, logging an error if it is not configured.
func (a *App) Publisher() (messaging.Messaging, error) {
	if a.messaging == nil {
		a.Logger().Error("Messaging client not configured")
		return nil, fmt.Errorf("messaging client not configured")
	}
	a.Logger().Info("Messaging publisher retrieved successfully")
	return a.messaging, nil
}

// GetSubscriber retrieves the messaging subscriber, logging an error if it is not configured.
func (a *App) Subscriber() (messaging.Messaging, error) {
	if a.messaging == nil {
		a.Logger().Error("Messaging client not configured")
		return nil, fmt.Errorf("messaging client not configured")
	}
	a.Logger().Info("Messaging subscriber retrieved successfully")
	return a.messaging, nil
}

// AddHTTPService registers HTTP service in container.
func (a *App) RegisterService(serviceName, serviceAddress string, service interface{}) {
	if a.container.Services == nil {
		a.container.Services = &registry.ServiceRegistry{}
	}

	if _, ok := a.container.Services.Get(serviceName); ok {
		a.Logger().Debug("Service already registered", zap.String("Name", serviceName))
	}

	a.container.Services.Register(serviceName, service)
}

// GET adds a Handler for HTTP GET method for a route pattern.
func (a *App) GET(pattern string, handler HandlerFunc) {
	a.httpServer.Router().GET(pattern, a.wrapHandler(handler))
}

// PUT adds a Handler for HTTP PUT method for a route pattern.
func (a *App) PUT(pattern string, handler HandlerFunc) {
	a.httpServer.Router().PUT(pattern, a.wrapHandler(handler))
}

// POST adds a Handler for HTTP POST method for a route pattern.
func (a *App) POST(pattern string, handler HandlerFunc) {
	a.httpServer.Router().POST(pattern, a.wrapHandler(handler))
}

// DELETE adds a Handler for HTTP DELETE method for a route pattern.
func (a *App) DELETE(pattern string, handler HandlerFunc) {
	a.httpServer.Router().DELETE(pattern, a.wrapHandler(handler))
}

// PATCH adds a Handler for HTTP PATCH method for a route pattern.
func (a *App) PATCH(pattern string, handler HandlerFunc) {
	a.httpServer.Router().PATCH(pattern, a.wrapHandler(handler))
}

// wrapHandler wraps a `HandlerFunc` to create a `fiber.Handler`.
func (a *App) wrapHandler(handler HandlerFunc) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := http.NewFiberRequestAdapter(c)
		res := http.NewFiberResponderAdapter(c)
		ctx := newContext(res, req, a.container)

		// Call the provided handler with the new context
		result, err := handler(ctx)

		// Send the response using the context's responder
		ctx.responder.Respond(result, err)
		return nil
	}
}

// Metrics returns the metrics manager associated with the App.
/* func (a *App) Metrics() metrics.Manager {
	return a.container.Metrics()
} */
