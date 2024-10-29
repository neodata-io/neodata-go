package neodata

import (
	"github.com/gofiber/fiber/v3"

	"go.uber.org/zap"
)

// RegisterRoute allows microservices to add custom routes with a simple syntax.
func RegisterRoute(ctx *NeoCtx, method, path string, handler func(*NeoCtx) (interface{}, error)) {
	httpServer, err := ctx.GetHTTPServer()
	if err != nil {
		ctx.Logger.Fatal("Err", zap.Error(err))
	}

	// Wrap handler function to match Fiber's handler signature
	fiberHandler := func(c fiber.Ctx) error {
		result, err := handler(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	}

	// Register the route based on the method
	switch method {
	case "GET":
		httpServer.Get(path, fiberHandler)
	case "POST":
		httpServer.Post(path, fiberHandler)
	case "PUT":
		httpServer.Put(path, fiberHandler)
	case "DELETE":
		httpServer.Delete(path, fiberHandler)
	default:
		ctx.Logger.Warn("Unsupported method", zap.String("method", method))
	}
}
