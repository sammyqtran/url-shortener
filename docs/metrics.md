# Metrics Overview

This document defines all Prometheus metrics collected across services.

---

## Gateway Service

### `http_requests_total`
- **Type**: Counter  
- **Description**: Total number of HTTP requests received  
- **Labels**: `service`, `method`, `endpoint`

### `http_request_duration_seconds`
- **Type**: Histogram  
- **Description**: Duration of HTTP requests  
- **Unit**: Seconds  
- **Labels**: `service`, `method`, `endpoint`

### `http_errors_total`
- **Type**: Counter  
- **Description**: Total number of failed HTTP requests  
- **Labels**: `service`, `method`, `endpoint`, `status_code`

---

### `grpc_calls_total`
- **Type**: Counter  
- **Description**: Total number of outbound gRPC calls  
- **Labels**: `service`, `method`

### `grpc_call_duration_seconds`
- **Type**: Histogram  
- **Description**: Duration of outbound gRPC calls  
- **Unit**: Seconds  
- **Labels**: `service`, `method`

### `grpc_errors_total`
- **Type**: Counter  
- **Description**: Total number of gRPC call failures  
- **Labels**: `service`, `method`

---

### `publish_events_total`
- **Type**: Counter  
- **Description**: Total number of events published  
- **Labels**: `service`, `event`

### `publish_event_errors_total`
- **Type**: Counter  
- **Description**: Total number of event publish failures  
- **Labels**: `service`, `event`

### `publish_event_duration_seconds`
- **Type**: Histogram  
- **Description**: Duration of event publishing  
- **Unit**: Seconds  
- **Labels**: `service`, `event`

---

## Analytics Service

### `consume_events_total`
- **Type**: Counter  
- **Description**: Total number of events consumed  
- **Labels**: `service`, `event`

### `consume_event_errors_total`
- **Type**: Counter  
- **Description**: Total number of errors while consuming events  
- **Labels**: `service`, `event`

### `consume_event_duration_seconds`
- **Type**: Histogram  
- **Description**: Duration of event processing  
- **Unit**: Seconds  
- **Labels**: `service`, `event`

---

## URL Service

### `cache_hits_total`
- **Type**: Counter  
- **Description**: Total number of cache hits  
- **Labels**: `service`, `cache`

### `cache_misses_total`
- **Type**: Counter  
- **Description**: Total number of cache misses  
- **Labels**: `service`, `cache`

### `cache_errors_total`
- **Type**: Counter  
- **Description**: Total number of cache operation errors  
- **Labels**: `service`, `cache`

---

### `db_operations_total`
- **Type**: Counter  
- **Description**: Total number of database operations  
- **Labels**: `service`, `operation`

### `db_errors_total`
- **Type**: Counter  
- **Description**: Total number of database errors  
- **Labels**: `service`, `operation`

### `db_operation_duration_seconds`
- **Type**: Histogram  
- **Description**: Duration of database operations  
- **Unit**: Seconds  
- **Labels**: `service`, `operation`

---
