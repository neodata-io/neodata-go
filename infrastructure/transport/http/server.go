// neodata-go/http/server.go
package http

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/neodata-io/neodata-go/config"
)

// SetupHTTPServer initializes a new HTTP server with the provided configuration and middleware.
func SetupHTTPServer(cfg *config.AppConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.App.ReadTimeout * time.Second,
		WriteTimeout: cfg.App.WriteTimeout * time.Second,
		AppName:      cfg.App.Name,
	})

	// Middleware setup (commented out; add as needed)
	// app.Use(LoggerMiddleware())                    // Log all incoming requests
	// app.Use(RateLimiterMiddleware(100, time.Minute)) // Rate limiting for protected endpoints

	return app
}

// StartServer starts the HTTP server on the specified port.
func StartServer(app *fiber.App, port string) {
	log.Printf("Starting server on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
