# Build stage
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency caching.
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code.
COPY . .

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./posts/cmd/app/main.go

# Final stage: use a minimal Docker image
FROM alpine:3.21

# Set the working directory and copy the built binary from the builder stage.
WORKDIR /root/
COPY --from=builder /app/app .

# Expose port (adjust accordingly to your app)
EXPOSE 50003

# Command to run the executable.
CMD ["./app"]
