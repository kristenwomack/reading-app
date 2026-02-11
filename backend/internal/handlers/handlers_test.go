package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kristenwomack/reading-app/backend/internal/books"
)

// setupTestBooks creates test data for handler tests
func setupTestBooks() []books.Book {
	return []books.Book{
		{
			Title:    "Test Book 2025",
			Author:   "Author One",
			DateRead: "2025/01/15",
			Pages:    200,
			Shelf:    "read",
			ISBN13:   "9781234567890",
		},
		{
			Title:    "Another Book 2025",
			Author:   "Author Two",
			DateRead: "2025/02/20",
			Pages:    300,
			Shelf:    "read",
			ISBN13:   "9780987654321",
		},
		{
			Title:    "Book 2024",
			Author:   "Author Three",
			DateRead: "2024/12/10",
			Pages:    150,
			Shelf:    "read",
			ISBN13:   "9781111111111",
		},
		{
			Title:    "Currently Reading",
			Author:   "Author Four",
			DateRead: "",
			Pages:    400,
			Shelf:    "currently-reading",
		},
	}
}

// TestGetYears verifies the GetYears handler
func TestGetYears(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/years", nil)
	w := httptest.NewRecorder()

	// Execute
	GetYears(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify content type
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	years, ok := response["years"].([]interface{})
	if !ok {
		t.Fatal("Expected years array in response")
	}

	// Should have 2 years (2024 and 2025) - currently-reading doesn't count
	if len(years) != 2 {
		t.Errorf("Expected 2 years, got %d", len(years))
	}
}

// TestGetYearsEmpty verifies GetYears with no books
func TestGetYearsEmpty(t *testing.T) {
	// Setup with empty books
	SetBooks([]books.Book{})

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/years", nil)
	w := httptest.NewRecorder()

	// Execute
	GetYears(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	years := response["years"]
	// When empty, years might be nil or an empty array
	if years != nil {
		if yearsArray, ok := years.([]interface{}); ok && len(yearsArray) != 0 {
			t.Errorf("Expected 0 years, got %d", len(yearsArray))
		}
	}
}

// TestGetBooks verifies the GetBooks handler
func TestGetBooks(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request for year 2025
	req := httptest.NewRequest(http.MethodGet, "/api/books?year=2025", nil)
	w := httptest.NewRecorder()

	// Execute
	GetBooks(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify content type
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	booksResp, ok := response["books"].([]interface{})
	if !ok {
		t.Fatal("Expected books array in response")
	}

	// Should have 2 books from 2025
	if len(booksResp) != 2 {
		t.Errorf("Expected 2 books for 2025, got %d", len(booksResp))
	}
}

// TestGetBooksMissingYear verifies error when year parameter is missing
func TestGetBooksMissingYear(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request without year parameter
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	w := httptest.NewRecorder()

	// Execute
	GetBooks(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetBooksInvalidYear verifies error when year parameter is invalid
func TestGetBooksInvalidYear(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request with invalid year
	req := httptest.NewRequest(http.MethodGet, "/api/books?year=invalid", nil)
	w := httptest.NewRecorder()

	// Execute
	GetBooks(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetBooksWithShelfFilter verifies filtering by shelf
func TestGetBooksWithShelfFilter(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request with shelf filter
	req := httptest.NewRequest(http.MethodGet, "/api/books?year=2025&shelf=read", nil)
	w := httptest.NewRecorder()

	// Execute
	GetBooks(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	booksResp, ok := response["books"].([]interface{})
	if !ok {
		t.Fatal("Expected books array in response")
	}

	// Should have 2 books with shelf "read" from 2025
	if len(booksResp) != 2 {
		t.Errorf("Expected 2 books with shelf 'read' for 2025, got %d", len(booksResp))
	}
}

// TestGetBooksWithMonthFilter verifies filtering by month
func TestGetBooksWithMonthFilter(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request with month filter (January = 1)
	req := httptest.NewRequest(http.MethodGet, "/api/books?year=2025&month=1", nil)
	w := httptest.NewRecorder()

	// Execute
	GetBooks(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	booksResp, ok := response["books"].([]interface{})
	if !ok {
		t.Fatal("Expected books array in response")
	}

	// Should have 1 book from January 2025
	if len(booksResp) != 1 {
		t.Errorf("Expected 1 book from January 2025, got %d", len(booksResp))
	}
}

// TestGetStats verifies the GetStats handler
func TestGetStats(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request for year 2025
	req := httptest.NewRequest(http.MethodGet, "/api/stats?year=2025", nil)
	w := httptest.NewRecorder()

	// Execute
	GetStats(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify content type
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check expected fields
	if year, ok := response["year"].(float64); !ok || int(year) != 2025 {
		t.Errorf("Expected year 2025, got %v", response["year"])
	}

	if totalBooks, ok := response["totalBooks"].(float64); !ok || int(totalBooks) != 2 {
		t.Errorf("Expected 2 total books, got %v", response["totalBooks"])
	}

	if totalPages, ok := response["totalPages"].(float64); !ok || int(totalPages) != 500 {
		t.Errorf("Expected 500 total pages, got %v", response["totalPages"])
	}

	if _, ok := response["averagePerMonth"]; !ok {
		t.Error("Expected averagePerMonth in response")
	}

	if _, ok := response["monthlyBreakdown"]; !ok {
		t.Error("Expected monthlyBreakdown in response")
	}
}

// TestGetStatsMissingYear verifies error when year parameter is missing
func TestGetStatsMissingYear(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request without year parameter
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	// Execute
	GetStats(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetStatsInvalidYear verifies error when year parameter is invalid
func TestGetStatsInvalidYear(t *testing.T) {
	// Setup
	SetBooks(setupTestBooks())

	// Create request with invalid year
	req := httptest.NewRequest(http.MethodGet, "/api/stats?year=notayear", nil)
	w := httptest.NewRecorder()

	// Execute
	GetStats(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetBooksUsesStoredCoverURL verifies that stored cover URLs are used
func TestGetBooksUsesStoredCoverURL(t *testing.T) {
	// Setup with a book that has a custom cover URL
	customCoverURL := "https://example.com/custom-cover.jpg"
	testBooks := []books.Book{
		{
			Title:    "Book With Custom Cover",
			Author:   "Test Author",
			DateRead: "2025/03/15",
			Pages:    250,
			Shelf:    "read",
			ISBN13:   "9781234567890",
			CoverURL: customCoverURL, // Custom cover URL should take priority
		},
		{
			Title:    "Book With ISBN Only",
			Author:   "Another Author",
			DateRead: "2025/04/20",
			Pages:    300,
			Shelf:    "read",
			ISBN13:   "9780987654321",
			CoverURL: "", // No custom cover, should generate from ISBN
		},
	}
	SetBooks(testBooks)

	// Create request for year 2025
	req := httptest.NewRequest(http.MethodGet, "/api/books?year=2025", nil)
	w := httptest.NewRecorder()

	// Execute
	GetBooks(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	booksResp, ok := response["books"].([]interface{})
	if !ok {
		t.Fatal("Expected books array in response")
	}

	// Should have 2 books
	if len(booksResp) != 2 {
		t.Fatalf("Expected 2 books, got %d", len(booksResp))
	}

	// Locate books by stable key (title) instead of relying on ordering
	var customBook map[string]interface{}
	var isbnBook map[string]interface{}

	for i, item := range booksResp {
		bookMap, ok := item.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected book at index %d to be an object, got %T", i, item)
		}

		titleVal, ok := bookMap["title"].(string)
		if !ok {
			t.Fatalf("Expected book at index %d to have string title, got %T", i, bookMap["title"])
		}

		switch titleVal {
		case "Book With Custom Cover":
			customBook = bookMap
		case "Book With ISBN Only":
			isbnBook = bookMap
		}
	}

	if customBook == nil {
		t.Fatalf("Did not find book with title %q in response", "Book With Custom Cover")
	}
	if isbnBook == nil {
		t.Fatalf("Did not find book with title %q in response", "Book With ISBN Only")
	}

	// Check custom-cover book - should use custom cover URL
	customCoverVal, ok := customBook["coverUrl"].(string)
	if !ok {
		t.Fatalf("Expected custom cover book to have string coverUrl, got %T", customBook["coverUrl"])
	}
	if customCoverVal != customCoverURL {
		t.Errorf("Expected custom cover URL %s, got %v", customCoverURL, customCoverVal)
	}

	// Check ISBN-only book - should generate from ISBN
	isbnCoverVal, ok := isbnBook["coverUrl"].(string)
	if !ok {
		t.Fatalf("Expected ISBN-only book to have string coverUrl, got %T", isbnBook["coverUrl"])
	}
	expectedURL := "https://covers.openlibrary.org/b/isbn/9780987654321-M.jpg"
	if isbnCoverVal != expectedURL {
		t.Errorf("Expected generated cover URL %s, got %v", expectedURL, isbnCoverVal)
	}
}
