package handlers

import (
	"strings"
	"testing"
)

func TestValidateYear(t *testing.T) {
	tests := []struct {
		name    string
		year    int
		wantErr bool
	}{
		{"valid 2025", 2025, false},
		{"valid 1900", 1900, false},
		{"valid 2100", 2100, false},
		{"too low", 1899, true},
		{"too high", 2101, true},
		{"zero", 0, true},
		{"negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateYear(tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateYear(%d) error = %v, wantErr %v", tt.year, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBookRequest_Create(t *testing.T) {
	tests := []struct {
		name       string
		req        BookRequest
		wantErrors []string
	}{
		{
			name: "valid minimal",
			req:  BookRequest{Title: "Test Book", Author: "Test Author"},
		},
		{
			name: "valid full",
			req: BookRequest{
				Title:    "Test Book",
				Author:   "Test Author",
				Pages:    200,
				ISBN:     "1234567890",
				DateRead: "2025/03/15",
				Shelf:    "read",
			},
		},
		{
			name:       "missing title",
			req:        BookRequest{Author: "Author"},
			wantErrors: []string{"title is required"},
		},
		{
			name:       "missing author",
			req:        BookRequest{Title: "Title"},
			wantErrors: []string{"author is required"},
		},
		{
			name:       "missing both",
			req:        BookRequest{},
			wantErrors: []string{"title is required", "author is required"},
		},
		{
			name:       "whitespace-only title",
			req:        BookRequest{Title: "   ", Author: "Author"},
			wantErrors: []string{"title is required"},
		},
		{
			name:       "title too long",
			req:        BookRequest{Title: strings.Repeat("a", 501), Author: "Author"},
			wantErrors: []string{"title must be 500 characters or less"},
		},
		{
			name:       "author too long",
			req:        BookRequest{Title: "Title", Author: strings.Repeat("a", 501)},
			wantErrors: []string{"author must be 500 characters or less"},
		},
		{
			name:       "negative pages",
			req:        BookRequest{Title: "Title", Author: "Author", Pages: -1},
			wantErrors: []string{"pages must be 0 or greater"},
		},
		{
			name: "zero pages valid",
			req:  BookRequest{Title: "Title", Author: "Author", Pages: 0},
		},
		{
			name:       "invalid ISBN too short",
			req:        BookRequest{Title: "Title", Author: "Author", ISBN: "12345"},
			wantErrors: []string{"isbn must be 10 or 13 digits"},
		},
		{
			name:       "invalid ISBN with letters",
			req:        BookRequest{Title: "Title", Author: "Author", ISBN: "123456789X"},
			wantErrors: []string{"isbn must be 10 or 13 digits"},
		},
		{
			name: "valid ISBN 10 digits",
			req:  BookRequest{Title: "Title", Author: "Author", ISBN: "1234567890"},
		},
		{
			name: "valid ISBN 13 digits",
			req:  BookRequest{Title: "Title", Author: "Author", ISBN: "9781234567890"},
		},
		{
			name:       "invalid dateRead format",
			req:        BookRequest{Title: "Title", Author: "Author", DateRead: "2025-03-15"},
			wantErrors: []string{"dateRead must be in YYYY/MM/DD format"},
		},
		{
			name: "valid dateRead",
			req:  BookRequest{Title: "Title", Author: "Author", DateRead: "2025/03/15"},
		},
		{
			name: "empty dateRead valid",
			req:  BookRequest{Title: "Title", Author: "Author", DateRead: ""},
		},
		{
			name:       "invalid shelf",
			req:        BookRequest{Title: "Title", Author: "Author", Shelf: "finished"},
			wantErrors: []string{"shelf must be one of: read, currently-reading, to-read"},
		},
		{
			name: "valid shelf read",
			req:  BookRequest{Title: "Title", Author: "Author", Shelf: "read"},
		},
		{
			name: "valid shelf currently-reading",
			req:  BookRequest{Title: "Title", Author: "Author", Shelf: "currently-reading"},
		},
		{
			name: "valid shelf to-read",
			req:  BookRequest{Title: "Title", Author: "Author", Shelf: "to-read"},
		},
		{
			name: "empty shelf valid",
			req:  BookRequest{Title: "Title", Author: "Author", Shelf: ""},
		},
		{
			name: "multiple errors",
			req: BookRequest{
				Title:    "",
				Author:   "",
				Pages:    -5,
				ISBN:     "bad",
				DateRead: "invalid",
				Shelf:    "unknown",
			},
			wantErrors: []string{
				"title is required",
				"author is required",
				"pages must be 0 or greater",
				"isbn must be 10 or 13 digits",
				"dateRead must be in YYYY/MM/DD format",
				"shelf must be one of: read, currently-reading, to-read",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateBookRequest(tt.req, false)
			if len(tt.wantErrors) == 0 {
				if len(errs) != 0 {
					t.Errorf("expected no errors, got %v", errs)
				}
				return
			}
			if len(errs) != len(tt.wantErrors) {
				t.Errorf("expected %d errors, got %d: %v", len(tt.wantErrors), len(errs), errs)
				return
			}
			for i, want := range tt.wantErrors {
				if errs[i] != want {
					t.Errorf("error[%d] = %q, want %q", i, errs[i], want)
				}
			}
		})
	}
}

func TestValidateBookRequest_Update(t *testing.T) {
	// For updates, title and author are not required
	errs := validateBookRequest(BookRequest{}, true)
	if len(errs) != 0 {
		t.Errorf("update with empty fields should not require title/author, got %v", errs)
	}

	// But other validations still apply
	errs = validateBookRequest(BookRequest{Pages: -1, ISBN: "bad"}, true)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors for invalid pages and ISBN, got %d: %v", len(errs), errs)
	}
}
