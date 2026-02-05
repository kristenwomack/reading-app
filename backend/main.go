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

// corsMiddleware adds CORS headers to allow browser access
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
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
	
	// Determine paths based on environment
	// DATA_DIR is set in container; locally we use parent directory
	dataDir := os.Getenv("DATA_DIR")
	var dbPath, booksPath, frontendDir string
	if dataDir != "" {
		// Running in container
		dbPath = filepath.Join(dataDir, "books.db")
		booksPath = "books.json" // In container root
		frontendDir = "frontend"
	} else {
		// Running locally from backend directory
		dbPath = filepath.Join("..", "books.db")
		booksPath = filepath.Join("..", "books.json")
		frontendDir = "../frontend"
	}
	
	// Initialize SQLite database
	dataStore, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dataStore.Close()
	handlers.SetStore(dataStore)
	handlers.SetFrontendDir(frontendDir)
	
	// Check if we need to import from books.json
	count, _ := dataStore.BookCount()
	if count == 0 {
		fmt.Println("Database empty, importing from books.json...")
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
	
	// Serve frontend static files
	fs := http.FileServer(http.Dir(frontendDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve admin.html for /admin path
		if r.URL.Path == "/admin" || strings.HasPrefix(r.URL.Path, "/admin/") {
			http.ServeFile(w, r, filepath.Join(frontendDir, "admin.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
	
	// Get port from environment (default 3000 for local, 8080 in container)
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
