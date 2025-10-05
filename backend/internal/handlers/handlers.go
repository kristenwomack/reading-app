package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kristenwomack/reading-app/backend/internal/books"
)

var cachedBooks []books.Book

// SetBooks sets the cached books for handlers to use
func SetBooks(b []books.Book) {
	cachedBooks = b
}

// GetYears returns available years with book counts
func GetYears(w http.ResponseWriter, r *http.Request) {
	yearCounts := make(map[int]int)
	
	for _, book := range cachedBooks {
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
	
	filtered, _ := books.FilterByYear(cachedBooks, year)
	readBooks := books.FilterByShelf(filtered, "read")
	
	type BookResponse struct {
		Title    string `json:"title"`
		Author   string `json:"author"`
		DateRead string `json:"dateRead"`
		Pages    int    `json:"pages"`
		Month    int    `json:"month"`
	}
	
	var responseBooks []BookResponse
	for _, book := range readBooks {
		date, _ := books.ParseDate(book.DateRead)
		responseBooks = append(responseBooks, BookResponse{
			Title:    book.GetTitle(),
			Author:   book.Author,
			DateRead: book.DateRead,
			Pages:    book.GetPages(),
			Month:    date.Month,
		})
	}
	
	response := map[string]interface{}{
		"books": responseBooks,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	
	filtered, _ := books.FilterByYear(cachedBooks, year)
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
