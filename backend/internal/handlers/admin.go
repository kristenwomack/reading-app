package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/kristenwomack/reading-app/backend/internal/auth"
	"github.com/kristenwomack/reading-app/backend/internal/store"
)

var dataStore *store.Store

// SetStore sets the database store for handlers
func SetStore(s *store.Store) {
	dataStore = s
}

// AuthMiddleware protects routes that require authentication
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Login handles POST /api/auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := auth.CheckPassword(req.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	auth.SetAuthCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Logout handles POST /api/auth/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth.ClearAuthCookie(w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// CheckAuth handles GET /api/auth/check
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	authenticated := auth.IsAuthenticated(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"authenticated": authenticated})
}

// BookRequest represents a book creation/update request
type BookRequest struct {
	Title                   string `json:"title"`
	Author                  string `json:"author"`
	AdditionalAuthors       string `json:"additionalAuthors"`
	ISBN                    string `json:"isbn"`
	ISBN13                  string `json:"isbn13"`
	Publisher               string `json:"publisher"`
	Pages                   int    `json:"pages"`
	YearPublished           int    `json:"yearPublished"`
	OriginalPublicationYear int    `json:"originalPublicationYear"`
	DateRead                string `json:"dateRead"`
	DateAdded               string `json:"dateAdded"`
	Shelf                   string `json:"shelf"`
	Review                  string `json:"review"`
	CoverURL                string `json:"coverUrl"`
}

// CreateBook handles POST /api/books
func CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Author == "" {
		http.Error(w, "Title and author are required", http.StatusBadRequest)
		return
	}

	if req.Shelf == "" {
		req.Shelf = "read"
	}

	book := &store.Book{
		Title:                   req.Title,
		Author:                  req.Author,
		AdditionalAuthors:       req.AdditionalAuthors,
		ISBN:                    req.ISBN,
		ISBN13:                  req.ISBN13,
		Publisher:               req.Publisher,
		Pages:                   req.Pages,
		YearPublished:           req.YearPublished,
		OriginalPublicationYear: req.OriginalPublicationYear,
		DateRead:                req.DateRead,
		DateAdded:               req.DateAdded,
		Shelf:                   req.Shelf,
		Review:                  req.Review,
		CoverURL:                req.CoverURL,
	}

	id, err := dataStore.CreateBook(book)
	if err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// UpdateBook handles PUT /api/books/{id}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/books/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	book := &store.Book{
		ID:                      id,
		Title:                   req.Title,
		Author:                  req.Author,
		AdditionalAuthors:       req.AdditionalAuthors,
		ISBN:                    req.ISBN,
		ISBN13:                  req.ISBN13,
		Publisher:               req.Publisher,
		Pages:                   req.Pages,
		YearPublished:           req.YearPublished,
		OriginalPublicationYear: req.OriginalPublicationYear,
		DateRead:                req.DateRead,
		DateAdded:               req.DateAdded,
		Shelf:                   req.Shelf,
		Review:                  req.Review,
		CoverURL:                req.CoverURL,
	}

	if err := dataStore.UpdateBook(book); err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// DeleteBook handles DELETE /api/books/{id}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/books/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := dataStore.DeleteBook(id); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// exportBook matches the books.json Goodreads-style format for round-trip compatibility
type exportBook struct {
	Title                   string `json:"Title"`
	Author                  string `json:"Author"`
	AdditionalAuthors       string `json:"Additional Authors"`
	ISBN                    string `json:"ISBN"`
	ISBN13                  string `json:"ISBN13"`
	Publisher               string `json:"Publisher"`
	Pages                   int    `json:"Number of Pages"`
	YearPublished           int    `json:"Year Published"`
	OriginalPublicationYear int    `json:"Original Publication Year"`
	DateRead                string `json:"Date Read"`
	DateAdded               string `json:"Date Added"`
	Shelf                   string `json:"Shelf"`
	MyReview                string `json:"My Review"`
	CoverURL                string `json:"CoverURL,omitempty"`
}

// ExportBooks handles GET /api/export
func ExportBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	books, err := dataStore.GetAllBooks()
	if err != nil {
		http.Error(w, "Failed to get books", http.StatusInternalServerError)
		return
	}

	exported := make([]exportBook, len(books))
	for i, b := range books {
		exported[i] = exportBook{
			Title:                   b.Title,
			Author:                  b.Author,
			AdditionalAuthors:       b.AdditionalAuthors,
			ISBN:                    b.ISBN,
			ISBN13:                  b.ISBN13,
			Publisher:               b.Publisher,
			Pages:                   b.Pages,
			YearPublished:           b.YearPublished,
			OriginalPublicationYear: b.OriginalPublicationYear,
			DateRead:                b.DateRead,
			DateAdded:               b.DateAdded,
			Shelf:                   b.Shelf,
			MyReview:                b.Review,
			CoverURL:                b.CoverURL,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=books.json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(exported)
}

// GetGoal handles GET /api/goals/:year
func GetGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract year from path
	path := strings.TrimPrefix(r.URL.Path, "/api/goals/")
	year, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	goal, err := dataStore.GetGoal(year)
	if err != nil {
		http.Error(w, "Failed to get goal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if goal == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"year": year, "target": nil})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"year": goal.Year, "target": goal.BookTarget})
	}
}

// SetGoal handles POST /api/goals
func SetGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Year   int `json:"year"`
		Target int `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Year < 2000 || req.Year > 2100 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	if req.Target < 0 {
		http.Error(w, "Target must be positive", http.StatusBadRequest)
		return
	}

	if err := dataStore.SetGoal(req.Year, req.Target); err != nil {
		http.Error(w, "Failed to set goal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
