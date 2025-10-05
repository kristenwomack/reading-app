package books

import (
	"testing"
)

// T015: TestCalculateStatistics verifies basic statistics calculation
func TestCalculateStatistics(t *testing.T) {
	// Given books for a year
	books := []Book{
		{Title: "Book 1", Pages: float64(300), DateRead: "2025/01/15"},
		{Title: "Book 2", Pages: float64(250), DateRead: "2025/02/20"},
		{Title: "Book 3", Pages: float64(400), DateRead: "2025/03/10"},
	}
	
	// When calculating statistics
	stats := CalculateStatistics(books, 2025)
	
	// Then total books should be 3
	if stats.TotalBooks != 3 {
		t.Errorf("TotalBooks: got %d, want 3", stats.TotalBooks)
	}
	
	// And total pages should be 950
	if stats.TotalPages != 950 {
		t.Errorf("TotalPages: got %d, want 950", stats.TotalPages)
	}
	
	// And average per month should be 0.25 (3 books / 12 months)
	expected := 3.0 / 12.0
	if stats.AveragePerMonth != expected {
		t.Errorf("AveragePerMonth: got %.2f, want %.2f", stats.AveragePerMonth, expected)
	}
}

// T015: TestCalculateStatisticsExcludesZeroPages verifies zero-page books excluded from total
func TestCalculateStatisticsExcludesZeroPages(t *testing.T) {
	// Given books with some having zero pages
	books := []Book{
		{Title: "Book 1", Pages: float64(300), DateRead: "2025/01/15"},
		{Title: "Book 2", Pages: 0, DateRead: "2025/02/20"},   // Zero pages
		{Title: "Book 3", Pages: float64(250), DateRead: "2025/03/10"},
		{Title: "Book 4", Pages: 0, DateRead: "2025/04/05"},   // Zero pages
	}
	
	// When calculating statistics
	stats := CalculateStatistics(books, 2025)
	
	// Then total books should include all 4
	if stats.TotalBooks != 4 {
		t.Errorf("TotalBooks: got %d, want 4", stats.TotalBooks)
	}
	
	// But total pages should exclude zero-page books
	if stats.TotalPages != 550 {
		t.Errorf("TotalPages: got %d, want 550 (300+250)", stats.TotalPages)
	}
}

// T015: TestCalculateMonthlyBreakdown verifies monthly aggregation
func TestCalculateMonthlyBreakdown(t *testing.T) {
	// Given books distributed across months
	books := []Book{
		{Title: "Jan Book 1", DateRead: "2025/01/05"},
		{Title: "Jan Book 2", DateRead: "2025/01/20"},
		{Title: "Feb Book 1", DateRead: "2025/02/15"},
		{Title: "Mar Book 1", DateRead: "2025/03/10"},
		{Title: "Mar Book 2", DateRead: "2025/03/25"},
		{Title: "Mar Book 3", DateRead: "2025/03/30"},
		// April-Dec have no books
	}
	
	// When calculating monthly breakdown
	breakdown := CalculateMonthlyBreakdown(books)
	
	// Then should have 12 months
	if len(breakdown) != 12 {
		t.Errorf("Expected 12 months, got %d", len(breakdown))
	}
	
	// And January should have 2 books
	if breakdown[0].Count != 2 {
		t.Errorf("January count: got %d, want 2", breakdown[0].Count)
	}
	if breakdown[0].MonthName != "Jan" {
		t.Errorf("January name: got %q, want 'Jan'", breakdown[0].MonthName)
	}
	
	// And February should have 1 book
	if breakdown[1].Count != 1 {
		t.Errorf("February count: got %d, want 1", breakdown[1].Count)
	}
	
	// And March should have 3 books
	if breakdown[2].Count != 3 {
		t.Errorf("March count: got %d, want 3", breakdown[2].Count)
	}
	
	// And remaining months should have 0 books
	for i := 3; i < 12; i++ {
		if breakdown[i].Count != 0 {
			t.Errorf("Month %d count: got %d, want 0", i+1, breakdown[i].Count)
		}
	}
}

// T015: TestCalculateStatisticsEmpty verifies handling of no books
func TestCalculateStatisticsEmpty(t *testing.T) {
	// Given no books
	books := []Book{}
	
	// When calculating statistics
	stats := CalculateStatistics(books, 2025)
	
	// Then all values should be zero
	if stats.TotalBooks != 0 {
		t.Errorf("TotalBooks: got %d, want 0", stats.TotalBooks)
	}
	
	if stats.TotalPages != 0 {
		t.Errorf("TotalPages: got %d, want 0", stats.TotalPages)
	}
	
	if stats.AveragePerMonth != 0 {
		t.Errorf("AveragePerMonth: got %.2f, want 0", stats.AveragePerMonth)
	}
}

// T015: TestCalculateMonthlyBreakdownEmpty verifies empty breakdown
func TestCalculateMonthlyBreakdownEmpty(t *testing.T) {
	// Given no books
	books := []Book{}
	
	// When calculating monthly breakdown
	breakdown := CalculateMonthlyBreakdown(books)
	
	// Then should still have 12 months
	if len(breakdown) != 12 {
		t.Errorf("Expected 12 months, got %d", len(breakdown))
	}
	
	// And all counts should be zero
	for i, month := range breakdown {
		if month.Count != 0 {
			t.Errorf("Month %d count: got %d, want 0", i+1, month.Count)
		}
	}
}

// T015: TestMonthNamesCorrect verifies month names are correct
func TestMonthNamesCorrect(t *testing.T) {
	expectedNames := []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}
	
	// When calculating breakdown with any books
	books := []Book{{Title: "Test", DateRead: "2025/01/01"}}
	breakdown := CalculateMonthlyBreakdown(books)
	
	// Then month names should match
	for i, expected := range expectedNames {
		if breakdown[i].MonthName != expected {
			t.Errorf("Month %d name: got %q, want %q", i+1, breakdown[i].MonthName, expected)
		}
		if breakdown[i].Month != i+1 {
			t.Errorf("Month %d number: got %d, want %d", i+1, breakdown[i].Month, i+1)
		}
	}
}
