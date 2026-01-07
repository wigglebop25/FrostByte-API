package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Config holds the application configuration
type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RateLimit       int
}

// Server wraps the HTTP server with middleware
type Server struct {
	config      Config
	rateLimiter *RateLimiter
	logger      *Logger
}

// Logger provides structured logging
type Logger struct {
	mu sync.Mutex
}

// Info logs informational messages
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["level"] = "info"
	fields["msg"] = msg
	fields["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	jsonData, err := json.Marshal(fields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling log: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// Error logs error messages
func (l *Logger) Error(msg string, err error, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["level"] = "error"
	fields["msg"] = msg
	fields["error"] = err.Error()
	fields["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	jsonData, marshalErr := json.Marshal(fields)
	if marshalErr != nil {
		fmt.Fprintf(os.Stderr, "error marshaling log: %v\n", marshalErr)
		return
	}
	fmt.Println(string(jsonData))
}

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientBucket
	limit   int
}

type clientBucket struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientBucket),
		limit:   limit,
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[ip]

	if !exists || now.Sub(bucket.lastRefill) > time.Minute {
		rl.clients[ip] = &clientBucket{
			tokens:     rl.limit - 1,
			lastRefill: now,
		}
		return true
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// NewServer creates a new server instance
func NewServer(config Config) *Server {
	return &Server{
		config:      config,
		rateLimiter: NewRateLimiter(config.RateLimit),
		logger:      &Logger{},
	}
}

// SecurityHeadersMiddleware adds security headers to all responses
func (s *Server) SecurityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next(w, r)
	}
}

// RateLimitMiddleware applies rate limiting
func (s *Server) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !s.rateLimiter.Allow(ip) {
			s.logger.Info("rate limit exceeded", map[string]interface{}{
				"ip":     ip,
				"path":   r.URL.Path,
				"method": r.Method,
			})
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

// LoggingMiddleware logs all requests
func (s *Server) LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		s.logger.Info("request", map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start).Milliseconds(),
			"ip":       r.RemoteAddr,
		})
	}
}

// HealthCheckHandler handles health check requests
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "FrostByte-API",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode response", err, nil)
	}
}

// TemperatureData represents temperature readings
type TemperatureData struct {
	Location    string    `json:"location"`
	Temperature float64   `json:"temperature"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Conditions  string    `json:"conditions"`
}

// ProcessTemperatureHandler processes temperature data
func (s *Server) ProcessTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data TemperatureData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if data.Location == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	if data.Unit != "celsius" && data.Unit != "fahrenheit" && data.Unit != "kelvin" {
		http.Error(w, "Invalid unit. Must be celsius, fahrenheit, or kelvin", http.StatusBadRequest)
		return
	}

	// Process temperature data (convert to Celsius for standardization)
	startProcess := time.Now()
	celsius := data.Temperature
	if data.Unit == "fahrenheit" {
		celsius = (data.Temperature - 32) * 5 / 9
	} else if data.Unit == "kelvin" {
		celsius = data.Temperature - 273.15
	}

	// Determine cold weather conditions
	isCold := celsius < 0
	severity := "normal"
	if celsius < -10 {
		severity = "extreme"
	} else if celsius < 0 {
		severity = "cold"
	}

	processingTime := time.Since(startProcess).Milliseconds()

	response := map[string]interface{}{
		"processed_at":       time.Now().UTC().Format(time.RFC3339),
		"location":           data.Location,
		"temperature_c":      celsius,
		"original_temp":      data.Temperature,
		"original_unit":      data.Unit,
		"is_cold":            isCold,
		"severity":           severity,
		"conditions":         data.Conditions,
		"processing_time_ms": processingTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode response", err, nil)
	}
}

// GetTemperatureStatsHandler returns temperature statistics
func (s *Server) GetTemperatureStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	location := r.URL.Query().Get("location")

	// In a real implementation, this would query a database
	// For now, return mock statistics
	response := map[string]interface{}{
		"location":     location,
		"avg_temp_c":   -5.2,
		"min_temp_c":   -15.8,
		"max_temp_c":   2.3,
		"samples":      1000,
		"last_updated": time.Now().UTC().Format(time.RFC3339),
		"cold_days":    45,
		"extreme_days": 12,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode response", err, nil)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register routes with middleware chain
	mux.HandleFunc("/health",
		s.SecurityHeadersMiddleware(
			s.LoggingMiddleware(
				s.HealthCheckHandler)))

	mux.HandleFunc("/api/v1/temperature/process",
		s.SecurityHeadersMiddleware(
			s.RateLimitMiddleware(
				s.LoggingMiddleware(
					s.ProcessTemperatureHandler))))

	mux.HandleFunc("/api/v1/temperature/stats",
		s.SecurityHeadersMiddleware(
			s.RateLimitMiddleware(
				s.LoggingMiddleware(
					s.GetTemperatureStatsHandler))))

	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      mux,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		s.logger.Info("starting server", map[string]interface{}{
			"port": s.config.Port,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("server error", err, nil)
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("shutting down server", nil)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("server forced to shutdown", err, nil)
		return err
	}

	s.logger.Info("server stopped", nil)
	return nil
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rateLimit := 100
	if rl := os.Getenv("RATE_LIMIT"); rl != "" {
		if parsed, err := strconv.Atoi(rl); err == nil {
			rateLimit = parsed
		}
	}

	return Config{
		Port:            port,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		RateLimit:       rateLimit,
	}
}

func main() {
	config := LoadConfig()
	server := NewServer(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
