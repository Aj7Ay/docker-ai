# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy all the source code
COPY . .

# Build the binary for a static, C-free build suitable for a minimal image
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker-ai ./cmd/docker-ai/main.go

# Stage 2: Create the minimal final image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/docker-ai .

# Set the entrypoint for the container
ENTRYPOINT ["/app/docker-ai"] 