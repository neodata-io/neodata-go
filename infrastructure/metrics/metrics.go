package neodata

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	RequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		}, []string{"path", "method", "status"}),
	}
}

func (n *NeoCtx) StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		n.Logger.Info("Starting metrics server on :9090/metrics")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			n.Logger.Error("Metrics server failed", zap.Error(err))
		}
	}()
}
