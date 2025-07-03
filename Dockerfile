# Start from the official Golang image for building
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o beerlund-api .

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/beerlund-api .

# Expose port (change if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./beerlund-api"]