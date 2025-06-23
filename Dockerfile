# Multi-stage build for Ospy
FROM golang:1.24-alpine AS builder

# Install git for version info
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o ospy ./cmd/ospy

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S ospy && \
    adduser -u 1001 -S ospy -G ospy

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/ospy .

# Copy configuration template
COPY --from=builder /app/configs/config.yaml ./configs/

# Create data and logs directories
RUN mkdir -p data logs && \
    chown -R ospy:ospy /app

# Switch to non-root user
USER ospy

# Expose web interface port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./ospy -version || exit 1

# Default command
CMD ["./ospy", "-config", "configs/config.yaml"]
