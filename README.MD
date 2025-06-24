# A URL Shortener

## Quick Summary 

Takes a URL and transforms it to a shortcode that can be used to retrieve the original URL effectively shortening the original URL.

Currently being refactored into a distributed architecture consisting of these services:

- URL-service completed

- Gateway TODO

- Analytics TODO

## Prerequisites

If you would like to build, run, and test this project with docker compose you will need to have the following items installed

1. [Docker](https://docs.docker.com/get-docker/) 

2. [Docker Compose](https://docs.docker.com/compose/install/) (may be included with docker)

## Endpoints:

The fields needed for the GRPC endpoints for url-service can be seen in proto/url_service.proto

    rpc CreateShortURL(CreateURLRequest) returns (CreateURLResponse);
    
    rpc GetOriginalURL(GetURLRequest) returns (GetURLResponse);
    
    rpc HealthCheck(HealthRequest) returns (HealthResponse);


## Docker Instructions

To run the url-service you can use the docker compose command to spin up the url service as well as the postgres database.

`
docker compose up -d --build
`



## Testing

Current unit tests exists for url-service under internal/service/url-service_test.go.

End to end testing exists under cmd/test-client/main.go

To test first ensure url-service and postgres are running. 
`
docker ps
`
If they are not running then run the command to build the docker images and start them

`
docker compose up -d --build
`

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

For full validation server needs to be hosted and endpoints tested manually end to end.

### Example Usage (Minimal)

Currently only the internal URL service exists. You can create and retrieve URLs by accessing the service through gRPC.

1. Build images and spin up docker containers  (See Docker Instructions)

2. Health Check

`
grpcurl -plaintext -d '{}' localhost:50051 urlservice.URLService.HealthCheck
`

3. Create a short URL (you need to pass a URL as well as a user id)


`
grpcurl -plaintext -d '{"original_url":"https://example.com","user_id":"user123"}' localhost:50051 urlservice.URLService.CreateShortURL
`

Response:

`
{
  "shortCode": "UBs3cx",
  "shortUrl": "http://localhost:8080/UBs3cx",
  "success": true
}
`

4. Retrieve the original URL from the short code generated 

`
grpcurl -plaintext -d '{"short_code":"UBs3cx"}' localhost:50051 urlservice.URLService.G
etOriginalURL
`

Response: 

`
{
  "originalUrl": "https://example.com",
  "found": true
}
`

## Goals and Roadmap 

### Original Goals Met

Sprint 1: Basic URL shortening service with core API endpoints and unit tests. (Complete)

Sprint 2: Add persistent storage (e.g., database) instead of in-memory maps. (Complete)

Sprint 3: Add analytics (click tracking, usage stats). (Complete)

Sprint 4: Implement user auth and management. (Pushed back)

Further goals include deployment to a cloud environment with CI/CD pipeline

### New Goals (Distributed Architecture)
Current goals are to build out this project with distributed architecture.
There will be 3 services

1. url-service 
    - internal service handling URL creation and deletion
    - database interactions and URL logic
2. gateway-service
    - external facing REST API
    - translates REST requests to gRPC calls to url-service
    - send events to analytics service on URL hits
3. analytics-service
    - internal service that will track user clicks and usage
    - aggregate data for reporting and insights.


### Sprint Plan
Sprint 1: Microservice architecture split from monolith (Complete)

Sprint 2: gRPC server setup for internal communication (Complete)

Sprint 3: url-service implementation with full DB integration (Complete)

Sprint 4: gateway-service implementation with REST to gRPC translation (In progress)

Sprint 5: analytics-service implementation and event pipeline setup

Sprint 6: User authentication and authorization service

Sprint 7: Cloud deployment with CI/CD pipelines and monitoring

Sprint 8: Performance optimization and horizontal scaling of services