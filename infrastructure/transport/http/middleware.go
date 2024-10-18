package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

const authServiceURL = "http://auth-service:3000/validate"

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
func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if len(authHeader) <= len("Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header missing or malformed"})
		}
		token := authHeader[len("Bearer "):]

		// Call auth microservice for validation
		validationResponse, err := validateTokenWithAuthService(token)
		if err != nil || !validationResponse.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		// Store user information in context for handler use
		c.Locals("userID", validationResponse.Sub)
		//ctx.Locals("role", validationResponse.Role)
		return c.Next()
	}
}

func validateTokenWithAuthService(token string) (*ValidationResponse, error) {
	client := NewHTTPClient(10 * time.Second)

	req, err := http.NewRequest("POST", authServiceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unauthorized")
	}

	var validationResponse ValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
		return nil, err
	}

	return &validationResponse, nil
}
