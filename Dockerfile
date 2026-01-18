# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Debug: List what Docker sees
COPY . .
RUN ls -laR

# Restore dependencies
RUN go mod download

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Build Seeder binary
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder ./cmd/seeder

# Run stage
FROM alpine:latest

# Install timezone data and certificates
RUN apk add --no-cache tzdata ca-certificates

# Ensure Go can find the timezone data
ENV ZONEINFO=/usr/share/zoneinfo

# Create a non-root user
RUN adduser -D frostbyte
USER frostbyte

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/seeder .

# NOTE: In production, mount .env, server.crt, and server.key using volumes.
# Do NOT copy them here to avoid baking secrets into the image.

EXPOSE 8080

CMD ["./main"]
