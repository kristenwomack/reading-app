package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kristenwomack/reading-app/backend/internal/covers"
)

// EnrichCoversResponse is the response from the cover enrichment endpoint
type EnrichCoversResponse struct {
	Total    int `json:"total"`
	Updated  int `json:"updated"`
	Skipped  int `json:"skipped"`
	Failed   int `json:"failed"`
}

// EnrichCovers scans all books and resolves missing or invalid cover URLs
// using the fallback chain: ISBN → title+author search.
// This is a protected endpoint (requires auth) since it makes many external API calls.
func EnrichCovers(w http.ResponseWriter, r *http.Request) {
	if dataStore == nil {
		http.Error(w, "Store not initialized", http.StatusInternalServerError)
		return
	}

	allBooks, err := dataStore.GetAllBooks()
	if err != nil {
		http.Error(w, "Failed to load books", http.StatusInternalServerError)
		return
	}

	resp := EnrichCoversResponse{Total: len(allBooks)}

	for _, book := range allBooks {
		isbn := covers.PreferISBN(book.ISBN, book.ISBN13)

		// Skip books that already have a non-ISBN-based cover URL
		// (these were likely manually set or from a search)
		if book.CoverURL != "" && !isISBNCoverURL(book.CoverURL) {
			resp.Skipped++
			continue
		}

		// Use the resolver to find the best cover
		newURL := coverResolver.Resolve(book.CoverURL, isbn, book.Title, book.Author)

		if newURL == book.CoverURL {
			resp.Skipped++
			continue
		}

		if newURL != "" {
			book.CoverURL = newURL
			if err := dataStore.UpdateBook(&book); err != nil {
				log.Printf("Failed to update cover for %q: %v", book.Title, err)
				resp.Failed++
				continue
			}
			resp.Updated++
			log.Printf("Updated cover for %q → %s", book.Title, newURL)
		} else {
			resp.Skipped++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// isISBNCoverURL checks if a cover URL is an ISBN-based Open Library URL
func isISBNCoverURL(url string) bool {
	return len(url) > 0 && 
		(contains(url, "covers.openlibrary.org/b/isbn/"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetCoverURL resolves a cover URL for a single book by ID.
// Uses the full fallback chain and updates the database if a better cover is found.
func GetCoverURL(w http.ResponseWriter, r *http.Request) {
	if dataStore == nil {
		http.Error(w, "Store not initialized", http.StatusInternalServerError)
		return
	}

	// Extract book ID from path: /api/books/{id}/cover
	path := r.URL.Path
	var bookID int64
	_, err := fmt.Sscanf(path, "/api/books/%d/cover", &bookID)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := dataStore.GetBook(bookID)
	if err != nil {
		http.Error(w, "Failed to load book", http.StatusInternalServerError)
		return
	}
	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	isbn := covers.PreferISBN(book.ISBN, book.ISBN13)
	newURL := coverResolver.Resolve(book.CoverURL, isbn, book.Title, book.Author)

	// Update if we found a better cover
	if newURL != "" && newURL != book.CoverURL {
		book.CoverURL = newURL
		dataStore.UpdateBook(book)
	}

	response := map[string]string{
		"coverUrl": newURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
