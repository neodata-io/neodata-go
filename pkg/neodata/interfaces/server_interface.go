package interfaces

import (
	"context"

	"github.com/neodata-io/neodata-go/infrastructure/transport/http"
)

// HTTPServer defines the methods for an HTTP server interface.
type HTTPServer interface {
	Start(port int) error
	Shutdown(ctx context.Context) error
	Router() *http.Router
}
