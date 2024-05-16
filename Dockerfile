# Use the official Golang image to create a build stage
FROM golang:1.17 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the main.go source file to the working directory
COPY main.go ./

# Build the Go app
RUN go build -o honeypot

# Start a new stage from scratch
FROM golang:1.17

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/honeypot .

# Expose port 2222 to the outside world
EXPOSE 2222

# Command to run the executable
CMD ["./honeypot"]