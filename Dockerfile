# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Run stage
FROM alpine:latest

# Create a non-root user
RUN adduser -D frostbyte
USER frostbyte

WORKDIR /app

COPY --from=builder /app/main .

# NOTE: In production, mount .env, server.crt, and server.key using volumes.
# Do NOT copy them here to avoid baking secrets into the image.

EXPOSE 8080

CMD ["./main"]
