package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rubenclaes/neodata-go/config"
	"github.com/rubenclaes/neodata-go/db"
	"github.com/rubenclaes/neodata-go/http"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database connection
	dbPool, err := db.NewPostgresPool(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Setup Fiber app with neodata-go middleware
	app := fiber.New()
	app.Use(http.LoggerMiddleware())
	app.Use(http.RateLimiterMiddleware(100, time.Minute))

	// Define routes
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("Service is running")
	})

	// Start server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
