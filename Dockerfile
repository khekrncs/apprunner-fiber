FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy and build the application
COPY . .
# Build a static binary for better compatibility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fiber-app

# Use a clean Alpine image for the final container
FROM alpine:latest
WORKDIR /app

# Add CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/fiber-app /app/
RUN chmod +x ./fiber-app

# Set port
EXPOSE 8080

# Define the command to run your service
CMD ["./fiber-app"]