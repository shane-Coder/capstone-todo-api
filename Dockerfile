# Dockerfile

# Stage 1: The build environment
FROM golang:1.24.6-alpine AS builder

# Install Git, which is needed to fetch Go modules
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static, CGO-disabled binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

# Stage 2: The final, production-ready container
# Use a minimal 'scratch' image for a tiny and secure final image
FROM scratch

# Copy the compiled server binary from the 'builder' stage
COPY --from=builder /server /server

# Tell Docker which port the container will listen on
EXPOSE 8888

# The command to run when the container starts
CMD ["/server"]