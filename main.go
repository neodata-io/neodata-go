package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/yourcompany/neodata-go/config"
	"github.com/yourcompany/neodata-go/database"
	"github.com/yourcompany/neodata-go/http"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database connection
	dbPool, err := database.NewPostgresPool(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Setup Fiber app with neodata-go middleware
	app := fiber.New()
	app.Use(http.LoggerMiddleware())
	app.Use(http.RateLimiterMiddleware(100, time.Minute))

	// Define routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is running")
	})

	// Start server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
