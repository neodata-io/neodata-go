package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/neodata-io/neodata-go/neodata/container"
)

// Router provides methods for registering HTTP routes
type Router struct {
	server    *fiber.App
	container *container.Container
}

// NewRouter initializes a new Router instance
func NewRouter(s *fiber.App, c *container.Container) *Router {
	return &Router{
		server:    s,
		container: c,
	}
}

// GET registers a GET route with a Fiber-compatible handler.
func (r *Router) GET(path string, handler fiber.Handler) {
	r.server.Get(path, handler)
}

// POST registers a POST route with a Fiber-compatible handler.
func (r *Router) POST(path string, handler fiber.Handler) {
	r.server.Post(path, handler)
}

// PUT registers a PUT route with a Fiber-compatible handler.
func (r *Router) PUT(path string, handler fiber.Handler) {
	r.server.Put(path, handler)
}

// DELETE registers a DELETE route with a Fiber-compatible handler.
func (r *Router) DELETE(path string, handler fiber.Handler) {
	r.server.Delete(path, handler)
}

// PATCH registers a PATCH route with a Fiber-compatible handler.
func (r *Router) PATCH(path string, handler fiber.Handler) {
	r.server.Patch(path, handler)
}
