package neodata

import (
	"time"

	"github.com/neodata-io/neodata-go/neodata/container"
)

type HandlerFunc func(c *Context) (interface{}, error)

type handler struct {
	function       HandlerFunc
	container      *container.Container
	requestTimeout time.Duration
}

func healthHandler(c *Context) (interface{}, error) {
	return c.Health(c), nil
}

func liveHandler(*Context) (interface{}, error) {
	return struct {
		Status string `json:"status"`
	}{Status: "UP"}, nil
}
