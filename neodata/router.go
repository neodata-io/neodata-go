package neodata

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// Router provides methods for registering HTTP routes
type Router struct {
	server *fiber.App
	ctx    *NeoCtx
}

// NewRouter initializes a new Router instance
func NewRouter(ctx *NeoCtx) *Router {
	httpServer, err := ctx.GetHTTPServer()
	if err != nil {
		ctx.Logger.Fatal("Err", zap.Error(err))
	}

	return &Router{
		server: httpServer,
		ctx:    ctx,
	}
}

// GET registers a GET route with a custom handler function
func (r *Router) GET(path string, handler func(*NeoCtx) (interface{}, error)) {
	r.server.Get(path, wrapHandler(r.ctx, handler))
}

// POST registers a POST route with a custom handler function
func (r *Router) POST(path string, handler func(*NeoCtx) (interface{}, error)) {
	r.server.Post(path, wrapHandler(r.ctx, handler))
}

// wrapHandler wraps a handler to match Fiber's handler signature
func wrapHandler(ctx *NeoCtx, handler func(*NeoCtx) (interface{}, error)) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	}
}
