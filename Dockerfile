# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Debug: List what Docker sees
COPY . .
RUN ls -laR

# Restore dependencies
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Run stage
FROM alpine:latest

# Install timezone data and certificates
RUN apk add --no-cache tzdata ca-certificates

# Create a non-root user
RUN adduser -D frostbyte
USER frostbyte

WORKDIR /app

COPY --from=builder /app/main .

# NOTE: In production, mount .env, server.crt, and server.key using volumes.
# Do NOT copy them here to avoid baking secrets into the image.

EXPOSE 8080

CMD ["./main"]
