// neodata-go/http/server.go
package http

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/neodata/container"
)

// FiberServer implements the HTTPServer interface using Fiber.
type FiberServer struct {
	app       *fiber.App
	router    *Router
	container *container.Container
}

// SetupHTTPServer initializes a new HTTP server with the provided configuration and middleware.
func NewFiberServer(cfg config.ConfigProvider, c *container.Container) *FiberServer {
	conf := cfg.GetAppConfig()

	// Create Fiber app with custom configuration
	app := fiber.New(fiber.Config{
		ReadTimeout:  conf.App.ReadTimeout * time.Second,
		WriteTimeout: conf.App.WriteTimeout * time.Second,
		AppName:      conf.App.Name,
	})

	// Middleware setup
	app.Use(CorrelationIDMiddleware()) // CorrelationIDMiddleware for all requests
	app.Use(ZapLoggerMiddleware(c.Logger))

	// Enable CORS if environment is development
	if conf.App.Env == "dev" {
		// Allow all methods and headers from localhost for development purposes
		app.Use(cors.New())
	}
	// app.Use(RateLimiterMiddleware(100, time.Minute)) // Rate limiting for protected endpoints

	// Initialize Router with Fiber app and shared container
	r := NewRouter(app, c)

	// Return the configured FiberServer
	return &FiberServer{
		app:       app,
		router:    r,
		container: c,
	}
}

// Start launches the Fiber server on the specified port.
func (s *FiberServer) Start(port int) error {
	if err := s.app.Listen(fmt.Sprintf(":%d", port), fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server
func (s *FiberServer) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// Router provides access to route registration functions.
func (s *FiberServer) Router() *Router {
	return s.router
}
