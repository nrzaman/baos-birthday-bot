# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bot .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs tzdata

# Create non-root user
RUN addgroup -g 1000 botuser && \
    adduser -D -u 1000 -G botuser botuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bot .

# Create directory for database
RUN mkdir -p /app/data && \
    chown -R botuser:botuser /app

# Switch to non-root user
USER botuser

# Run the bot
CMD ["./bot"]
