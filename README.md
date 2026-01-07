# FrostByte-API

A high-performance, minimalist backend API built in Go, designed for blistering speed and reliability in cold-weather data processing scenarios. Built with the Go philosophy of simplicity and efficiency with security-first architecture.

## Features

- **High Performance**: Lightweight HTTP server with minimal overhead
- **Security-First**: Built-in security headers, rate limiting, and input validation
- **Structured Logging**: JSON-formatted logs for easy parsing and monitoring
- **Graceful Shutdown**: Proper cleanup of resources on termination
- **Cold Weather Processing**: Specialized endpoints for temperature data processing
- **Rate Limiting**: Token bucket algorithm to prevent abuse
- **Health Checks**: Ready-to-use health endpoint for monitoring
- **Docker Support**: Containerized deployment with multi-stage builds

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

### Installation

```bash
# Clone the repository
git clone https://github.com/wigglebop25/FrostByte-API.git
cd FrostByte-API

# Install dependencies
go mod download
```

### Running the API

```bash
# Run directly
go run .

# Or use Make
make run

# Or build and run
make build
./frostbyte-api
```

### Using Docker

```bash
# Build the Docker image
make docker-build

# Run the container
make docker-run
```

## API Endpoints

### Health Check

```bash
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2026-01-07T16:22:58Z",
  "service": "FrostByte-API"
}
```

### Process Temperature Data

```bash
POST /api/v1/temperature/process
Content-Type: application/json

{
  "location": "Antarctica",
  "temperature": -15.5,
  "unit": "celsius",
  "timestamp": "2026-01-07T16:22:58Z",
  "conditions": "clear"
}
```

Response:
```json
{
  "processed_at": "2026-01-07T16:22:58Z",
  "location": "Antarctica",
  "temperature_c": -15.5,
  "original_temp": -15.5,
  "original_unit": "celsius",
  "is_cold": true,
  "severity": "extreme",
  "conditions": "clear",
  "processing_time_ms": 0
}
```

Supported units: `celsius`, `fahrenheit`, `kelvin`

Severity levels:
- `normal`: Temperature >= 0°C
- `cold`: Temperature < 0°C
- `extreme`: Temperature < -10°C

### Get Temperature Statistics

```bash
GET /api/v1/temperature/stats?location=Alaska
```

Response:
```json
{
  "location": "Alaska",
  "avg_temp_c": -5.2,
  "min_temp_c": -15.8,
  "max_temp_c": 2.3,
  "samples": 1000,
  "last_updated": "2026-01-07T16:22:58Z",
  "cold_days": 45,
  "extreme_days": 12
}
```

## Configuration

The API can be configured using environment variables:

- `PORT`: Server port (default: `8080`)
- `RATE_LIMIT`: Requests per minute per IP (default: `100`)

Example:
```bash
PORT=3000 RATE_LIMIT=50 go run .
```

## Security Features

### Security Headers

All responses include the following security headers:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy: default-src 'self'`
- `Strict-Transport-Security: max-age=31536000; includeSubDomains`

### Rate Limiting

Token bucket rate limiting is applied per IP address to prevent abuse. Default: 100 requests per minute.

### Input Validation

All input is validated and sanitized to prevent injection attacks.

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage
```

### Linting

```bash
# Format code and run vet
make lint
```

### Build

```bash
# Build binary
make build

# Clean build artifacts
make clean

# Run all checks and build
make all
```

## Architecture

### Design Principles

1. **Simplicity**: Minimal dependencies, standard library first
2. **Performance**: Optimized for speed with efficient algorithms
3. **Security**: Defense-in-depth with multiple security layers
4. **Reliability**: Graceful error handling and shutdown
5. **Observability**: Structured logging for monitoring

### Components

- **HTTP Server**: Standard library `net/http` with custom middleware
- **Rate Limiter**: Token bucket algorithm for request throttling
- **Logger**: Structured JSON logging for production use
- **Temperature Processor**: Core business logic for cold-weather data

## Performance

The API is designed for high performance:

- Zero-allocation middleware chains where possible
- Efficient JSON encoding/decoding
- Minimal memory footprint
- Fast startup time
- Low latency response times

## Deployment

### Docker Deployment

The included Dockerfile uses multi-stage builds for minimal image size:

```bash
docker build -t frostbyte-api .
docker run -p 8080:8080 frostbyte-api
```

### Production Recommendations

1. Use HTTPS in production (configure a reverse proxy like nginx)
2. Set appropriate rate limits based on your needs
3. Monitor logs for security events
4. Use container orchestration (Kubernetes, Docker Swarm) for scaling
5. Configure health checks for load balancers

## License

BSD 3-Clause License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please ensure:
- Code follows Go best practices
- Tests are included for new features
- Security considerations are addressed
- Documentation is updated