# Use the official Go image as a base
FROM golang:1.22.3 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o blogging-platform-api .

# Create a minimal image for running the application
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/blogging-platform-api .

# Expose the port the app runs on
EXPOSE 8081

# Command to run the executable
CMD ["./blogging-platform-api"]
