# Use the official Golang image as the base image
FROM golang:1.21rc2-alpine3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose a port for the application to listen on
EXPOSE 8080

# Set the command to run the executable
CMD ["./main"]
