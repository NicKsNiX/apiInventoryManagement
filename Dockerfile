# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git and build-base for any dependencies
RUN apk add --no-cache git build-base

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o apiInventoryManagement .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for HTTPS and other runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/apiInventoryManagement .

# Copy uploads directory
RUN mkdir -p /root/uploads
COPY --from=builder /app/uploads /root/uploads 2>/dev/null || true

# Copy environment template if exists
COPY .env* /root/ 2>/dev/null || true

EXPOSE 4002

CMD ["./apiInventoryManagement"]
