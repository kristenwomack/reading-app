package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kristenwomack/reading-app/backend/internal/auth"
	"github.com/kristenwomack/reading-app/backend/internal/store"
)

// setupTestStore creates a test SQLite database in memory
func setupTestStore(t *testing.T) *store.Store {
	t.Helper()

	// Create in-memory database for testing
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	SetStore(s)
	return s
}

// teardownTestStore closes the test database
func teardownTestStore(t *testing.T, s *store.Store) {
	t.Helper()
	if err := s.Close(); err != nil {
		t.Errorf("Failed to close test store: %v", err)
	}
	SetStore(nil)
}

// TestLogin verifies successful login
func TestLogin(t *testing.T) {
	// Setup password
	os.Setenv("READING_APP_PASSWORD", "testpassword")
	defer os.Unsetenv("READING_APP_PASSWORD")

	// Create request with valid password
	body := bytes.NewBufferString(`{"password":"testpassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	w := httptest.NewRecorder()

	// Execute
	Login(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["success"] {
		t.Error("Expected success to be true")
	}

	// Verify cookie was set
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			found = true
			if cookie.Value == "" {
				t.Error("Expected auth_token cookie to have a value")
			}
		}
	}
	if !found {
		t.Error("Expected auth_token cookie to be set")
	}
}

// TestLoginInvalidPassword verifies login with wrong password
func TestLoginInvalidPassword(t *testing.T) {
	// Setup password
	os.Setenv("READING_APP_PASSWORD", "correctpassword")
	defer os.Unsetenv("READING_APP_PASSWORD")

	// Create request with invalid password
	body := bytes.NewBufferString(`{"password":"wrongpassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	w := httptest.NewRecorder()

	// Execute
	Login(w, req)

	// Verify response code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// TestLoginMissingPassword verifies login when no password is set
func TestLoginMissingPassword(t *testing.T) {
	// Ensure no password is set
	os.Unsetenv("READING_APP_PASSWORD")

	// Create request
	body := bytes.NewBufferString(`{"password":"anypassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	w := httptest.NewRecorder()

	// Execute
	Login(w, req)

	// Verify response code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// TestLoginInvalidJSON verifies login with malformed JSON
func TestLoginInvalidJSON(t *testing.T) {
	// Create request with invalid JSON
	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	w := httptest.NewRecorder()

	// Execute
	Login(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestLoginWrongMethod verifies login with wrong HTTP method
func TestLoginWrongMethod(t *testing.T) {
	// Create GET request instead of POST
	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	w := httptest.NewRecorder()

	// Execute
	Login(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestLogout verifies successful logout
func TestLogout(t *testing.T) {
	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	// Execute
	Logout(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["success"] {
		t.Error("Expected success to be true")
	}

	// Verify cookie was cleared
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			found = true
			if cookie.MaxAge != -1 {
				t.Error("Expected auth_token cookie to be cleared (MaxAge=-1)")
			}
		}
	}
	if !found {
		t.Error("Expected auth_token cookie to be set for clearing")
	}
}

// TestLogoutWrongMethod verifies logout with wrong HTTP method
func TestLogoutWrongMethod(t *testing.T) {
	// Create GET request instead of POST
	req := httptest.NewRequest(http.MethodGet, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	// Execute
	Logout(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestCheckAuthAuthenticated verifies CheckAuth with valid token
func TestCheckAuthAuthenticated(t *testing.T) {
	// Generate a valid token
	token, err := auth.GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create request with auth cookie
	req := httptest.NewRequest(http.MethodGet, "/api/auth/check", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
	})
	w := httptest.NewRecorder()

	// Execute
	CheckAuth(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["authenticated"] {
		t.Error("Expected authenticated to be true")
	}
}

// TestCheckAuthUnauthenticated verifies CheckAuth without token
func TestCheckAuthUnauthenticated(t *testing.T) {
	// Create request without auth cookie
	req := httptest.NewRequest(http.MethodGet, "/api/auth/check", nil)
	w := httptest.NewRecorder()

	// Execute
	CheckAuth(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["authenticated"] {
		t.Error("Expected authenticated to be false")
	}
}

// TestCreateBook verifies creating a new book
func TestCreateBook(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request body
	bookData := map[string]interface{}{
		"title":    "New Book",
		"author":   "New Author",
		"pages":    250,
		"dateRead": "2025/03/15",
		"shelf":    "read",
	}
	body, _ := json.Marshal(bookData)
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute
	CreateBook(w, req)

	// Verify response code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify response body
	var response map[string]int64
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["id"] <= 0 {
		t.Error("Expected positive book ID")
	}
}

// TestCreateBookMissingTitle verifies error when title is missing
func TestCreateBookMissingTitle(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request body without title
	bookData := map[string]interface{}{
		"author": "New Author",
		"pages":  250,
	}
	body, _ := json.Marshal(bookData)
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute
	CreateBook(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateBookInvalidJSON verifies error with malformed JSON
func TestCreateBookInvalidJSON(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request with invalid JSON
	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/api/books", body)
	w := httptest.NewRecorder()

	// Execute
	CreateBook(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateBookWrongMethod verifies error with wrong HTTP method
func TestCreateBookWrongMethod(t *testing.T) {
	// Create GET request instead of POST
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	w := httptest.NewRecorder()

	// Execute
	CreateBook(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestUpdateBook verifies updating an existing book
func TestUpdateBook(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// First, create a book
	book := &store.Book{
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
	bookData := map[string]interface{}{
		"title":  "Updated Title",
		"author": "Updated Author",
		"pages":  300,
		"shelf":  "read",
	}
	body, _ := json.Marshal(bookData)
	req := httptest.NewRequest(http.MethodPut, "/api/books/"+string(rune(id+'0')), bytes.NewBuffer(body))
	req.URL.Path = "/api/books/" + string(rune(id+'0'))
	w := httptest.NewRecorder()

	// Execute
	UpdateBook(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["success"] {
		t.Error("Expected success to be true")
	}
}

// TestUpdateBookInvalidID verifies error with invalid book ID
func TestUpdateBookInvalidID(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request with invalid ID
	bookData := map[string]interface{}{
		"title":  "Updated Title",
		"author": "Updated Author",
	}
	body, _ := json.Marshal(bookData)
	req := httptest.NewRequest(http.MethodPut, "/api/books/invalid", bytes.NewBuffer(body))
	req.URL.Path = "/api/books/invalid"
	w := httptest.NewRecorder()

	// Execute
	UpdateBook(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateBookWrongMethod verifies error with wrong HTTP method
func TestUpdateBookWrongMethod(t *testing.T) {
	// Create GET request instead of PUT
	req := httptest.NewRequest(http.MethodGet, "/api/books/1", nil)
	w := httptest.NewRecorder()

	// Execute
	UpdateBook(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestDeleteBook verifies deleting a book
func TestDeleteBook(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// First, create a book
	book := &store.Book{
		Title:  "Book to Delete",
		Author: "Author",
		Pages:  200,
		Shelf:  "read",
	}
	id, err := s.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Delete the book
	req := httptest.NewRequest(http.MethodDelete, "/api/books/"+string(rune(id+'0')), nil)
	req.URL.Path = "/api/books/" + string(rune(id+'0'))
	w := httptest.NewRecorder()

	// Execute
	DeleteBook(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["success"] {
		t.Error("Expected success to be true")
	}
}

// TestDeleteBookInvalidID verifies error with invalid book ID
func TestDeleteBookInvalidID(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request with invalid ID
	req := httptest.NewRequest(http.MethodDelete, "/api/books/notanumber", nil)
	req.URL.Path = "/api/books/notanumber"
	w := httptest.NewRecorder()

	// Execute
	DeleteBook(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestDeleteBookWrongMethod verifies error with wrong HTTP method
func TestDeleteBookWrongMethod(t *testing.T) {
	// Create GET request instead of DELETE
	req := httptest.NewRequest(http.MethodGet, "/api/books/1", nil)
	w := httptest.NewRecorder()

	// Execute
	DeleteBook(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestGetGoal verifies getting a reading goal
func TestGetGoal(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Set a goal
	err := s.SetGoal(2025, 50)
	if err != nil {
		t.Fatalf("Failed to set goal: %v", err)
	}

	// Get the goal
	req := httptest.NewRequest(http.MethodGet, "/api/goals/2025", nil)
	req.URL.Path = "/api/goals/2025"
	w := httptest.NewRecorder()

	// Execute
	GetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if year, ok := response["year"].(float64); !ok || int(year) != 2025 {
		t.Errorf("Expected year 2025, got %v", response["year"])
	}

	if target, ok := response["target"].(float64); !ok || int(target) != 50 {
		t.Errorf("Expected target 50, got %v", response["target"])
	}
}

// TestGetGoalNotFound verifies getting a goal that doesn't exist
func TestGetGoalNotFound(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Get a goal that doesn't exist
	req := httptest.NewRequest(http.MethodGet, "/api/goals/2099", nil)
	req.URL.Path = "/api/goals/2099"
	w := httptest.NewRecorder()

	// Execute
	GetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body contains nil target
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["target"] != nil {
		t.Errorf("Expected target to be nil, got %v", response["target"])
	}
}

// TestGetGoalInvalidYear verifies error with invalid year
func TestGetGoalInvalidYear(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request with invalid year
	req := httptest.NewRequest(http.MethodGet, "/api/goals/notayear", nil)
	req.URL.Path = "/api/goals/notayear"
	w := httptest.NewRecorder()

	// Execute
	GetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetGoalWrongMethod verifies error with wrong HTTP method
func TestGetGoalWrongMethod(t *testing.T) {
	// Create POST request instead of GET
	req := httptest.NewRequest(http.MethodPost, "/api/goals/2025", nil)
	w := httptest.NewRecorder()

	// Execute
	GetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestSetGoal verifies setting a reading goal
func TestSetGoal(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request body
	goalData := map[string]int{
		"year":   2025,
		"target": 50,
	}
	body, _ := json.Marshal(goalData)
	req := httptest.NewRequest(http.MethodPost, "/api/goals", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute
	SetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var response map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["success"] {
		t.Error("Expected success to be true")
	}

	// Verify the goal was set in the database
	goal, err := s.GetGoal(2025)
	if err != nil {
		t.Fatalf("Failed to get goal: %v", err)
	}
	if goal == nil || goal.BookTarget != 50 {
		t.Errorf("Expected goal target to be 50, got %v", goal)
	}
}

// TestSetGoalInvalidYear verifies error with invalid year
func TestSetGoalInvalidYear(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request body with invalid year
	goalData := map[string]int{
		"year":   1999,
		"target": 50,
	}
	body, _ := json.Marshal(goalData)
	req := httptest.NewRequest(http.MethodPost, "/api/goals", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute
	SetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestSetGoalNegativeTarget verifies error with negative target
func TestSetGoalNegativeTarget(t *testing.T) {
	// Setup test database
	s := setupTestStore(t)
	defer teardownTestStore(t, s)

	// Create request body with negative target
	goalData := map[string]int{
		"year":   2025,
		"target": -10,
	}
	body, _ := json.Marshal(goalData)
	req := httptest.NewRequest(http.MethodPost, "/api/goals", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute
	SetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestSetGoalWrongMethod verifies error with wrong HTTP method
func TestSetGoalWrongMethod(t *testing.T) {
	// Create GET request instead of POST
	req := httptest.NewRequest(http.MethodGet, "/api/goals", nil)
	w := httptest.NewRecorder()

	// Execute
	SetGoal(w, req)

	// Verify response code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestAuthMiddleware verifies the auth middleware
func TestAuthMiddleware(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with auth middleware
	protectedHandler := AuthMiddleware(testHandler)

	t.Run("authenticated", func(t *testing.T) {
		// Generate a valid token
		token, err := auth.GenerateToken()
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Create request with auth cookie
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: token,
		})
		w := httptest.NewRecorder()

		// Execute
		protectedHandler(w, req)

		// Verify response code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("unauthenticated", func(t *testing.T) {
		// Create request without auth cookie
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		w := httptest.NewRecorder()

		// Execute
		protectedHandler(w, req)

		// Verify response code
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}
