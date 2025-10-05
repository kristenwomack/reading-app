package books

import (
	"testing"
)

// T014: TestFilterByYear verifies filtering books by year
func TestFilterByYear(t *testing.T) {
	// Given a collection of books from multiple years
	books := []Book{
		{Title: "Book 2025-1", Author: "Author A", DateRead: "2025/09/19", Shelf: "read"},
		{Title: "Book 2025-2", Author: "Author B", DateRead: "2025/03/15", Shelf: "read"},
		{Title: "Book 2024-1", Author: "Author C", DateRead: "2024/12/01", Shelf: "read"},
		{Title: "Book 2024-2", Author: "Author D", DateRead: "2024/08/20", Shelf: "read"},
		{Title: "Book 2023-1", Author: "Author E", DateRead: "2023/06/10", Shelf: "read"},
	}
	
	// When filtering by year 2025
	filtered, err := FilterByYear(books, 2025)
	
	// Then no error should occur
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// And should return only 2025 books
	if len(filtered) != 2 {
		t.Errorf("Expected 2 books for 2025, got %d", len(filtered))
	}
	
	// And all filtered books should be from 2025
	for _, book := range filtered {
		date, _ := ParseDate(book.DateRead)
		if date.Year != 2025 {
			t.Errorf("Expected year 2025, got %d for book %q", date.Year, book.Title)
		}
	}
}

// T014: TestFilterByYearNoMatch verifies empty result when no books match
func TestFilterByYearNoMatch(t *testing.T) {
	// Given books from years other than 2026
	books := []Book{
		{Title: "Book 2025", Author: "Author A", DateRead: "2025/09/19", Shelf: "read"},
		{Title: "Book 2024", Author: "Author B", DateRead: "2024/03/15", Shelf: "read"},
	}
	
	// When filtering by year 2026 (no books)
	filtered, err := FilterByYear(books, 2026)
	
	// Then no error should occur
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// And result should be empty
	if len(filtered) != 0 {
		t.Errorf("Expected 0 books for 2026, got %d", len(filtered))
	}
}

// T014: TestFilterByShelf verifies filtering by shelf status
func TestFilterByShelf(t *testing.T) {
	// Given books with various shelf statuses
	books := []Book{
		{Title: "Read Book 1", Author: "Author A", DateRead: "2025/09/19", Shelf: "read"},
		{Title: "Read Book 2", Author: "Author B", DateRead: "2025/03/15", Shelf: "read"},
		{Title: "Reading Now", Author: "Author C", DateRead: "", Shelf: "currently-reading"},
		{Title: "To Read", Author: "Author D", DateRead: "", Shelf: "to-read"},
		{Title: "Read Book 3", Author: "Author E", DateRead: "2025/06/10", Shelf: "read"},
	}
	
	// When filtering by "read" shelf
	filtered := FilterByShelf(books, "read")
	
	// Then should return only "read" books
	if len(filtered) != 3 {
		t.Errorf("Expected 3 'read' books, got %d", len(filtered))
	}
	
	// And all filtered books should have "read" shelf
	for _, book := range filtered {
		if book.Shelf != "read" {
			t.Errorf("Expected shelf 'read', got %q for book %q", book.Shelf, book.Title)
		}
	}
}

// T014: TestFilterByShelfExcludesOthers verifies other shelf statuses are excluded
func TestFilterByShelfExcludesOthers(t *testing.T) {
	// Given books with various shelf statuses
	books := []Book{
		{Title: "Currently Reading", Author: "Author A", Shelf: "currently-reading"},
		{Title: "To Read", Author: "Author B", Shelf: "to-read"},
		{Title: "Did Not Finish", Author: "Author C", Shelf: "dnf"},
	}
	
	// When filtering by "read" shelf
	filtered := FilterByShelf(books, "read")
	
	// Then result should be empty
	if len(filtered) != 0 {
		t.Errorf("Expected 0 'read' books, got %d", len(filtered))
	}
}

// T014: TestFilterByYearAndShelf verifies combined filtering
func TestFilterByYearAndShelf(t *testing.T) {
	// Given books with various years and shelf statuses
	books := []Book{
		{Title: "2025 Read", Author: "Author A", DateRead: "2025/09/19", Shelf: "read"},
		{Title: "2025 Reading", Author: "Author B", DateRead: "2025/03/15", Shelf: "currently-reading"},
		{Title: "2024 Read", Author: "Author C", DateRead: "2024/12/01", Shelf: "read"},
		{Title: "2025 Read 2", Author: "Author D", DateRead: "2025/06/10", Shelf: "read"},
	}
	
	// When filtering by year 2025
	byYear, err := FilterByYear(books, 2025)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// And then filtering by "read" shelf
	filtered := FilterByShelf(byYear, "read")
	
	// Then should return only 2025 read books
	if len(filtered) != 2 {
		t.Errorf("Expected 2 books (2025 + read), got %d", len(filtered))
	}
	
	// And all should match both criteria
	for _, book := range filtered {
		date, _ := ParseDate(book.DateRead)
		if date.Year != 2025 {
			t.Errorf("Expected year 2025, got %d", date.Year)
		}
		if book.Shelf != "read" {
			t.Errorf("Expected shelf 'read', got %q", book.Shelf)
		}
	}
}
