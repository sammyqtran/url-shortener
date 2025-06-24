.PHONY: proto build run-server test clean

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/url_service.proto

# Build the URL service
build:
	go build -o bin/url-service cmd/url-service/main.go

# Run the gRPC server
run-server:
	go run cmd/url-service/main.go

# Test the service
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
