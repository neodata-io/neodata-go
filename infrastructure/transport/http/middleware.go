package http

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/neodata-io/neodata-go/domain/entities"
	"go.uber.org/zap"
)

type ValidationResponse struct {
	Valid bool   `json:"valid"`
	Sub   string `json:"sub"`
	Role  string `json:"role"`
}

// ZapLoggerMiddleware logs request details using Zap logger
func ZapLoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now() // Capture start time

		err := c.Next() // Process request

		// Calculate latency
		latency := time.Since(start)

		// Log details of the request
		logger.Info("Request",
			zap.String("method", c.Method()),                                  // HTTP method (GET, POST, etc.)
			zap.String("path", c.Path()),                                      // Request path
			zap.Int("status", c.Response().StatusCode()),                      // Status code (200, 404, etc.)
			zap.Duration("latency", latency),                                  // Time taken to process the request
			zap.String("correlation_id", c.Locals("correlation_id").(string)), // Example: custom request ID
		)

		return err // Return any errors from the next handler
	}
}

// LoggerMiddleware provides structured and stylized request logging.
func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency} ${locals:requestid} ${body} \n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout, // Ensures output to standard log console.
	})
}

// RateLimiterMiddleware provides rate limiting based on request count per time unit.
func RateLimiterMiddleware(maxRequests int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: duration,
	})
}

func CorrelationIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Check if Correlation ID exists in the incoming request header
		correlationID := c.Get("X-Correlation-ID")
		if correlationID == "" {
			// Generate a new Correlation ID if not provided
			correlationID = uuid.New().String()
			c.Set("X-Correlation-ID", correlationID)
		}

		// Attach the Correlation ID to the request context for use in other functions
		c.Locals("correlation_id", correlationID)

		return c.Next()
	}
}

// AuthMiddleware validates tokens by calling the auth service.
func AuthMiddleware(secretKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if len(authHeader) <= len("Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header missing or malformed"})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token.
		token, err := jwt.ParseWithClaims(tokenString, &entities.Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what we expect (HS256).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key for token validation.
			return []byte(secretKey), nil
		})

		// Handle token parsing or validation errors.
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// Extract claims and store user data in context
		if claims, ok := token.Claims.(*entities.Claims); ok && token.Valid {
			c.Locals("userID", claims.UserID)
			c.Locals("abilities", claims.Abilities)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// Continue to the next middleware or handler.
		return c.Next()
	}
}
