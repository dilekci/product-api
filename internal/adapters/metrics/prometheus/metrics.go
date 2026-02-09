package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	httpRequests *prometheus.CounterVec
	httpLatency  *prometheus.HistogramVec
}

func New() *Metrics {
	m := &Metrics{
		httpRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		httpLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	prometheus.MustRegister(
		m.httpRequests,
		m.httpLatency,
	)

	return m
}

func (m *Metrics) IncHTTPRequests(method, path, status string) {
	m.httpRequests.WithLabelValues(method, path, status).Inc()
}

func (m *Metrics) ObserveHTTPLatency(method, path string, seconds float64) {
	m.httpLatency.WithLabelValues(method, path).Observe(seconds)
}
