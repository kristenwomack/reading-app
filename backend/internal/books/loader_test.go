package books

import (
	"testing"
)

// T012: TestLoadBooks verifies that books can be loaded from a JSON file
func TestLoadBooks(t *testing.T) {
	// Given a valid test data file
	filepath := "../testdata/books_test.json"
	
	// When loading books
	books, err := LoadBooks(filepath)
	
	// Then no error should occur
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// And books should be loaded
	if len(books) == 0 {
		t.Error("Expected books to be loaded, got empty slice")
	}
	
	// And each book should have required fields
	if len(books) > 0 {
		book := books[0]
		if book.Title == "" {
			t.Error("Expected book to have Title")
		}
		if book.Author == "" {
			t.Error("Expected book to have Author")
		}
	}
}

// T012: TestLoadBooksFileNotFound verifies error handling when file doesn't exist
func TestLoadBooksFileNotFound(t *testing.T) {
	// Given a non-existent file path
	filepath := "../testdata/nonexistent.json"
	
	// When trying to load books
	books, err := LoadBooks(filepath)
	
	// Then an error should occur
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	// And books should be empty
	if len(books) != 0 {
		t.Errorf("Expected empty books slice, got %d books", len(books))
	}
}

// T012: TestLoadBooksInvalidJSON verifies error handling for malformed JSON
func TestLoadBooksInvalidJSON(t *testing.T) {
	// Given a file with invalid JSON
	filepath := "../testdata/invalid.json"
	
	// When trying to load books
	books, err := LoadBooks(filepath)
	
	// Then an error should occur
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
	
	// And books should be empty
	if len(books) != 0 {
		t.Errorf("Expected empty books slice, got %d books", len(books))
	}
}

// T013: TestParseDate verifies parsing of full date format YYYY/MM/DD
func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantYear int
		wantMonth int
		wantDay  int
	}{
		{
			name:     "full date",
			input:    "2025/09/19",
			wantYear: 2025,
			wantMonth: 9,
			wantDay:  19,
		},
		{
			name:     "different full date",
			input:    "2024/12/31",
			wantYear: 2024,
			wantMonth: 12,
			wantDay:  31,
		},
		{
			name:     "january first",
			input:    "2023/01/01",
			wantYear: 2023,
			wantMonth: 1,
			wantDay:  1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When parsing the date
			date, err := ParseDate(tt.input)
			
			// Then no error should occur
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			
			// And year should match
			if date.Year != tt.wantYear {
				t.Errorf("Year: got %d, want %d", date.Year, tt.wantYear)
			}
			
			// And month should match
			if date.Month != tt.wantMonth {
				t.Errorf("Month: got %d, want %d", date.Month, tt.wantMonth)
			}
			
			// And day should match
			if date.Day != tt.wantDay {
				t.Errorf("Day: got %d, want %d", date.Day, tt.wantDay)
			}
		})
	}
}

// T013: TestParseDateYearOnly verifies parsing of year-only dates
func TestParseDateYearOnly(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantYear int
	}{
		{
			name:     "year only 2025",
			input:    "2025",
			wantYear: 2025,
		},
		{
			name:     "year only 2024",
			input:    "2024",
			wantYear: 2024,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When parsing year-only date
			date, err := ParseDate(tt.input)
			
			// Then no error should occur
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			
			// And year should match
			if date.Year != tt.wantYear {
				t.Errorf("Year: got %d, want %d", date.Year, tt.wantYear)
			}
			
			// And month should be 0 (not present)
			if date.Month != 0 {
				t.Errorf("Month: got %d, want 0 for year-only date", date.Month)
			}
			
			// And day should be 0 (not present)
			if date.Day != 0 {
				t.Errorf("Day: got %d, want 0 for year-only date", date.Day)
			}
		})
	}
}

// T013: TestParseDateYearMonth verifies parsing of year/month dates
func TestParseDateYearMonth(t *testing.T) {
	// Given a year/month format date
	input := "2025/09"
	
	// When parsing the date
	date, err := ParseDate(input)
	
	// Then no error should occur
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// And year should be correct
	if date.Year != 2025 {
		t.Errorf("Year: got %d, want 2025", date.Year)
	}
	
	// And month should be correct
	if date.Month != 9 {
		t.Errorf("Month: got %d, want 9", date.Month)
	}
	
	// And day should be 0 (not present)
	if date.Day != 0 {
		t.Errorf("Day: got %d, want 0", date.Day)
	}
}

// T013: TestParseDateEmpty verifies handling of empty dates
func TestParseDateEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty string", input: ""},
		{name: "whitespace", input: "   "},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When parsing empty date
			_, err := ParseDate(tt.input)
			
			// Then an error should occur
			if err == nil {
				t.Error("Expected error for empty date, got nil")
			}
		})
	}
}

// T013: TestParseDateInvalid verifies handling of invalid date formats
func TestParseDateInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "invalid format", input: "not-a-date"},
		{name: "invalid year", input: "abc/12/25"},
		{name: "invalid month", input: "2025/13/01"},
		{name: "invalid day", input: "2025/12/32"},
		{name: "old year", input: "1899/01/01"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When parsing invalid date
			_, err := ParseDate(tt.input)
			
			// Then an error should occur
			if err == nil {
				t.Errorf("Expected error for invalid date %q, got nil", tt.input)
			}
		})
	}
}
