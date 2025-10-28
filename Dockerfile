# Start from the official Go image
FROM golang:1.25 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o shortener ./cmd/shortener

# Start a minimal image
FROM ubuntu:latest
WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/shortener ./shortener

# Expose the port (change if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./shortener"]
