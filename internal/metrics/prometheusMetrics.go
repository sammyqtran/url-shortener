package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics implements Metrics interface using Prometheus.
type PrometheusMetrics struct {
	httpRequests        *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpErrors          *prometheus.CounterVec
	grpcCalls           *prometheus.CounterVec
	grpcErrors          *prometheus.CounterVec
	grpcLatency         *prometheus.HistogramVec
	publishEvents       *prometheus.CounterVec
	publishEventErrors  *prometheus.CounterVec
	publishEventLatency *prometheus.HistogramVec
	consumeEvents       *prometheus.CounterVec
	consumeEventErrors  *prometheus.CounterVec
	consumeEventLatency *prometheus.HistogramVec
	cacheHits           *prometheus.CounterVec
	cacheMisses         *prometheus.CounterVec
	cacheErrors         *prometheus.CounterVec
	dbOperations        *prometheus.CounterVec
	dbErrors            *prometheus.CounterVec
	dbOperationDuration *prometheus.HistogramVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		httpRequests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		}, []string{"service", "method", "endpoint"}),
		httpRequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "method", "endpoint"}),
		httpErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total HTTP errors",
		}, []string{"service", "method", "endpoint", "status_code"}),
		grpcCalls: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "grpc_calls_total",
			Help: "Total gRPC calls",
		}, []string{"service", "method"}),
		grpcErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "grpc_errors_total",
			Help: "Total gRPC errors",
		}, []string{"service", "method"}),
		grpcLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "grpc_call_duration_seconds",
			Help:    "gRPC call duration in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "method"}),
		publishEvents: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "publish_events_total",
			Help: "Total events published",
		}, []string{"service", "event"}),
		publishEventErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "publish_event_errors_total",
			Help: "Total errors publishing events",
		}, []string{"service", "event"}),
		publishEventLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "publish_event_duration_seconds",
			Help:    "Event publishing duration in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "event"}),
		consumeEvents: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "consume_events_total",
			Help: "Total events consumed",
		}, []string{"service", "event"}),
		consumeEventErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "consume_event_errors_total",
			Help: "Total errors consuming events",
		}, []string{"service", "event"}),
		consumeEventLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "consume_event_duration_seconds",
			Help:    "Event consuming duration in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "event"}),
		cacheHits: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		}, []string{"service", "cache"}),
		cacheMisses: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total cache misses",
		}, []string{"service", "cache"}),
		cacheErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "cache_errors_total",
			Help: "Total cache errors",
		}, []string{"service", "cache", "operation"}),
		dbOperations: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "db_operations_total",
			Help: "Total DB operations",
		}, []string{"service", "operation"}),
		dbErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total DB errors",
		}, []string{"service", "operation"}),
		dbOperationDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets, // or customize: prometheus.ExponentialBuckets(0.001, 2, 12)
		}, []string{"service", "operation"}),
	}
}

// Implement Metrics interface methods

func (m *PrometheusMetrics) IncHTTPRequest(service, method, endpoint string) {
	m.httpRequests.WithLabelValues(service, method, endpoint).Inc()
}

func (m *PrometheusMetrics) ObserveHTTPRequestDuration(service, method, endpoint string, seconds float64) {
	m.httpRequestDuration.WithLabelValues(service, method, endpoint).Observe(seconds)
}

func (m *PrometheusMetrics) IncHTTPError(service, method, endpoint string, statusCode int) {
	m.httpErrors.WithLabelValues(service, method, endpoint, fmt.Sprintf("%d", statusCode)).Inc()
}

func (m *PrometheusMetrics) IncGRPCCall(service, method string) {
	m.grpcCalls.WithLabelValues(service, method).Inc()
}

func (m *PrometheusMetrics) IncGRPCError(service, method string) {
	m.grpcErrors.WithLabelValues(service, method).Inc()
}

func (m *PrometheusMetrics) ObserveGRPCLatency(service, method string, seconds float64) {
	m.grpcLatency.WithLabelValues(service, method).Observe(seconds)
}

func (m *PrometheusMetrics) IncPublishEvent(service, eventName string) {
	m.publishEvents.WithLabelValues(service, eventName).Inc()
}

func (m *PrometheusMetrics) IncPublishEventError(service, eventName string) {
	m.publishEventErrors.WithLabelValues(service, eventName).Inc()
}

func (m *PrometheusMetrics) ObservePublishEventLatency(service, eventName string, seconds float64) {
	m.publishEventLatency.WithLabelValues(service, eventName).Observe(seconds)
}

func (m *PrometheusMetrics) IncConsumeEvent(service, eventName string) {
	m.consumeEvents.WithLabelValues(service, eventName).Inc()
}

func (m *PrometheusMetrics) IncConsumeEventError(service, eventName string) {
	m.consumeEventErrors.WithLabelValues(service, eventName).Inc()
}

func (m *PrometheusMetrics) ObserveConsumeEventLatency(service, eventName string, seconds float64) {
	m.consumeEventLatency.WithLabelValues(service, eventName).Observe(seconds)
}

func (m *PrometheusMetrics) IncCacheHit(service, cacheName string) {
	m.cacheHits.WithLabelValues(service, cacheName).Inc()
}

func (m *PrometheusMetrics) IncCacheMiss(service, cacheName string) {
	m.cacheMisses.WithLabelValues(service, cacheName).Inc()
}

func (m *PrometheusMetrics) IncCacheError(service, cacheName, operation string) {
	m.cacheErrors.WithLabelValues(service, cacheName, operation).Inc()
}

func (m *PrometheusMetrics) IncDBOperation(service, operation string) {
	m.dbOperations.WithLabelValues(service, operation).Inc()
}

func (m *PrometheusMetrics) IncDBError(service, operation string) {
	m.dbErrors.WithLabelValues(service, operation).Inc()
}

func (m *PrometheusMetrics) ObserveDBOperationDuration(service, operation string, seconds float64) {
	m.dbOperationDuration.WithLabelValues(service, operation).Observe(seconds)
}
