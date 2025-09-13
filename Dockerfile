# Stage 1: The build environment
FROM golang:1.24.6-alpine AS builder

# Install the C compiler (gcc) needed to build CGO packages
RUN apk add --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application, allowing CGO this time
RUN go build -o /server


# Stage 2: The final, production-ready container
# Use Alpine as the base image because it contains the C libraries our app needs
FROM alpine:latest

# Copy the compiled server binary from the 'builder' stage
COPY --from=builder /server /server

# Copy the database file into the container
# NOTE: In a real-world app, the database would live outside the container.
# This is fine for our learning project.
COPY todos.db /todos.db

# Tell Docker which port the container will listen on
EXPOSE 8888

# The command to run when the container starts
CMD ["/server"]