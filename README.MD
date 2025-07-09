# URL Shortener - Distributed Microservices Platform

A Go-based microservices URL shortener.

## Overview

A production-simulated microservices platform written in Go, containerized with Docker, deployed to Kubernetes, and tested via CI/CD. Includes REST/gRPC interfaces, asynchronous event logging, and observability hooks.

## Project Structure
```
url-shortener/
├── cmd/                         # Entry points for services and clients
│   ├── gateway-service/         # REST API → gRPC gateway
│   ├── url-service/             # Core URL logic and persistence
│   ├── analytics-service/       # Asynchronous event consumer
│   └── test-client/             # Manual and integration test clients
├── proto/                       # gRPC protobufs
├── internal/                    # Core application code
│   ├── analytics/               # Analytics logic and tests
│   ├── database/                # Database connection & handling
│   ├── events/                  # Event definitions and handling
│   ├── gateway/                 # Gateway handlers and tests
│   ├── models/                  # Data models
│   ├── queue/                   # Queue interface and Redis streams
│   ├── repository/              # Data repositories (including Postgres)
│   └── service/                 # URL Service logic and tests
├── deployments/                 # Kubernetes manifests
├── .github/workflows/           # GitHub Actions CI config
├── docker-compose.yml           # Local container orchestration
└── README.md                    # Project overview
```

## Tech Stack:
- Language: Go
- Protocol: REST, gRPC
- Data & Messaging: PostgreSQL, Redis Streams
- Infra & Platform: Docker, Kubernetes, GitHub Actions
- Testing: Unit tests, integration test clients, CI coverage

## Architecture

### Component Breakdown:
- **gateway-service**: Accepts REST requests and forwards them via gRPC to internal services.
- **url-service**: Core logic service that generates shortcodes and persists mappings to PostgreSQL.
- **analytics-service**: Subscribes to Redis Streams to consume and log access events asynchronously.
- **Redis**: Acts as the message broker using Redis Streams for event delivery.
- **PostgreSQL**: Stores URL mappings and access statistics.


### Data Flow:
- HTTP REST request -> Gateway -> gRPC -> url-service
- gateway-service publishes access events to Redis Stream → analytics-service consumes and logs them

## Intent 
- Build realistic backend infra with decoupled components
- Apply containerization, orchestration, and CI/CD principles
- Practice patterns like event-driven design and infra-as-code
- Simulate scalable deployment with K8s and Prometheus-ready endpoints

## Running the Project

### Local: Docker Compose

#### Prerequisites

1. [Docker](https://docs.docker.com/get-docker/) 

2. [Docker Compose](https://docs.docker.com/compose/install/) (may be included with docker)

#### Run Services

To run all of the services you can use the following docker compose command.

```
docker compose up -d --build
```

### Local: Kubernetes: Minikube

#### Prerequisites

[Minikube](https://minikube.sigs.k8s.io)

[kubectl](https://kubernetes.io/docs/tasks/tools/)

#### Run Services

The following commands will run all of the services and open the port to your local machine for testing.
You will need to port forward url-service to directly test gRPC endpoints. (Done in test clients)
```
minikube start
kubectl apply -f deployments/
kubectl port-forward svc/gateway-service 8080:8080
kubectl port-forward svc/url-service 50051:50051 
```

## Testing

### Integration Test Clients

Run after services are up. Confirm with `docker ps`, or `kubectl get pods`

To run all integration tests.
```
./testingscripts/run-test-clients.sh
```
To run the test client run this command

`
go run cmd/test-client/main.go
`

To run unit tests with coverage 

`
go test ./... -coverprofile=coverage.out
`

To see an HTML report

`
go tool cover -html=coverage.out -o coverage.html
`



## Endpoints:

Defined in proto/url_service.proto

    rpc CreateShortURL(CreateURLRequest) returns (CreateURLResponse);
    
    rpc GetOriginalURL(GetURLRequest) returns (GetURLResponse);
    
    rpc HealthCheck(HealthRequest) returns (HealthResponse);

Example usage with grpcurl:
```
grpcurl -plaintext -d '{"original_url":"https://example.com","user_id":"user123"}' localhost:50051 urlservice.URLService.CreateShortURL
```


## Roadmap 

### Completed 
- Core shortening logic (in-memory → Postgres + Redis Cache)

- Microservice decomposition

- gRPC + REST translation via gateway

- Redis Streams → analytics event processing

- CI/CD pipeline + integration testing

- Kubernetes deployment manifests

### In Progress
- Prometheus metrics exposure

- Service resource tuning for K8s

- Observability dashboards

## Appendix
| Method | Path           | Description              |
| ------ | -------------- | ------------------------ |
| POST   | `/create`      | Create shortened URL     |
| GET    | `/{shortcode}` | Redirect to original URL |
| GET    | `/healthz`     | Service health check     |

Example usage:

```
curl -X POST -H "Content-Type: application/json" -d '{"url": "https://example.com"}' http://localhost:8080/create
```