// http/middleware.go
package http

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// LoggerMiddleware provides standardized request logging.
func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	})
}

// RateLimiterMiddleware provides rate limiting based on request count per time unit.
func RateLimiterMiddleware(maxRequests int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: duration,
	})
}
