package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCorsMiddlewareDefaultOrigin(t *testing.T) {
	// Test default behavior when no ALLOWED_ORIGINS is set
	os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	
	handler.ServeHTTP(rec, req)
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be http://localhost:3000, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestCorsMiddlewareDefaultOriginRejection(t *testing.T) {
	// Test that non-default origins are rejected when no ALLOWED_ORIGINS is set
	os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec := httptest.NewRecorder()
	
	handler.ServeHTTP(rec, req)
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected no Access-Control-Allow-Origin header, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

func TestCorsMiddlewareConfiguredOrigins(t *testing.T) {
	// Test with multiple configured origins
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,https://reading-app-production-5700.up.railway.app")
	defer os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	
	// Test first allowed origin
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be http://localhost:3000, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	// Test second allowed origin
	req = httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "https://reading-app-production-5700.up.railway.app")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://reading-app-production-5700.up.railway.app" {
		t.Errorf("Expected Access-Control-Allow-Origin to be https://reading-app-production-5700.up.railway.app, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestCorsMiddlewareRejectUnknownOrigin(t *testing.T) {
	// Test that unknown origins are rejected
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,https://reading-app-production-5700.up.railway.app")
	defer os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec := httptest.NewRecorder()
	
	handler.ServeHTTP(rec, req)
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected no Access-Control-Allow-Origin header, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

func TestCorsMiddlewareNoOriginHeader(t *testing.T) {
	// Test that requests without Origin header are allowed (same-origin requests)
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	// No Origin header set
	rec := httptest.NewRecorder()
	
	handler.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 for request without Origin header, got %d", rec.Code)
	}
	// No CORS headers should be set for same-origin requests
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected no Access-Control-Allow-Origin header for same-origin request, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCorsMiddlewarePreflight(t *testing.T) {
	// Test preflight OPTIONS request
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS requests")
	}))
	
	req := httptest.NewRequest("OPTIONS", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	
	handler.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS request, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin header for preflight, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected Access-Control-Allow-Methods header for preflight")
	}
}

func TestCorsMiddlewareWhitespaceInOrigins(t *testing.T) {
	// Test that whitespace in ALLOWED_ORIGINS is handled correctly
	os.Setenv("ALLOWED_ORIGINS", " http://localhost:3000 , https://example.com ")
	defer os.Unsetenv("ALLOWED_ORIGINS")
	initAllowedOrigins()
	
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	
	handler.ServeHTTP(rec, req)
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be http://localhost:3000, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}
