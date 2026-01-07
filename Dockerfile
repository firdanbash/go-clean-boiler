# Build stage
FROM golang:1.22-alpine AS builder


WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy source code
COPY . .

# Download dependencies and build
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api


# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
