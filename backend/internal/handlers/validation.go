package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	isbnRegex    = regexp.MustCompile(`^\d{10}$|^\d{13}$`)
	dateRegex    = regexp.MustCompile(`^\d{4}/\d{2}/\d{2}$`)
	validShelves = map[string]bool{
		"read":              true,
		"currently-reading": true,
		"to-read":           true,
	}
)

// validateYear checks that a year is within a reasonable range.
func validateYear(year int) error {
	if year < 1900 || year > 2100 {
		return fmt.Errorf("year must be between 1900 and 2100")
	}
	return nil
}

// validateBookRequest validates a book creation/update request and returns
// a list of field-level error messages. If the list is empty, the request is valid.
// When isUpdate is true, title and author are not required (partial updates).
func validateBookRequest(req BookRequest, isUpdate bool) []string {
	var errors []string

	if !isUpdate {
		if strings.TrimSpace(req.Title) == "" {
			errors = append(errors, "title is required")
		}
		if strings.TrimSpace(req.Author) == "" {
			errors = append(errors, "author is required")
		}
	}

	if len(req.Title) > 500 {
		errors = append(errors, "title must be 500 characters or less")
	}
	if len(req.Author) > 500 {
		errors = append(errors, "author must be 500 characters or less")
	}

	if req.Pages < 0 {
		errors = append(errors, "pages must be 0 or greater")
	}

	if req.ISBN != "" && !isbnRegex.MatchString(req.ISBN) {
		errors = append(errors, "isbn must be 10 or 13 digits")
	}
	if req.ISBN13 != "" && !isbnRegex.MatchString(req.ISBN13) {
		errors = append(errors, "isbn13 must be 10 or 13 digits")
	}

	if req.DateRead != "" && !dateRegex.MatchString(req.DateRead) {
		errors = append(errors, "dateRead must be in YYYY/MM/DD format")
	}

	if req.Shelf != "" && !validShelves[req.Shelf] {
		errors = append(errors, "shelf must be one of: read, currently-reading, to-read")
	}

	return errors
}

// writeValidationError writes a 400 response with a JSON body listing validation errors.
func writeValidationError(w http.ResponseWriter, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "Validation failed",
		"fields": errors,
	})
}
