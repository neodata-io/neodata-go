// http/client.go
package http

import (
	"net/http"
	"time"
)

// NewHTTPClient initializes a shared HTTP client with custom settings.
// timeout: timeout duration for requests, e.g., 10 * time.Second
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}
}
