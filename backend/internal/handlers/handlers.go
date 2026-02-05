package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kristenwomack/reading-app/backend/internal/books"
	"github.com/kristenwomack/reading-app/backend/internal/store"
)

var cachedBooks []books.Book

// SetBooks sets the cached books for handlers to use (legacy support)
func SetBooks(b []books.Book) {
	cachedBooks = b
}

// getBooks returns books from store if available, otherwise from cache
func getBooks() []books.Book {
	if dataStore != nil {
		storeBooks, err := dataStore.GetAllBooks()
		if err == nil && len(storeBooks) > 0 {
			return convertStoreBooks(storeBooks)
		}
	}
	return cachedBooks
}

// convertStoreBooks converts store.Book slice to books.Book slice
func convertStoreBooks(storeBooks []store.Book) []books.Book {
	result := make([]books.Book, len(storeBooks))
	for i, sb := range storeBooks {
		result[i] = books.Book{
			Title:                   sb.Title,
			Author:                  sb.Author,
			AdditionalAuthors:       sb.AdditionalAuthors,
			ISBN:                    sb.ISBN,
			ISBN13:                  sb.ISBN13,
			Publisher:               sb.Publisher,
			Pages:                   sb.Pages,
			YearPublished:           sb.YearPublished,
			OriginalPublicationYear: sb.OriginalPublicationYear,
			DateRead:                sb.DateRead,
			DateAdded:               sb.DateAdded,
			Bookshelves:             sb.Shelf,
			Shelf:                   sb.Shelf,
			MyReview:                sb.Review,
		}
	}
	return result
}

// GetYears returns available years with book counts
func GetYears(w http.ResponseWriter, r *http.Request) {
	yearCounts := make(map[int]int)
	
	for _, book := range getBooks() {
		if book.DateRead == "" || book.Shelf != "read" {
			continue
		}
		
		date, err := books.ParseDate(book.DateRead)
		if err != nil {
			continue
		}
		
		yearCounts[date.Year]++
	}
	
	type YearInfo struct {
		Year  int `json:"year"`
		Count int `json:"count"`
	}
	
	var years []YearInfo
	for year, count := range yearCounts {
		years = append(years, YearInfo{Year: year, Count: count})
	}
	
	response := map[string]interface{}{
		"years": years,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBooks returns books for a specific year
func GetBooks(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "year parameter required", http.StatusBadRequest)
		return
	}
	
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year parameter", http.StatusBadRequest)
		return
	}

	// Optional filters
	shelfFilter := r.URL.Query().Get("shelf")
	monthStr := r.URL.Query().Get("month")
	var monthFilter int
	if monthStr != "" {
		monthFilter, _ = strconv.Atoi(monthStr)
	}
	
	allBooks := getBooks()
	filtered, _ := books.FilterByYear(allBooks, year)
	
	// Filter by shelf if specified, otherwise return all shelves
	var booksToReturn []books.Book
	if shelfFilter != "" {
		booksToReturn = books.FilterByShelf(filtered, shelfFilter)
	} else {
		booksToReturn = filtered
	}
	
	type BookResponse struct {
		ID       int64  `json:"id,omitempty"`
		Title    string `json:"title"`
		Author   string `json:"author"`
		DateRead string `json:"dateRead"`
		Pages    int    `json:"pages"`
		Month    int    `json:"month"`
		Shelf    string `json:"shelf"`
		ISBN     string `json:"isbn,omitempty"`
		CoverURL string `json:"coverUrl,omitempty"`
	}
	
	var responseBooks []BookResponse
	for _, book := range booksToReturn {
		date, _ := books.ParseDate(book.DateRead)
		
		// Filter by month if specified
		if monthFilter > 0 && date.Month != monthFilter {
			continue
		}
		
		// Build cover URL from ISBN if available
		isbn := getISBN(book)
		coverURL := ""
		if isbn != "" {
			coverURL = "https://covers.openlibrary.org/b/isbn/" + isbn + "-M.jpg"
		}
		
		responseBooks = append(responseBooks, BookResponse{
			Title:    book.GetTitle(),
			Author:   book.Author,
			DateRead: book.DateRead,
			Pages:    book.GetPages(),
			Month:    date.Month,
			Shelf:    book.Shelf,
			ISBN:     isbn,
			CoverURL: coverURL,
		})
	}
	
	response := map[string]interface{}{
		"books": responseBooks,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getISBN extracts ISBN from book
func getISBN(book books.Book) string {
	if s, ok := book.ISBN13.(string); ok && s != "" {
		return s
	}
	if f, ok := book.ISBN13.(float64); ok && f > 0 {
		return strconv.FormatInt(int64(f), 10)
	}
	if s, ok := book.ISBN.(string); ok && s != "" {
		return s
	}
	return ""
}

// GetStats returns statistics for a specific year
func GetStats(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "year parameter required", http.StatusBadRequest)
		return
	}
	
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year parameter", http.StatusBadRequest)
		return
	}
	
	filtered, _ := books.FilterByYear(getBooks(), year)
	readBooks := books.FilterByShelf(filtered, "read")
	
	stats := books.CalculateStatistics(readBooks, year)
	breakdown := books.CalculateMonthlyBreakdown(readBooks)
	
	response := map[string]interface{}{
		"year":             stats.Year,
		"totalBooks":       stats.TotalBooks,
		"totalPages":       stats.TotalPages,
		"averagePerMonth":  stats.AveragePerMonth,
		"monthlyBreakdown": breakdown,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
