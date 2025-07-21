.PHONY: ci build test staticAnalysis docker-ci tidy docker-build-url docker-build-gateway docker-build-analytics docker-build-all

ci: tidy build test staticAnalysis

docker-ci: test staticAnalysis docker-build-all docker-push-all

tidy:
	go mod tidy

build: 
	@echo "Verifying Compilation"
	go build ./cmd/url-service
	go build ./cmd/gateway-service
	go build ./cmd/analytics-service

test:
	@echo "Running unit tests"
	go test ./...

staticAnalysis:
	@echo "Running static analysis..."
	go vet ./...

docker-build-url:
	@echo "Building url-service"
	docker build -t ghcr.io/sammyqtran/url-service:latest -f ./cmd/url-service/Dockerfile .

docker-build-gateway:
	@echo "Building gateway-service"
	docker build -t ghcr.io/sammyqtran/gateway-service:latest -f ./cmd/gateway-service/Dockerfile .

docker-build-analytics:
	@echo "Building analytics-service"
	docker build -t ghcr.io/sammyqtran/analytics-service:latest -f ./cmd/analytics-service/Dockerfile .


docker-build-all: docker-build-url docker-build-gateway docker-build-analytics

docker-push-url:
	docker push ghcr.io/sammyqtran/url-service:latest

docker-push-gateway:
	docker push ghcr.io/sammyqtran/gateway-service:latest

docker-push-analytics:
	docker push ghcr.io/sammyqtran/analytics-service:latest

docker-push-all: docker-push-url docker-push-gateway docker-push-analytics


