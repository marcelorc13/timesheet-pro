# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and go.sum files for dependency caching
COPY go.mod go.sum ./

# Download dependencies (will be cached if go.mod/go.sum don't change)
RUN go mod download

# Copy the entire project
COPY . .

# Generate templ templates using go tool (matching Makefile)
RUN go tool templ generate

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o main ./cmd/api/main.go

# Final stage - minimal production image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port (Render uses PORT env variable)
EXPOSE 8080

# Run the application
CMD ["./main"]
