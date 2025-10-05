package books

// FilterByYear filters books by year from DateRead field
func FilterByYear(books []Book, year int) ([]Book, error) {
	var filtered []Book
	
	for _, book := range books {
		if book.DateRead == "" {
			continue
		}
		
		date, err := ParseDate(book.DateRead)
		if err != nil {
			continue // Skip invalid dates
		}
		
		if date.Year == year {
			filtered = append(filtered, book)
		}
	}
	
	return filtered, nil
}

// FilterByShelf filters books by shelf status
func FilterByShelf(books []Book, shelf string) []Book {
	var filtered []Book
	
	for _, book := range books {
		if book.Shelf == shelf {
			filtered = append(filtered, book)
		}
	}
	
	return filtered
}
