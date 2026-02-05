package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrNoPassword      = errors.New("READING_APP_PASSWORD environment variable not set")
	jwtSecret          []byte
)

func init() {
	// Generate a random JWT secret on startup
	jwtSecret = make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		panic("failed to generate JWT secret: " + err.Error())
	}
}

// Claims represents JWT claims
type Claims struct {
	jwt.RegisteredClaims
}

// GetPasswordHash returns the bcrypt hash of the configured password
func GetPasswordHash() (string, error) {
	password := os.Getenv("READING_APP_PASSWORD")
	if password == "" {
		return "", ErrNoPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword verifies the provided password against the environment variable
func CheckPassword(password string) error {
	expected := os.Getenv("READING_APP_PASSWORD")
	if expected == "" {
		return ErrNoPassword
	}
	if password != expected {
		return ErrInvalidPassword
	}
	return nil
}

// GenerateToken creates a new JWT token valid for 30 days
func GenerateToken() (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "reading-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken checks if a JWT token is valid
func ValidateToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// SetAuthCookie sets the JWT token as an httpOnly cookie
func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60, // 30 days
	})
}

// ClearAuthCookie removes the auth cookie
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// GetTokenFromRequest extracts the auth token from request cookies
func GetTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// IsAuthenticated checks if the request has a valid auth token
func IsAuthenticated(r *http.Request) bool {
	token := GetTokenFromRequest(r)
	if token == "" {
		return false
	}
	return ValidateToken(token) == nil
}

// GenerateCSRFToken generates a random CSRF token
func GenerateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
