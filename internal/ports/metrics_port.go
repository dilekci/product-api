package ports

type MetricsPort interface {
	IncHTTPRequests(method, path, status string)
	ObserveHTTPLatency(method, path string, seconds float64)
}
