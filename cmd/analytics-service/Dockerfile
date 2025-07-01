# cmd/analytics-service/Dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the analytics service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o analytics-service ./cmd/analytics-service

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/analytics-service .

# Expose port (not needed for this service but good practice)
EXPOSE 8082

# Run the analytics service
CMD ["./analytics-service"]