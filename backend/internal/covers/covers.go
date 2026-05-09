package covers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// PlaceholderMaxBytes is the size of Open Library's 1×1 pixel placeholder image
	PlaceholderMaxBytes = 100

	openLibraryCoversByISBN = "https://covers.openlibrary.org/b/isbn/%s-M.jpg"
	openLibraryCoversByID   = "https://covers.openlibrary.org/b/id/%d-M.jpg"
	openLibrarySearchAPI    = "https://openlibrary.org/search.json"
)

// HTTPClient allows injecting a mock HTTP client for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Resolver resolves book cover URLs using a fallback chain:
// 1. ISBN-based Open Library cover
// 2. Title+Author search via Open Library Search API
type Resolver struct {
	client HTTPClient
}

// New creates a new cover Resolver with the given HTTP client.
// If client is nil, a default client with a 10-second timeout is used.
func New(client HTTPClient) *Resolver {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Resolver{client: client}
}

// CoverURLByISBN returns the Open Library cover URL for an ISBN
func CoverURLByISBN(isbn string) string {
	if isbn == "" {
		return ""
	}
	return fmt.Sprintf(openLibraryCoversByISBN, isbn)
}

// searchResult represents relevant fields from the Open Library Search API
type searchResult struct {
	Docs []searchDoc `json:"docs"`
}

type searchDoc struct {
	CoverI    int    `json:"cover_i"`
	Title     string `json:"title"`
	AuthorKey []string `json:"author_key"`
}

// IsValidCover checks if a cover URL returns a real image (not a placeholder).
// Open Library returns a ~43-byte 1×1 transparent GIF for missing covers.
func (r *Resolver) IsValidCover(coverURL string) bool {
	if coverURL == "" {
		return false
	}

	req, err := http.NewRequest("GET", coverURL, nil)
	if err != nil {
		return false
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	// Read the response to check size
	body, err := io.ReadAll(io.LimitReader(resp.Body, PlaceholderMaxBytes+1))
	if err != nil {
		return false
	}

	return len(body) > PlaceholderMaxBytes
}

// SearchByTitleAuthor searches Open Library for a book by title and author,
// returning a cover URL if found. Uses the cover_i field from search results
// for more accurate cover matching than ISBN-based lookups.
func (r *Resolver) SearchByTitleAuthor(title, author string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title is required")
	}

	params := url.Values{}
	params.Set("title", title)
	if author != "" {
		params.Set("author", author)
	}
	params.Set("fields", "cover_i,title,author_key")
	params.Set("limit", "3")

	reqURL := fmt.Sprintf("%s?%s", openLibrarySearchAPI, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search returned status %d", resp.StatusCode)
	}

	var result searchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode search response: %w", err)
	}

	// Find first result with a cover
	for _, doc := range result.Docs {
		if doc.CoverI > 0 {
			return fmt.Sprintf(openLibraryCoversByID, doc.CoverI), nil
		}
	}

	return "", nil
}

// Resolve attempts to find a valid cover URL using the fallback chain:
// 1. If existingURL is set and valid, keep it
// 2. Try ISBN-based cover URL
// 3. Search by title+author
func (r *Resolver) Resolve(existingURL, isbn, title, author string) string {
	// 1. Check existing URL
	if existingURL != "" && r.IsValidCover(existingURL) {
		return existingURL
	}

	// 2. Try ISBN-based cover
	if isbn != "" {
		isbnURL := CoverURLByISBN(isbn)
		if r.IsValidCover(isbnURL) {
			return isbnURL
		}
	}

	// 3. Fallback to title+author search
	searchURL, err := r.SearchByTitleAuthor(title, author)
	if err == nil && searchURL != "" {
		return searchURL
	}

	return ""
}

// PreferISBN returns the best ISBN from isbn13 and isbn10 values.
// Prefers ISBN-13 over ISBN-10.
func PreferISBN(isbn, isbn13 string) string {
	isbn13 = strings.TrimSpace(isbn13)
	isbn = strings.TrimSpace(isbn)

	if isbn13 != "" && isbn13 != "0" {
		return isbn13
	}
	if isbn != "" && isbn != "0" {
		return isbn
	}
	return ""
}
