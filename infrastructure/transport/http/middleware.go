package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neodata-io/neodata-go/domain/entities"
)

type ValidationResponse struct {
	Valid bool   `json:"valid"`
	Sub   string `json:"sub"`
	Role  string `json:"role"`
}

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

		// Extract claims and store user information in context for use in handlers.
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
