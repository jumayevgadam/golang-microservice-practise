package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	ObserveLatency(path string, duration float64)
	IncError(path string)
}

var _ Metrics = &AppMetrics{}

type AppMetrics struct {
	ResponseLatency *prometheus.HistogramVec
	ErrorsTotal     *prometheus.CounterVec
}

func RegisterMetrics() *AppMetrics {
	responseLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_duration_seconds",
			Help: "Histogram of response latencies",
		},
		[]string{"path"},
	)

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failed_requests_total",
			Help: "Increments on every failed HTTP request",
		},
		[]string{"path"},
	)

	prometheus.MustRegister(responseLatency, errorCounter)

	return &AppMetrics{
		ResponseLatency: responseLatency,
		ErrorsTotal:     errorCounter,
	}
}

func (a *AppMetrics) ObserveLatency(path string, duration float64) {
	a.ResponseLatency.With(prometheus.Labels{"path": path}).Observe(duration)
}

func (a *AppMetrics) IncError(path string) {
	a.ErrorsTotal.With(prometheus.Labels{"path": path}).Inc()
}
