package metrics

type NoopMetrics struct{}

func (m *NoopMetrics) IncHTTPRequest(service, method, endpoint string) {}

func (m *NoopMetrics) ObserveHTTPRequestDuration(service, method, endpoint string, seconds float64) {}

func (m *NoopMetrics) IncHTTPError(service, method, endpoint string, statusCode int) {}

func (m *NoopMetrics) IncGRPCCall(service, method string) {}

func (m *NoopMetrics) IncGRPCError(service, method string) {}

func (m *NoopMetrics) ObserveGRPCLatency(service, method string, seconds float64) {}

func (m *NoopMetrics) IncPublishEvent(service, eventName string) {}

func (m *NoopMetrics) IncPublishEventError(service, eventName string) {}

func (m *NoopMetrics) ObservePublishEventLatency(service, eventName string, seconds float64) {}

func (m *NoopMetrics) IncConsumeEvent(service, eventName string) {}

func (m *NoopMetrics) IncConsumeEventError(service, eventName string) {}

func (m *NoopMetrics) ObserveConsumeEventLatency(service, eventName string, seconds float64) {}

func (m *NoopMetrics) IncCacheHit(service, cacheName string) {}

func (m *NoopMetrics) IncCacheMiss(service, cacheName string) {}

func (m *NoopMetrics) IncCacheError(service, cacheName, operation string) {}

func (m *NoopMetrics) IncDBOperation(service, operation string) {}

func (m *NoopMetrics) IncDBError(service, operation string) {}

func (m *NoopMetrics) ObserveDBOperationDuration(service, operation string, seconds float64) {}
