# Use the official Golang image with Alpine as a base image for building the app
FROM golang:1.24.1-alpine AS build

# Set the working directory inside the container to /app
WORKDIR /app

# Copy the go.mod and go.sum files into the container (this is for Go dependency management)
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Build the Go application from the cmd/api/main.go file
# The binary will be named "main" and located at /app/main
RUN go build -o main cmd/api/main.go

# Create a new stage using the Alpine image for the production environment
FROM alpine:3.20.1 AS prod

# Set the working directory in the production container to /app
WORKDIR /app

# Copy the built binary "main" from the build stage to the production stage
COPY --from=build /app/main /app/main

# Expose the port that the application will run on
# The PORT variable should be set in the environment or docker-compose
EXPOSE ${PORT}

# Define the command to run when the container starts (running the main binary)
CMD ["./main"]
