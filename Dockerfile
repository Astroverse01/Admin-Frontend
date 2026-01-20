# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o admin-be ./main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Set timezone (optional, adjust as needed)
ENV TZ=UTC

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/admin-be .

# Expose port
EXPOSE 8080

# Set environment variables (can be overridden)
ENV PORT=8080
ENV GIN_MODE=release

# Run the application
CMD ["./admin-be"]

