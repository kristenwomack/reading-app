# Data Model: Reading Tracker Webpage

**Feature**: Reading Tracker Webpage  
**Date**: 2025-01-16  
**Phase**: 1 - Data Model Design

## Overview

This document defines the data structures for the reading tracker. Since this is a frontend-only feature reading from a static JSON file, the model focuses on the input schema (books.json) and the computed view models for statistics and visualization.

## Input Schema: books.json

### Book Entity

The existing books.json file contains an array of Book objects with the following structure:

```typescript
interface Book {
  Title: string;                    // Book title
  Author: string;                   // Primary author name
  "Additional Authors": string;     // Comma-separated additional authors (may be empty)
  ISBN: string;                     // ISBN-10 (may be empty)
  ISBN13: number | string;          // ISBN-13 (may be 0 or empty)
  Publisher: string;                // Publisher name (may be empty)
  "Number of Pages": number;        // Page count (may be 0)
  "Year Published": number;         // Publication year
  "Original Publication Year": number; // Original publication year
  "Date Read": string;              // Completion date in "YYYY/MM/DD" format (may be empty)
  "Date Added": string;             // Date added to list (may be empty)
  Bookshelves: string;              // Shelf name (e.g., "read", "currently-reading", "to-read")
  "Bookshelves with positions": string; // Extended shelf info (may be empty)
  Shelf: string;                    // Primary shelf status
  "My Review": string;              // User's review text (may be empty)
}
```

**Field Usage for Reading Tracker**:
- **Title**: ✓ Display in book list
- **Author**: ✓ Display in book list
- **Date Read**: ✓ PRIMARY FIELD - Filter by year, extract month
- **Shelf**: ✓ Filter for "read" status
- **Number of Pages**: ✓ Sum for total pages statistic
- Other fields: Not used in this feature

**Data Quality Notes**:
- Some books may have `"Number of Pages": 0` - exclude from page count calculations
- "Date Read" may be empty - exclude these books from tracker
- "Date Read" format is consistently "YYYY/MM/DD" (validated by existing data)

### Date Read Parsing

```typescript
interface ParsedDate {
  year: number;      // Always present (e.g., 2025)
  month: number;     // 1-12 (may be null for year-only dates)
  day: number;       // 1-31 (may be null)
  raw: string;       // Original date string
}
```

**Parsing Rules**:
1. Split by "/" delimiter
2. First part is year (required)
3. Second part is month (optional)
4. Third part is day (optional)
5. Convert strings to integers
6. Validate ranges: year > 1900, month 1-12, day 1-31

**Edge Case Handling**:
- `"2025"` → `{ year: 2025, month: null, day: null }`
- `"2025/09"` → `{ year: 2025, month: 9, day: null }`
- `"2025/09/19"` → `{ year: 2025, month: 9, day: 19 }`
- `""` or `null` → Skip book entirely

## Computed View Models

### YearStatistics

Aggregated statistics for a single year.

```typescript
interface YearStatistics {
  year: number;                      // Target year (e.g., 2025)
  totalBooks: number;                // Total books completed
  totalPages: number;                // Sum of all pages (excluding 0-page books)
  averageBooksPerMonth: number;      // totalBooks / 12, rounded to 1 decimal
  monthlyBreakdown: MonthlyCount[];  // Books per month (Jan-Dec)
  books: BookSummary[];              // List of books read this year
  isEmpty: boolean;                  // True if totalBooks === 0
}
```

### MonthlyCount

Count of books for a specific month, used in bar chart.

```typescript
interface MonthlyCount {
  month: number;          // 1-12 (January=1)
  monthName: string;      // "Jan", "Feb", ..., "Dec"
  count: number;          // Number of books completed
}
```

**Chart.js Mapping**:
- `labels`: Array of monthName values → ["Jan", "Feb", ..., "Dec"]
- `data`: Array of count values → [3, 5, 2, 0, ...]

### BookSummary

Simplified book data for display in the tracker.

```typescript
interface BookSummary {
  title: string;          // Book.Title
  author: string;         // Book.Author
  dateRead: string;       // Book["Date Read"]
  pages: number;          // Book["Number of Pages"]
  month: number | null;   // Extracted month (1-12) or null
}
```

**Display Usage**:
- Optional feature: Show list of books read
- Sorted by date read (most recent first)
- Group by month if desired

### YearOption

Available years for the year selector dropdown.

```typescript
interface YearOption {
  year: number;           // Year value
  label: string;          // Display label (e.g., "2025")
  bookCount: number;      // Number of books read this year
}
```

**Extraction Logic**:
1. Map all books → extract year from "Date Read"
2. Count books per year
3. Filter years with count > 0
4. Sort descending (newest first)

### UIState

Application state management (simple vanilla JS).

```typescript
interface UIState {
  selectedYear: number;         // Currently selected year
  availableYears: YearOption[]; // All years with books
  allBooks: Book[];             // Loaded from books.json
  currentStats: YearStatistics; // Computed for selectedYear
  isLoading: boolean;           // Data loading state
  hasError: boolean;            // Error loading data
  errorMessage: string;         // Error text to display
}
```

## Data Flow

```
1. Load books.json
   ↓
2. Parse JSON → Book[]
   ↓
3. Extract available years → YearOption[]
   ↓
4. Set selectedYear (default: current year or latest year with data)
   ↓
5. Filter books by selectedYear
   ↓
6. Filter by Shelf === "read"
   ↓
7. Parse "Date Read" for each book
   ↓
8. Calculate YearStatistics:
   - totalBooks: count filtered books
   - totalPages: sum pages (exclude 0)
   - averageBooksPerMonth: totalBooks / 12
   - monthlyBreakdown: group by month, count per month
   ↓
9. Render UI:
   - Statistics cards
   - Monthly bar chart
   - Year selector
   - Empty state (if no books)
```

## Validation Rules

### Book Filtering
- ✓ Include: `"Date Read"` contains selectedYear
- ✓ Include: `Shelf === "read"`
- ✗ Exclude: `"Date Read"` is empty or null
- ✗ Exclude: Year cannot be parsed
- ✗ Exclude: Shelf !== "read"

### Statistics Calculations
- **Total Books**: Count all books passing filters
- **Total Pages**: Sum `"Number of Pages"` where value > 0
- **Average Books/Month**: `totalBooks / 12`, round to 1 decimal place
- **Monthly Breakdown**: 
  - Initialize array with 12 entries (Jan-Dec), all counts = 0
  - For each book with valid month, increment count for that month
  - Books without month (year-only dates) not included in monthly breakdown

### Empty State Trigger
- Display when `totalBooks === 0` for selectedYear
- Show message: "No books tracked yet for [YEAR]. Start reading!"
- Statistics show zeros
- Chart hidden or displays empty bars

## Performance Considerations

**Dataset Size**: Up to 1000 books (per NFR-003)

**Operations**:
- Load JSON: O(1) - single fetch
- Parse JSON: O(n) - native JSON.parse
- Filter by year: O(n) - single pass
- Calculate statistics: O(n) - single pass
- Group by month: O(n) - single pass

**Total Complexity**: O(n) where n ≤ 1000

**Expected Performance**: <100ms for all computations (well under 2s requirement)

## Example Data Structure

```javascript
// Example YearStatistics for 2025
{
  year: 2025,
  totalBooks: 12,
  totalPages: 3240,
  averageBooksPerMonth: 1.0,
  isEmpty: false,
  monthlyBreakdown: [
    { month: 1, monthName: "Jan", count: 2 },
    { month: 2, monthName: "Feb", count: 1 },
    { month: 3, monthName: "Mar", count: 0 },
    // ... months 4-8
    { month: 9, monthName: "Sep", count: 3 },
    // ... months 10-12
  ],
  books: [
    {
      title: "Co-Intelligence",
      author: "Ethan Mollick",
      dateRead: "2025/09/19",
      pages: 0,
      month: 9
    },
    // ... more books
  ]
}
```

## Schema Validation

While not implemented in code (keeping it simple), the following validations should be considered during testing:

- books.json is valid JSON array
- Each book has required fields: Title, Author, "Date Read", Shelf
- "Date Read" format matches /^\d{4}(\/\d{2})?(\/\d{2})?$/
- "Number of Pages" is a number (including 0)

Invalid data should be logged to console but not break the application (graceful degradation).
