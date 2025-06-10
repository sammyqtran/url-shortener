# Start from a minimal Go image
FROM golang:1.24.3-alpine

# Install git (for go get), and set up working directory
RUN apk add --no-cache git
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o app .

# Set entrypoint
CMD ["./app"]
