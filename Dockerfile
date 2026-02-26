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

# Copy uploads directory if it exists
RUN mkdir -p /root/uploads

EXPOSE 4002

CMD ["./apiInventoryManagement"]
