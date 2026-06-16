# Start with a Go image to build the application
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server/main.go

# Start a new stage from a minimal base image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/bin/server .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./server"]
