# Use the official Golang image to create a build artifact.
FROM golang:1.20

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container.
COPY . .
# Set the environment variable for the base path
ENV BASE_PATH=/app/internal/server


# Build the Go app
RUN go build -o server ./internal/server

# Expose port 1323 to the outside world.
EXPOSE 1323

# Command to run the executable.
CMD ["./server"]