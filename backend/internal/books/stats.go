package books

// Statistics represents reading statistics for a year
type Statistics struct {
	Year            int
	TotalBooks      int
	TotalPages      int
	AveragePerMonth float64
}

// MonthlyCount represents book count for a month
type MonthlyCount struct {
	Month     int
	MonthName string
	Count     int
}

// CalculateStatistics calculates reading statistics for a year
func CalculateStatistics(books []Book, year int) Statistics {
	stats := Statistics{
		Year:       year,
		TotalBooks: len(books),
	}
	
	// Calculate total pages (exclude zero-page books)
	for _, book := range books {
		pages := book.GetPages()
		if pages > 0 {
			stats.TotalPages += pages
		}
	}
	
	// Calculate average per month
	if stats.TotalBooks > 0 {
		stats.AveragePerMonth = float64(stats.TotalBooks) / 12.0
	}
	
	return stats
}

// CalculateMonthlyBreakdown calculates book count per month
func CalculateMonthlyBreakdown(books []Book) []MonthlyCount {
	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	
	// Initialize all 12 months
	breakdown := make([]MonthlyCount, 12)
	for i := 0; i < 12; i++ {
		breakdown[i] = MonthlyCount{
			Month:     i + 1,
			MonthName: monthNames[i],
			Count:     0,
		}
	}
	
	// Count books per month
	for _, book := range books {
		if book.DateRead == "" {
			continue
		}
		
		date, err := ParseDate(book.DateRead)
		if err != nil || date.Month == 0 {
			continue
		}
		
		if date.Month >= 1 && date.Month <= 12 {
			breakdown[date.Month-1].Count++
		}
	}
	
	return breakdown
}
