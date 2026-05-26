# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod .

# Copy source code
COPY . .

# Build the application
RUN go build -o app .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/app .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"]
