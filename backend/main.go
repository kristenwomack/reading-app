package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kristenwomack/reading-app/backend/internal/books"
	"github.com/kristenwomack/reading-app/backend/internal/handlers"
	"github.com/kristenwomack/reading-app/backend/internal/store"
)

// allowedOrigins is a map of allowed CORS origins for O(1) lookup
var allowedOrigins map[string]bool

// initAllowedOrigins parses the ALLOWED_ORIGINS environment variable once at startup
func initAllowedOrigins() {
	allowedOrigins = make(map[string]bool)
	
	// Get allowed origins from environment variable, default to localhost:3000
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		allowedOriginsStr = "http://localhost:3000"
	}
	
	// Parse comma-separated origins and trim whitespace
	origins := strings.Split(allowedOriginsStr, ",")
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowedOrigins[trimmed] = true
		}
	}
}

// corsMiddleware adds CORS headers to allow browser access from configured origins
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// If no Origin header, this is a same-origin request - allow it
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Check if the request origin is allowed using O(1) map lookup
		if !allowedOrigins[origin] {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}
		
		// Set CORS headers for allowed origin
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Reading Tracker API")
	
	// Initialize allowed CORS origins from environment
	initAllowedOrigins()
	
	// Determine database path (configurable for Railway persistent volume)
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("..", "books.db")
	}
	
	// Initialize SQLite database
	dataStore, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dataStore.Close()
	handlers.SetStore(dataStore)
	
	// Check if we need to import from books.json
	count, _ := dataStore.BookCount()
	if count == 0 {
		fmt.Println("Database empty, importing from books.json...")
		// Try local path first (Docker), then relative path (dev)
		booksPath := "books.json"
		if _, err := os.Stat(booksPath); os.IsNotExist(err) {
			booksPath = filepath.Join("..", "books.json")
		}
		allBooks, err := books.LoadBooks(booksPath)
		if err != nil {
			log.Fatalf("Failed to load books.json: %v", err)
		}
		imported, err := dataStore.ImportFromJSON(allBooks)
		if err != nil {
			log.Fatalf("Failed to import books: %v", err)
		}
		fmt.Printf("Imported %d books into database\n", imported)
	} else {
		fmt.Printf("Database contains %d books\n", count)
	}
	
	// Check if password is configured
	if os.Getenv("READING_APP_PASSWORD") == "" {
		fmt.Println("WARNING: READING_APP_PASSWORD not set - admin features disabled")
	} else {
		fmt.Println("Admin authentication enabled")
	}
	
	// Set up API routes (public)
	http.HandleFunc("/api/years", handlers.GetYears)
	http.HandleFunc("/api/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AuthMiddleware(handlers.CreateBook)(w, r)
		} else {
			handlers.GetBooks(w, r)
		}
	})
	http.HandleFunc("/api/books/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handlers.AuthMiddleware(handlers.UpdateBook)(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.AuthMiddleware(handlers.DeleteBook)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/stats", handlers.GetStats)
	
	// Goals routes
	http.HandleFunc("/api/goals/", handlers.GetGoal)
	http.HandleFunc("/api/goals", handlers.AuthMiddleware(handlers.SetGoal))
	
	// Auth routes
	http.HandleFunc("/api/auth/login", handlers.Login)
	http.HandleFunc("/api/auth/logout", handlers.Logout)
	http.HandleFunc("/api/auth/check", handlers.CheckAuth)
	
	// Export route (protected)
	http.HandleFunc("/api/export", handlers.AuthMiddleware(handlers.ExportBooks))
	
	// Health check endpoint
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	
	// Serve frontend static files (try Docker path first, then dev path)
	frontendDir := "frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "../frontend"
	}
	fs := http.FileServer(http.Dir(frontendDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve admin.html for /admin path
		if r.URL.Path == "/admin" || strings.HasPrefix(r.URL.Path, "/admin/") {
			http.ServeFile(w, r, filepath.Join(frontendDir, "admin.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
	
	// Determine port (Railway sets PORT env var)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	fmt.Printf("Server starting on http://localhost:%s\n", port)
	fmt.Println("API endpoints:")
	fmt.Println("  GET  /api/years")
	fmt.Println("  GET  /api/books?year=2025")
	fmt.Println("  POST /api/books (auth required)")
	fmt.Println("  PUT  /api/books/:id (auth required)")
	fmt.Println("  DELETE /api/books/:id (auth required)")
	fmt.Println("  GET  /api/stats?year=2025")
	fmt.Println("  POST /api/auth/login")
	fmt.Println("  GET  /admin (book entry form)")
	
	// Wrap with CORS middleware
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}
