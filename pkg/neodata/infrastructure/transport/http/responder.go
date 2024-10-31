package http

import "github.com/gofiber/fiber/v3"

// FiberResponderAdapter wraps fiber.Ctx to implement the Responder interface
type FiberResponderAdapter struct {
	ctx fiber.Ctx
}

// NewFiberResponderAdapter initializes a new FiberResponderAdapter
func NewFiberResponderAdapter(ctx fiber.Ctx) *FiberResponderAdapter {
	return &FiberResponderAdapter{ctx: ctx}
}

// Respond sends a JSON response with the result or an error message
func (r *FiberResponderAdapter) Respond(data interface{}, err error) {
	if err != nil {
		r.ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	} else {
		r.ctx.JSON(data)
	}
}
