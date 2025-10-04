package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO: Load books.json on startup
	// TODO: Set up API routes
	// TODO: Set up static file server for frontend
	// TODO: Start HTTP server on :8080
	
	fmt.Println("Reading Tracker API")
	fmt.Println("Server starting on http://localhost:8080")
	
	// Placeholder - will be implemented in T029
	log.Fatal(http.ListenAndServe(":8080", nil))
}
