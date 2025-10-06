# Data Model: Summary Statistics Display

**Feature**: 002-style-the-table  
**Date**: 2025-01-16  
**Status**: Complete

## Entities

### 1. ReadingStatistics

**Description**: Aggregated reading data for a specific year, displayed in the summary table.

**Source**: Backend `/api/stats?year={year}` endpoint (existing, no changes)

**Fields**:
| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `year` | integer | Required, 4-digit | Year for which statistics are calculated |
| `total_books` | integer | Required, >= 0 | Total number of books finished in the year |
| `total_pages` | integer | Required, >= 0 | Sum of page counts across all books |
| `avg_pages_per_book` | integer | Required, >= 0 | Average pages per book (rounded) |

**Validation Rules**:
- `year` must be a valid 4-digit year (e.g., 2025)
- `total_books` = 0 is valid (empty state)
- If `total_books` = 0, then `total_pages` = 0 and `avg_pages_per_book` = 0
- `avg_pages_per_book` = floor(`total_pages` / `total_books`) when `total_books` > 0

**Example Instances**:
```json
// Normal year with books
{
  "year": 2025,
  "total_books": 42,
  "total_pages": 12450,
  "avg_pages_per_book": 296
}

// Empty year (no books)
{
  "year": 2024,
  "total_books": 0,
  "total_pages": 0,
  "avg_pages_per_book": 0
}
```

**Relationships**: None (standalone entity, not related to individual Book entities)

---

## UI Components

### 1. SummaryCard Component

**Purpose**: Display ReadingStatistics in a condensed, scannable format occupying 25% of screen width (desktop).

**HTML Structure**:
```html
<div class="summary-card" data-testid="summary-card">
  <h2 class="summary-title">2025 Summary</h2>
  <div class="stats-list">
    <div class="stat-row" data-testid="stat-total-books">
      <span class="stat-label">Total Books</span>
      <span class="stat-value">42</span>
    </div>
    <div class="stat-row" data-testid="stat-total-pages">
      <span class="stat-label">Total Pages</span>
      <span class="stat-value">12,450</span>
    </div>
    <div class="stat-row" data-testid="stat-avg-pages">
      <span class="stat-label">Avg Pages/Book</span>
      <span class="stat-value">296</span>
    </div>
  </div>
</div>
```

**CSS Classes** (defined in `frontend/styles/main.css`):
- `.summary-card`: Container styling (border, padding, background)
- `.summary-title`: Year heading (e.g., "2025 Summary")
- `.stats-list`: Wrapper for stat rows (vertical stack)
- `.stat-row`: Individual stat container (flexbox for label/value alignment)
- `.stat-label`: Left-aligned label text
- `.stat-value`: Right-aligned numeric value (bold/emphasized)

**Empty State** (when `total_books` = 0):
```html
<div class="summary-card empty-state" data-testid="summary-card-empty">
  <h2 class="summary-title">2025 Summary</h2>
  <p class="empty-message">No books tracked for this year</p>
</div>
```

**State Transitions**:
1. **Loading**: Show skeleton or spinner (out of scope for this feature)
2. **Populated**: Display stat rows (normal state)
3. **Empty**: Show empty state message

---

### 2. DashboardLayout Component

**Purpose**: Container for 25%/75% grid layout (summary card + charts section).

**HTML Structure**:
```html
<div class="dashboard-layout" data-testid="dashboard-layout">
  <aside class="summary-section">
    <!-- SummaryCard component here -->
  </aside>
  <main class="charts-section">
    <!-- Existing charts (Chart.js canvas elements) -->
  </main>
</div>
```

**CSS Layout** (defined in `frontend/styles/main.css`):
```css
/* Desktop: 25% summary, 75% charts */
@media (min-width: 769px) {
  .dashboard-layout {
    display: grid;
    grid-template-columns: 1fr 3fr; /* 25% / 75% ratio */
    gap: 1.5rem;
    align-items: start;
  }
}

/* Mobile: Stacked (summary full width, then charts) */
@media (max-width: 768px) {
  .dashboard-layout {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .summary-section {
    width: 100%;
  }
}
```

---

## Data Flow

### Rendering Flow
1. User selects year (via year selector UI, existing functionality)
2. Frontend calls `/api/stats?year=2025` (existing API client)
3. API returns `ReadingStatistics` JSON
4. UI.js updates SummaryCard with stat values
5. CSS Grid layout positions SummaryCard at 25% width (desktop) or full width (mobile)

**No changes to data fetching logic** - existing `api-client.js` already handles stats endpoint.

---

## Display Formatting

### Number Formatting Rules
| Field | Format | Example Input | Example Output |
|-------|--------|---------------|----------------|
| `total_books` | Integer | 42 | "42" |
| `total_pages` | Comma-separated | 12450 | "12,450" |
| `avg_pages_per_book` | Integer | 296 | "296" |

**Implementation**:
```javascript
// In ui.js
function formatNumber(num) {
  return num.toLocaleString('en-US');
}
```

---

## Accessibility Considerations

- **Semantic HTML**: Use `<aside>` for summary section (complementary content)
- **Heading Hierarchy**: `<h2>` for summary title (assuming page has `<h1>`)
- **Screen Reader Labels**: Ensure stat labels are readable ("Total Books: 42")
- **Focus Management**: Summary card should not interfere with keyboard navigation to charts

---

## Testing Data

**Test Cases**:
1. **Normal year**: `{ year: 2025, total_books: 42, total_pages: 12450, avg_pages_per_book: 296 }`
2. **Empty year**: `{ year: 2024, total_books: 0, total_pages: 0, avg_pages_per_book: 0 }`
3. **Single book**: `{ year: 2023, total_books: 1, total_pages: 320, avg_pages_per_book: 320 }`
4. **Large numbers**: `{ year: 2022, total_books: 150, total_pages: 45000, avg_pages_per_book: 300 }`

---

## Summary

**Entities**: 1 (ReadingStatistics - existing backend model, no changes)  
**UI Components**: 2 (SummaryCard, DashboardLayout)  
**API Changes**: None (reuses existing `/api/stats` endpoint)  
**New Fields**: None (all data already available)  
**Database Changes**: None (data derived from existing books.json)

This feature is purely a UI/UX layout change with no backend or data model modifications.
