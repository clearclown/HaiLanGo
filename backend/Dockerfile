# Build stage
FROM docker.io/library/golang:alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/server/main.go

# Runtime stage
FROM docker.io/library/alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/main .

# Copy mocks directory (for development/testing)
COPY --from=builder /build/mocks ./mocks

# Create necessary directories
RUN mkdir -p /app/uploads /app/logs

# Expose port
EXPOSE 8080

# Set timezone
ENV TZ=Asia/Tokyo

# Run the application
CMD ["./main"]
