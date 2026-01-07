package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheckHandler(t *testing.T) {
	config := Config{
		Port:            "8080",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		RateLimit:       100,
	}
	server := NewServer(config)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.HealthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%v'", response["status"])
	}

	if response["service"] != "FrostByte-API" {
		t.Errorf("expected service 'FrostByte-API', got '%v'", response["service"])
	}
}

func TestProcessTemperatureHandler(t *testing.T) {
	config := Config{
		Port:            "8080",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		RateLimit:       100,
	}
	server := NewServer(config)

	tests := []struct {
		name           string
		data           TemperatureData
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid celsius temperature",
			data: TemperatureData{
				Location:    "Antarctica",
				Temperature: -15.5,
				Unit:        "celsius",
				Timestamp:   time.Now(),
				Conditions:  "clear",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "valid fahrenheit temperature",
			data: TemperatureData{
				Location:    "Alaska",
				Temperature: 32.0,
				Unit:        "fahrenheit",
				Timestamp:   time.Now(),
				Conditions:  "snowy",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "missing location",
			data: TemperatureData{
				Temperature: -10.0,
				Unit:        "celsius",
				Timestamp:   time.Now(),
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid unit",
			data: TemperatureData{
				Location:    "Iceland",
				Temperature: 100.0,
				Unit:        "invalid",
				Timestamp:   time.Now(),
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.data)
			if err != nil {
				t.Fatalf("failed to marshal test data: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/temperature/process", bytes.NewReader(body))
			w := httptest.NewRecorder()

			server.ProcessTemperatureHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError && w.Code == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response["location"] != tt.data.Location {
					t.Errorf("expected location '%s', got '%v'", tt.data.Location, response["location"])
				}
			}
		})
	}
}

func TestGetTemperatureStatsHandler(t *testing.T) {
	config := Config{
		Port:            "8080",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		RateLimit:       100,
	}
	server := NewServer(config)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/temperature/stats?location=Alaska", nil)
	w := httptest.NewRecorder()

	server.GetTemperatureStatsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["location"] != "Alaska" {
		t.Errorf("expected location 'Alaska', got '%v'", response["location"])
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	config := Config{
		Port:            "8080",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		RateLimit:       100,
	}
	server := NewServer(config)

	handler := server.SecurityHeadersMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Content-Security-Policy":   "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, expectedValue := range expectedHeaders {
		if w.Header().Get(header) != expectedValue {
			t.Errorf("expected header %s to be '%s', got '%s'", header, expectedValue, w.Header().Get(header))
		}
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(2)

	// First two requests should pass
	if !limiter.Allow("127.0.0.1") {
		t.Error("expected first request to be allowed")
	}

	if !limiter.Allow("127.0.0.1") {
		t.Error("expected second request to be allowed")
	}

	// Third request should be blocked
	if limiter.Allow("127.0.0.1") {
		t.Error("expected third request to be blocked")
	}

	// Different IP should be allowed
	if !limiter.Allow("127.0.0.2") {
		t.Error("expected request from different IP to be allowed")
	}
}

func TestTemperatureConversion(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		unit        string
		expectedC   float64
	}{
		{
			name:        "celsius to celsius",
			temperature: -10.0,
			unit:        "celsius",
			expectedC:   -10.0,
		},
		{
			name:        "fahrenheit to celsius",
			temperature: 32.0,
			unit:        "fahrenheit",
			expectedC:   0.0,
		},
		{
			name:        "kelvin to celsius",
			temperature: 273.15,
			unit:        "kelvin",
			expectedC:   0.0,
		},
	}

	config := Config{
		Port:            "8080",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		RateLimit:       100,
	}
	server := NewServer(config)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := TemperatureData{
				Location:    "Test",
				Temperature: tt.temperature,
				Unit:        tt.unit,
				Timestamp:   time.Now(),
				Conditions:  "test",
			}

			body, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("failed to marshal test data: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/temperature/process", bytes.NewReader(body))
			w := httptest.NewRecorder()

			server.ProcessTemperatureHandler(w, req)

			var response map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			tempC := response["temperature_c"].(float64)
			if tempC < tt.expectedC-0.1 || tempC > tt.expectedC+0.1 {
				t.Errorf("expected temperature %.2f°C, got %.2f°C", tt.expectedC, tempC)
			}
		})
	}
}
