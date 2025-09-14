# Build
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum for the main module
COPY go.mod go.sum ./

# Copy shared packages
COPY pkg/ ./pkg/

# Copy service-specific go.mod and go.sum
COPY services/worker/go.mod services/worker/go.sum ./services/worker/

# Set working directory to the service
WORKDIR /app/services/worker

# Download dependencies
RUN go mod download

# Copy service source code
COPY services/worker/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

# Final
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/services/worker/main /app/services/worker/main
COPY --from=builder /app/services/worker/.env.dist /app/services/worker/.env
RUN chown -R appuser:appgroup /app
USER appuser
CMD ["/app/services/worker/main"]
