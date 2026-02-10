package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestGenerateToken verifies token generation
func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}
}

// TestValidateToken verifies token validation
func TestValidateToken(t *testing.T) {
	// Generate a token
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	err = ValidateToken(token)
	if err != nil {
		t.Errorf("Expected valid token, got error: %v", err)
	}
}

// TestValidateTokenInvalid verifies rejecting invalid tokens
func TestValidateTokenInvalid(t *testing.T) {
	// Try to validate an invalid token
	err := ValidateToken("invalid.token.string")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

// TestValidateTokenEmpty verifies rejecting empty tokens
func TestValidateTokenEmpty(t *testing.T) {
	// Try to validate an empty token
	err := ValidateToken("")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

// TestCheckPassword verifies password checking
func TestCheckPassword(t *testing.T) {
	// Set password
	os.Setenv("READING_APP_PASSWORD", "correctpassword")
	defer os.Unsetenv("READING_APP_PASSWORD")

	// Check correct password
	err := CheckPassword("correctpassword")
	if err != nil {
		t.Errorf("Expected no error for correct password, got: %v", err)
	}
}

// TestCheckPasswordWrong verifies rejecting wrong password
func TestCheckPasswordWrong(t *testing.T) {
	// Set password
	os.Setenv("READING_APP_PASSWORD", "correctpassword")
	defer os.Unsetenv("READING_APP_PASSWORD")

	// Check wrong password
	err := CheckPassword("wrongpassword")
	if err != ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got: %v", err)
	}
}

// TestCheckPasswordMissing verifies error when no password is set
func TestCheckPasswordMissing(t *testing.T) {
	// Ensure no password is set
	os.Unsetenv("READING_APP_PASSWORD")

	// Check password
	err := CheckPassword("anypassword")
	if err != ErrNoPassword {
		t.Errorf("Expected ErrNoPassword, got: %v", err)
	}
}

// TestGetPasswordHash verifies getting password hash
func TestGetPasswordHash(t *testing.T) {
	// Set password
	os.Setenv("READING_APP_PASSWORD", "testpassword")
	defer os.Unsetenv("READING_APP_PASSWORD")

	// Get hash
	hash, err := GetPasswordHash()
	if err != nil {
		t.Fatalf("Failed to get password hash: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

// TestGetPasswordHashMissing verifies error when no password is set
func TestGetPasswordHashMissing(t *testing.T) {
	// Ensure no password is set
	os.Unsetenv("READING_APP_PASSWORD")

	// Try to get hash
	_, err := GetPasswordHash()
	if err != ErrNoPassword {
		t.Errorf("Expected ErrNoPassword, got: %v", err)
	}
}

// TestSetAuthCookie verifies setting auth cookie
func TestSetAuthCookie(t *testing.T) {
	// Create response recorder
	w := httptest.NewRecorder()

	// Set cookie
	token := "test.jwt.token"
	SetAuthCookie(w, token)

	// Verify cookie was set
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookie to be set")
	}

	cookie := cookies[0]
	if cookie.Name != "auth_token" {
		t.Errorf("Expected cookie name 'auth_token', got %q", cookie.Name)
	}

	if cookie.Value != token {
		t.Errorf("Expected cookie value %q, got %q", token, cookie.Value)
	}

	if !cookie.HttpOnly {
		t.Error("Expected cookie to be HttpOnly")
	}

	if cookie.MaxAge != 30*24*60*60 {
		t.Errorf("Expected MaxAge 30 days, got %d", cookie.MaxAge)
	}
}

// TestClearAuthCookie verifies clearing auth cookie
func TestClearAuthCookie(t *testing.T) {
	// Create response recorder
	w := httptest.NewRecorder()

	// Clear cookie
	ClearAuthCookie(w)

	// Verify cookie was cleared
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookie to be set for clearing")
	}

	cookie := cookies[0]
	if cookie.Name != "auth_token" {
		t.Errorf("Expected cookie name 'auth_token', got %q", cookie.Name)
	}

	if cookie.Value != "" {
		t.Errorf("Expected empty cookie value, got %q", cookie.Value)
	}

	if cookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge -1, got %d", cookie.MaxAge)
	}
}

// TestGetTokenFromRequest verifies extracting token from request
func TestGetTokenFromRequest(t *testing.T) {
	// Create request with cookie
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "test.token.value",
	})

	// Get token
	token := GetTokenFromRequest(req)
	if token != "test.token.value" {
		t.Errorf("Expected token %q, got %q", "test.token.value", token)
	}
}

// TestGetTokenFromRequestMissing verifies handling missing token
func TestGetTokenFromRequestMissing(t *testing.T) {
	// Create request without cookie
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Get token
	token := GetTokenFromRequest(req)
	if token != "" {
		t.Errorf("Expected empty token, got %q", token)
	}
}

// TestIsAuthenticated verifies checking authentication
func TestIsAuthenticated(t *testing.T) {
	// Generate a valid token
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create request with valid token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	// Check authentication
	if !IsAuthenticated(req) {
		t.Error("Expected request to be authenticated")
	}
}

// TestIsAuthenticatedInvalid verifies rejecting invalid authentication
func TestIsAuthenticatedInvalid(t *testing.T) {
	// Create request with invalid token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "invalid.token",
	})

	// Check authentication
	if IsAuthenticated(req) {
		t.Error("Expected request to not be authenticated")
	}
}

// TestIsAuthenticatedMissing verifies rejecting missing authentication
func TestIsAuthenticatedMissing(t *testing.T) {
	// Create request without token
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Check authentication
	if IsAuthenticated(req) {
		t.Error("Expected request to not be authenticated")
	}
}

// TestGenerateCSRFToken verifies CSRF token generation
func TestGenerateCSRFToken(t *testing.T) {
	token := GenerateCSRFToken()
	if token == "" {
		t.Error("Expected non-empty CSRF token")
	}

	// Generate another token and ensure they're different
	token2 := GenerateCSRFToken()
	if token == token2 {
		t.Error("Expected different CSRF tokens")
	}
}
