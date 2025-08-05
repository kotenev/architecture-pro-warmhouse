# Build stage
FROM golang:1.24.4-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=auto go build -o device_simulator

# Runtime stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/device_simulator .

# Command to run the application
CMD ["./device_simulator"]