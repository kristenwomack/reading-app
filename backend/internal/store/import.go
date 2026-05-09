package store

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kristenwomack/reading-app/backend/internal/books"
	"github.com/kristenwomack/reading-app/backend/internal/covers"
)

// ImportFromJSON imports books from the old JSON format into the database
func (s *Store) ImportFromJSON(jsonBooks []books.Book) (int, error) {
	// Check if already imported
	count, err := s.BookCount()
	if err != nil {
		return 0, fmt.Errorf("failed to check book count: %w", err)
	}
	if count > 0 {
		return 0, nil // Already has books, skip import
	}

	imported := 0
	for _, jb := range jsonBooks {
		b := Book{
			Title:                   jb.GetTitle(),
			Author:                  jb.Author,
			AdditionalAuthors:       toString(jb.AdditionalAuthors),
			ISBN:                    toString(jb.ISBN),
			ISBN13:                  toString(jb.ISBN13),
			Publisher:               toString(jb.Publisher),
			Pages:                   jb.GetPages(),
			YearPublished:           toInt(jb.YearPublished),
			OriginalPublicationYear: toInt(jb.OriginalPublicationYear),
			DateRead:                jb.DateRead,
			DateAdded:               toString(jb.DateAdded),
			Shelf:                   jb.Shelf,
			Review:                  toString(jb.MyReview),
			CoverURL:                covers.CoverURLByISBN(covers.PreferISBN(toString(jb.ISBN), toString(jb.ISBN13))),
		}

		if b.Title == "" || b.Author == "" {
			continue // Skip invalid entries
		}

		if _, err := s.CreateBook(&b); err != nil {
			return imported, fmt.Errorf("failed to import book %q: %w", b.Title, err)
		}
		imported++
	}

	// Mark import as complete
	s.SetSetting("imported_at", time.Now().Format(time.RFC3339))

	return imported, nil
}

// toString converts various types to string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// toInt converts various types to int
func toInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}

// buildCoverURL creates an Open Library cover URL from ISBN (kept for backward compatibility)
func buildCoverURL(isbn, isbn13 interface{}) string {
	return covers.CoverURLByISBN(covers.PreferISBN(toString(isbn), toString(isbn13)))
}
