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

# Install goose for database migrations
RUN go install -ldflags="-s -w" -tags="no_libsql no_mssql no_vertica no_clickhouse no_mysql no_sqlite3 no_ydb" github.com/pressly/goose/v3/cmd/goose@latest

# Final stage - minimal production image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy goose binary from builder
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration files
COPY --from=builder /app/internal/repository/migrations /app/migrations

# Copy entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Expose port (Render uses PORT env variable)
EXPOSE 8080

# Use entrypoint script to run migrations and start the app
ENTRYPOINT ["docker-entrypoint.sh"]
