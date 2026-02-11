package store

import (
	"testing"
)

// setupTestStore creates a new in-memory SQLite database for testing
func setupTestStore(t *testing.T) *Store {
	t.Helper()

	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	return s
}

// TestNew verifies creating a new store
func TestNew(t *testing.T) {
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	if s.db == nil {
		t.Error("Expected database connection to be initialized")
	}
}

// TestCreateBook verifies creating a book
func TestCreateBook(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	book := &Book{
		Title:    "Test Book",
		Author:   "Test Author",
		Pages:    200,
		DateRead: "2025/01/15",
		Shelf:    "read",
	}

	id, err := s.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	if id <= 0 {
		t.Error("Expected positive book ID")
	}
}

// TestGetBook verifies retrieving a book by ID
func TestGetBook(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Create a book
	book := &Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  200,
		Shelf:  "read",
	}
	id, err := s.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Get the book
	retrieved, err := s.GetBook(id)
	if err != nil {
		t.Fatalf("Failed to get book: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected book to be retrieved")
	}

	if retrieved.Title != book.Title {
		t.Errorf("Expected title %q, got %q", book.Title, retrieved.Title)
	}

	if retrieved.Author != book.Author {
		t.Errorf("Expected author %q, got %q", book.Author, retrieved.Author)
	}
}

// TestGetBookNotFound verifies getting a non-existent book
func TestGetBookNotFound(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Try to get a book that doesn't exist
	retrieved, err := s.GetBook(999999)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if retrieved != nil {
		t.Error("Expected nil for non-existent book")
	}
}

// TestGetAllBooks verifies retrieving all books
func TestGetAllBooks(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Create multiple books
	books := []Book{
		{Title: "Book 1", Author: "Author 1", Shelf: "read"},
		{Title: "Book 2", Author: "Author 2", Shelf: "currently-reading"},
		{Title: "Book 3", Author: "Author 3", Shelf: "read"},
	}

	for i := range books {
		_, err := s.CreateBook(&books[i])
		if err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
	}

	// Get all books
	allBooks, err := s.GetAllBooks()
	if err != nil {
		t.Fatalf("Failed to get all books: %v", err)
	}

	if len(allBooks) != 3 {
		t.Errorf("Expected 3 books, got %d", len(allBooks))
	}
}

// TestUpdateBook verifies updating a book
func TestUpdateBook(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Create a book
	book := &Book{
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  200,
		Shelf:  "read",
	}
	id, err := s.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Update the book
	book.ID = id
	book.Title = "Updated Title"
	book.Author = "Updated Author"
	book.Pages = 300

	err = s.UpdateBook(book)
	if err != nil {
		t.Fatalf("Failed to update book: %v", err)
	}

	// Verify the update
	updated, err := s.GetBook(id)
	if err != nil {
		t.Fatalf("Failed to get updated book: %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("Expected title %q, got %q", "Updated Title", updated.Title)
	}

	if updated.Author != "Updated Author" {
		t.Errorf("Expected author %q, got %q", "Updated Author", updated.Author)
	}

	if updated.Pages != 300 {
		t.Errorf("Expected pages 300, got %d", updated.Pages)
	}
}

// TestDeleteBook verifies deleting a book
func TestDeleteBook(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Create a book
	book := &Book{
		Title:  "Book to Delete",
		Author: "Author",
		Shelf:  "read",
	}
	id, err := s.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Delete the book
	err = s.DeleteBook(id)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}

	// Verify the book is gone
	deleted, err := s.GetBook(id)
	if err != nil {
		t.Fatalf("Failed to get deleted book: %v", err)
	}

	if deleted != nil {
		t.Error("Expected book to be deleted")
	}
}

// TestBookCount verifies counting books
func TestBookCount(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Initially should have 0 books
	count, err := s.BookCount()
	if err != nil {
		t.Fatalf("Failed to get book count: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 books, got %d", count)
	}

	// Create some books
	for i := 0; i < 5; i++ {
		book := &Book{
			Title:  "Book",
			Author: "Author",
			Shelf:  "read",
		}
		_, err := s.CreateBook(book)
		if err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
	}

	// Should now have 5 books
	count, err = s.BookCount()
	if err != nil {
		t.Fatalf("Failed to get book count: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected 5 books, got %d", count)
	}
}

// TestGetGoal verifies getting a reading goal
func TestGetGoal(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Set a goal
	err := s.SetGoal(2025, 50)
	if err != nil {
		t.Fatalf("Failed to set goal: %v", err)
	}

	// Get the goal
	goal, err := s.GetGoal(2025)
	if err != nil {
		t.Fatalf("Failed to get goal: %v", err)
	}

	if goal == nil {
		t.Fatal("Expected goal to be retrieved")
	}

	if goal.Year != 2025 {
		t.Errorf("Expected year 2025, got %d", goal.Year)
	}

	if goal.BookTarget != 50 {
		t.Errorf("Expected target 50, got %d", goal.BookTarget)
	}
}

// TestGetGoalNotFound verifies getting a non-existent goal
func TestGetGoalNotFound(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Try to get a goal that doesn't exist
	goal, err := s.GetGoal(2099)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if goal != nil {
		t.Error("Expected nil for non-existent goal")
	}
}

// TestSetGoal verifies setting a reading goal
func TestSetGoal(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Set a goal
	err := s.SetGoal(2025, 50)
	if err != nil {
		t.Fatalf("Failed to set goal: %v", err)
	}

	// Verify the goal
	goal, err := s.GetGoal(2025)
	if err != nil {
		t.Fatalf("Failed to get goal: %v", err)
	}

	if goal == nil || goal.BookTarget != 50 {
		t.Error("Expected goal to be set")
	}
}

// TestSetGoalUpdate verifies updating an existing goal
func TestSetGoalUpdate(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Set initial goal
	err := s.SetGoal(2025, 50)
	if err != nil {
		t.Fatalf("Failed to set goal: %v", err)
	}

	// Update the goal
	err = s.SetGoal(2025, 75)
	if err != nil {
		t.Fatalf("Failed to update goal: %v", err)
	}

	// Verify the update
	goal, err := s.GetGoal(2025)
	if err != nil {
		t.Fatalf("Failed to get goal: %v", err)
	}

	if goal == nil {
		t.Fatal("Expected goal to exist")
	}

	if goal.BookTarget != 75 {
		t.Errorf("Expected target 75, got %d", goal.BookTarget)
	}
}

// TestGetSetting verifies getting a setting
func TestGetSetting(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Set a setting
	err := s.SetSetting("theme", "dark")
	if err != nil {
		t.Fatalf("Failed to set setting: %v", err)
	}

	// Get the setting
	value, err := s.GetSetting("theme")
	if err != nil {
		t.Fatalf("Failed to get setting: %v", err)
	}

	if value != "dark" {
		t.Errorf("Expected value %q, got %q", "dark", value)
	}
}

// TestGetSettingNotFound verifies getting a non-existent setting
func TestGetSettingNotFound(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Try to get a setting that doesn't exist
	value, err := s.GetSetting("nonexistent")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if value != "" {
		t.Errorf("Expected empty string, got %q", value)
	}
}

// TestSetSetting verifies setting a value
func TestSetSetting(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Set a setting
	err := s.SetSetting("language", "en")
	if err != nil {
		t.Fatalf("Failed to set setting: %v", err)
	}

	// Verify the setting
	value, err := s.GetSetting("language")
	if err != nil {
		t.Fatalf("Failed to get setting: %v", err)
	}

	if value != "en" {
		t.Errorf("Expected value %q, got %q", "en", value)
	}
}

// TestSetSettingUpdate verifies updating an existing setting
func TestSetSettingUpdate(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Set initial value
	err := s.SetSetting("theme", "light")
	if err != nil {
		t.Fatalf("Failed to set setting: %v", err)
	}

	// Update the value
	err = s.SetSetting("theme", "dark")
	if err != nil {
		t.Fatalf("Failed to update setting: %v", err)
	}

	// Verify the update
	value, err := s.GetSetting("theme")
	if err != nil {
		t.Fatalf("Failed to get setting: %v", err)
	}

	if value != "dark" {
		t.Errorf("Expected value %q, got %q", "dark", value)
	}
}
