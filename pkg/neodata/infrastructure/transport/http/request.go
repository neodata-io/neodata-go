package http

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

// Request is an abstraction over the underlying http.Request. This abstraction is useful because it allows us
// to create applications without being aware of the transport. cmd.Request is another such abstraction.
type FiberRequestAdapter struct {
	req fiber.Ctx // HTTP request and response context
}

func NewFiberRequestAdapter(r fiber.Ctx) *FiberRequestAdapter {
	return &FiberRequestAdapter{
		req: r,
	}
}

// to access the context associated with the incoming request
func (r *FiberRequestAdapter) Context() context.Context {
	return r.req.Context()
}

// to access the query parameters present in the request, it returns the value of the key provided
func (r *FiberRequestAdapter) Param(key string) string {
	return r.req.Query(key)
}

// to retrieve the path parameters
func (r *FiberRequestAdapter) PathParam(key string) string {
	return r.req.Params(key)
}

// to access a decoded format of the request body, the body is mapped to the interface provided
func (r *FiberRequestAdapter) Bind(v interface{}) error {
	return r.req.Bind().Body(v)
}

// to access the host name for the incoming request
func (r *FiberRequestAdapter) HostName() string {
	return r.req.Hostname()
}

// to access all query parameters for a given key returning slice of string
func (r *FiberRequestAdapter) Params() map[string]string {
	return r.req.Queries()
}
