package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kristenwomack/reading-app/backend/internal/books"
	"github.com/kristenwomack/reading-app/backend/internal/handlers"
)

// corsMiddleware adds CORS headers to allow browser access
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
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
	fmt.Println("Loading books.json...")
	
	// Load books.json from repository root
	booksPath := filepath.Join("..", "books.json")
	allBooks, err := books.LoadBooks(booksPath)
	if err != nil {
		log.Fatalf("Failed to load books: %v", err)
	}
	
	fmt.Printf("Loaded %d books\n", len(allBooks))
	handlers.SetBooks(allBooks)
	
	// Set up API routes
	http.HandleFunc("/api/years", handlers.GetYears)
	http.HandleFunc("/api/books", handlers.GetBooks)
	http.HandleFunc("/api/stats", handlers.GetStats)
	
	// Serve frontend static files
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)
	
	fmt.Println("Server starting on http://localhost:3000")
	fmt.Println("API endpoints:")
	fmt.Println("  GET /api/years")
	fmt.Println("  GET /api/books?year=2025")
	fmt.Println("  GET /api/stats?year=2025")
	
	// Wrap with CORS middleware
	log.Fatal(http.ListenAndServe(":3000", corsMiddleware(http.DefaultServeMux)))
}
