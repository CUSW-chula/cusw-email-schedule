# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Set timezone to Asia/Bangkok (adjust as needed)
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy .env file if it exists (optional)
COPY .env* ./

# Make the binary executable
RUN chmod +x ./main

# Expose port if needed (adjust as required)
# EXPOSE 8080

# Run the binary
CMD ["./main"]
