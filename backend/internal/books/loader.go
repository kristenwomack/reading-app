package books

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Book represents a book entity from books.json
type Book struct {
	Title                    interface{} `json:"Title"`
	Author                   string      `json:"Author"`
	AdditionalAuthors        interface{} `json:"Additional Authors"`
	ISBN                     interface{} `json:"ISBN"`
	ISBN13                   interface{} `json:"ISBN13"`
	Publisher                interface{} `json:"Publisher"`
	Pages                    interface{} `json:"Number of Pages"`
	YearPublished            interface{} `json:"Year Published"`
	OriginalPublicationYear  interface{} `json:"Original Publication Year"`
	DateRead                 string      `json:"Date Read"`
	DateAdded                interface{} `json:"Date Added"`
	Bookshelves              string      `json:"Bookshelves"`
	BookshelvesWithPositions interface{} `json:"Bookshelves with positions"`
	Shelf                    string      `json:"Shelf"`
	MyReview                 interface{} `json:"My Review"`
	CoverURL                 string      `json:"CoverURL,omitempty"`
}

// GetTitle returns the title as a string
func (b *Book) GetTitle() string {
	switch v := b.Title.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return ""
	}
}

// GetPages returns the page count as an integer
func (b *Book) GetPages() int {
	switch v := b.Pages.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		// Try to parse string to int
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}

// ParsedDate represents a parsed date
type ParsedDate struct {
	Year  int
	Month int
	Day   int
}

// LoadBooks loads books from a JSON file
func LoadBooks(filepath string) ([]Book, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return books, nil
}

// ParseDate parses a date string in format YYYY/MM/DD or YYYY
func ParseDate(dateStr string) (ParsedDate, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ParsedDate{}, fmt.Errorf("empty date string")
	}

	parts := strings.Split(dateStr, "/")
	
	// Parse year (required)
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return ParsedDate{}, fmt.Errorf("invalid year: %w", err)
	}
	if year < 1900 {
		return ParsedDate{}, fmt.Errorf("year must be >= 1900")
	}

	date := ParsedDate{Year: year}

	// Parse month (optional)
	if len(parts) > 1 && parts[1] != "" {
		month, err := strconv.Atoi(parts[1])
		if err != nil {
			return ParsedDate{}, fmt.Errorf("invalid month: %w", err)
		}
		if month < 1 || month > 12 {
			return ParsedDate{}, fmt.Errorf("month must be 1-12")
		}
		date.Month = month
	}

	// Parse day (optional)
	if len(parts) > 2 && parts[2] != "" {
		day, err := strconv.Atoi(parts[2])
		if err != nil {
			return ParsedDate{}, fmt.Errorf("invalid day: %w", err)
		}
		if day < 1 || day > 31 {
			return ParsedDate{}, fmt.Errorf("day must be 1-31")
		}
		date.Day = day
	}

	return date, nil
}
