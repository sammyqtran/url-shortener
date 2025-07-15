package metrics

type Metrics interface {
	// HTTP requests
	IncHTTPRequest(service, method, endpoint string)
	ObserveHTTPRequestDuration(service, method, endpoint string, seconds float64)
	IncHTTPError(service, method, endpoint string, statusCode int)

	// gRPC calls
	IncGRPCCall(service, method string)
	IncGRPCError(service, method string)
	ObserveGRPCLatency(service, method string, seconds float64)

	// Event publishing (gateway)
	IncPublishEvent(service, eventName string)
	IncPublishEventError(service, eventName string)
	ObservePublishEventLatency(service, eventName string, seconds float64)

	// Event consuming (analytics)
	IncConsumeEvent(service, eventName string)
	IncConsumeEventError(service, eventName string)
	ObserveConsumeEventLatency(service, eventName string, seconds float64)

	// Cache operations (URL service)
	IncCacheHit(service, cacheName string)
	IncCacheMiss(service, cacheName string)
	IncCacheError(service, cacheName string)

	// DB operations (URL service)
	IncDBOperation(service, operation string)
	IncDBError(service, operation string)
}
