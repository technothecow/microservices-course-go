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

# Build Migrator
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./utils/migrator/cmd/app/main.go

# Final stage: use a minimal Docker image
FROM alpine:latest

# Set the working directory and copy the built binary from the builder stage.
WORKDIR /root/
COPY --from=builder /app .

# Command to run the executable.
CMD ["./migrator"]
