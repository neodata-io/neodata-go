package neodata

import (
	"context"

	"github.com/neodata-io/neodata-go/neodata/container"
	"github.com/neodata-io/neodata-go/neodata/interfaces"
)

type Context struct {
	context.Context // Embeds Go's context for cancellation and timeouts

	// Request needs to be public because handlers need to access request details. Else, we need to provide all
	// functionalities of the Request as a method on context. This is not needed because Request here is an interface
	// So, internals are not exposed anyway.
	interfaces.Request

	// Same logic as above.
	*container.Container // Provides access to shared dependencies

	// responder is private as Handlers do not need to worry about how to respond. But it is still an abstraction over
	// normal response writer as we want to keep the context independent of http. Will help us in writing CMD application
	// or gRPC servers etc using the same handler signature.
	responder interfaces.Responder
}

/* func (c *Context) Trace(name string) trace.Span {
	tr := otel.GetTracerProvider().Tracer("gofr-context")
	ctx, span := tr.Start(c.Context, name)
	// TODO: If we don't close the span using `defer` and run the http-server example by hitting `/trace` endpoint, we are
	// getting incomplete redis spans when viewing the trace using correlationID. If we remove assigning the ctx to GoFr
	// context then spans are coming correct but then parent-child span relationship is being hindered.

	c.Context = ctx

	return span
}

func (c *Context) Bind(i interface{}) error {
	return c.Request.Bind(i)
} */

func newContext(w interfaces.Responder, r interfaces.Request, c *container.Container) *Context {
	return &Context{
		Context:   r.Context(),
		Request:   r,
		responder: w,
		Container: c,
	}
}
