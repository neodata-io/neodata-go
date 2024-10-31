package container

import (
	"fmt"

	"github.com/neodata-io/neodata-go/config"
	"github.com/neodata-io/neodata-go/logger"
	"github.com/neodata-io/neodata-go/neodata/registry"
	"go.uber.org/zap"
)

// Container holds shared concerns across the app, such as logging and config and dynamic service registry.
type Container struct {
	Logger   *zap.Logger               // Injected from the main application to enable structured logging
	Services *registry.ServiceRegistry // Add a dynamic service registry
}

// NewContainer initializes a container with common dependencies.
func NewContainer(cfg config.ConfigProvider) (*Container, error) {
	log, err := logger.InitServiceLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize default logger: %w", err)
	}

	return &Container{
		Services: &registry.ServiceRegistry{},
		Logger:   log,
	}, nil
}

func (c *Container) GetService(serviceName string) (interface{}, bool) {
	return c.Services.Get(serviceName)
}

func (c *Container) Close() {
	if c.Logger != nil {
		c.Logger.Sync()
	}
}
