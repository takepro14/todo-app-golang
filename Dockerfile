# Dockerfile for Go TODO App

# Use the official Go image as a base
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application's default port
EXPOSE 8080

# Start the application
CMD ["./main"]
