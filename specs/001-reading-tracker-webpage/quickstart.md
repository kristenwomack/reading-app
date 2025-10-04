# Quickstart: Reading Tracker Webpage

**Feature**: Reading Tracker Webpage  
**Date**: 2025-01-16  
**Purpose**: Test scenarios and integration examples

## Test Scenarios

### Scenario 1: Happy Path - Books Read in Current Year

**Given**:
- Current date is January 16, 2025
- books.json contains 12 books with "Date Read" in 2025
- All books have Shelf = "read"
- Page counts range from 200-400

**When**: User opens the reading tracker webpage

**Then**:
- Page displays "2025" as the selected year
- Year selector dropdown includes 2025 and other years with books
- Total books: 12
- Average books per month: 1.0
- Total pages: Sum of all page counts
- Monthly bar chart shows distribution across 12 months
- Each bar represents correct book count for that month
- Books are visible in a list (if implemented)

**Acceptance Criteria**:
- ✓ Statistics match calculations
- ✓ Chart renders with correct data
- ✓ Year selector is populated
- ✓ Page loads within 2 seconds

### Scenario 2: Empty State - No Books in Selected Year

**Given**:
- books.json contains books
- No books have "Date Read" in 2023
- User selects year 2023 from dropdown

**When**: User changes year selector to 2023

**Then**:
- Page displays "2023" as the selected year
- Empty state message appears: "No books tracked yet for 2023. Start reading!"
- Total books: 0
- Average books per month: 0.0
- Total pages: 0
- Monthly bar chart is hidden or shows all zero bars

**Acceptance Criteria**:
- ✓ Empty state message is friendly and visible
- ✓ Statistics show zeros
- ✓ No JavaScript errors
- ✓ Year selector still functional

### Scenario 3: Multiple Years - Switching Between Years

**Given**:
- books.json contains books from 2023, 2024, and 2025
- Each year has different book counts

**When**: 
1. Page loads (defaults to 2025)
2. User selects 2024 from dropdown
3. User selects 2023 from dropdown

**Then**:
- Each year selection updates all statistics
- Chart redraws with new data
- Statistics reflect correct year's books
- No page reload required (smooth update)
- Year selector shows current selection

**Acceptance Criteria**:
- ✓ Statistics update correctly for each year
- ✓ Chart updates smoothly
- ✓ No flicker or layout shifts
- ✓ Previous year data doesn't leak into new year

### Scenario 4: Data Loading Error

**Given**:
- books.json file is missing or inaccessible
- Network error or CORS issue prevents loading

**When**: User opens the reading tracker webpage

**Then**:
- Error message displays: "Unable to load reading data"
- No JavaScript console errors break the page
- Page remains functional (no crash)
- Year selector is disabled or empty
- Statistics section is hidden or shows error state

**Acceptance Criteria**:
- ✓ Error message is clear and visible
- ✓ No uncaught exceptions
- ✓ Page doesn't show misleading zero statistics
- ✓ User knows data failed to load

### Scenario 5: Books with Partial Dates

**Given**:
- books.json contains books with different date formats:
  - "2025/09/19" (full date)
  - "2025/09" (year and month only)
  - "2025" (year only)
- User selects year 2025

**When**: Page loads and calculates statistics

**Then**:
- All three books are counted in total books (all have year 2025)
- Full date book: Counted in September
- Year-month book: Counted in September
- Year-only book: Counted in yearly total but not in monthly breakdown
- Total books: 3
- September bar chart: 2 books (full date + year-month)
- Other months: Based on year-only books (implementation dependent)

**Acceptance Criteria**:
- ✓ All 2025 books included in total count
- ✓ Monthly breakdown handles partial dates gracefully
- ✓ No JavaScript errors from date parsing

### Scenario 6: Books with Zero Pages

**Given**:
- books.json contains 5 books in 2025
- 3 books have "Number of Pages" > 0 (total: 900 pages)
- 2 books have "Number of Pages" = 0

**When**: Page loads for year 2025

**Then**:
- Total books: 5 (all books counted)
- Total pages: 900 (zero-page books excluded)
- Average books per month: 0.4 (5 / 12)
- Chart shows all 5 books distributed by month

**Acceptance Criteria**:
- ✓ Zero-page books counted in book total
- ✓ Zero-page books excluded from page total
- ✓ Both statistics displayed correctly

### Scenario 7: Books with Non-"read" Shelf Status

**Given**:
- books.json contains books in 2025:
  - 8 books with Shelf = "read"
  - 3 books with Shelf = "currently-reading"
  - 2 books with Shelf = "to-read"

**When**: Page loads for year 2025

**Then**:
- Only "read" shelf books are counted
- Total books: 8 (excludes currently-reading and to-read)
- Statistics calculated from 8 books only
- Chart shows distribution of 8 books

**Acceptance Criteria**:
- ✓ Only "read" status books included
- ✓ Other shelf statuses filtered out
- ✓ Filtering logic documented in code

### Scenario 8: Large Dataset - 1000 Books

**Given**:
- books.json contains 1000 book entries
- Books distributed across multiple years
- Year 2025 has 200 books

**When**: Page loads with year 2025 selected

**Then**:
- Page loads within 2 seconds (per NFR-001)
- All 200 books for 2025 are processed
- Statistics calculate correctly
- Chart renders without lag
- Year selector includes all years (10+ years)
- No browser performance issues

**Acceptance Criteria**:
- ✓ Load time < 2 seconds
- ✓ No UI freeze or jank
- ✓ Memory usage reasonable
- ✓ Scrolling smooth (if book list shown)

### Scenario 9: Browser Compatibility

**Given**:
- Reading tracker deployed and accessible
- User has Chrome/Firefox/Safari/Edge (last 2 versions)

**When**: User opens page in each browser

**Then**:
- All browsers display identical functionality
- Chart renders correctly in all browsers
- ES6 features work (fetch, async/await, modules)
- CSS Grid/Flexbox layouts consistent
- No browser-specific errors

**Acceptance Criteria**:
- ✓ Chrome: Full functionality
- ✓ Firefox: Full functionality
- ✓ Safari: Full functionality
- ✓ Edge: Full functionality
- ✓ Visual consistency across browsers

### Scenario 10: Year Selector Population

**Given**:
- books.json contains books from years: 2020, 2022, 2024, 2025
- No books in 2021 or 2023

**When**: Page loads

**Then**:
- Year selector dropdown contains only: 2025, 2024, 2022, 2020
- Years listed in descending order (newest first)
- 2025 is selected by default (current year)
- Years without books (2021, 2023) are not in dropdown

**Acceptance Criteria**:
- ✓ Only years with books appear
- ✓ Years sorted descending
- ✓ Current year (or latest year) selected by default
- ✓ Dropdown is functional and accessible

## Integration Points

### 1. books.json File Access

**Location**: `/books.json` at repository root

**Access Method**: HTTP fetch request

**Requirements**:
- File must be served via HTTP (not file://)
- CORS not an issue (same origin)
- File must be valid JSON array

**Development Setup**:
```bash
# From repository root
cd reading-app
python3 -m http.server 8000

# Or with Node.js
npx http-server -p 8000

# Access at: http://localhost:8000/frontend/index.html
```

### 2. Frontend File Structure

**Entry Point**: `/frontend/index.html`

**Module Loading**:
```html
<script type="module" src="src/main.js"></script>
```

**ES6 Module Imports**:
```javascript
import { loadBooks } from './data-loader.js';
import { calculateStatistics } from './statistics.js';
import { renderChart } from './chart.js';
```

### 3. Chart.js Integration

**Option A - CDN** (Recommended for MVP):
```html
<script src="https://cdn.jsdelivr.net/npm/chart.js@4"></script>
```

**Option B - npm**:
```bash
npm install chart.js
```

### 4. Testing Integration

**Run Unit Tests**:
```bash
cd frontend
npm test
```

**Expected Test Structure**:
- `data-loader.test.js` → Test fetch and parsing
- `date-utils.test.js` → Test date extraction
- `statistics.test.js` → Test calculations
- `integration.test.js` → Test full data flow

**Coverage Goal**: 100% for new code (per TDD principle)

## Manual Testing Checklist

Before considering feature complete, manually verify:

- [ ] Page loads in all target browsers
- [ ] Year selector populates correctly
- [ ] Selecting different years updates all content
- [ ] Empty state shows for years without books
- [ ] Error state shows when books.json unavailable
- [ ] Chart displays correctly with accurate data
- [ ] Statistics calculations are accurate
- [ ] Page load time < 2 seconds (measure with DevTools)
- [ ] No console errors or warnings
- [ ] Layout is clean and readable
- [ ] Books with partial dates handled gracefully
- [ ] Zero-page books excluded from page totals
- [ ] Only "read" shelf books counted

## Performance Testing

**Dataset Sizes to Test**:
- Small: 10 books
- Medium: 100 books
- Large: 1000 books (required by NFR-003)

**Metrics to Measure** (using Chrome DevTools):
- Initial load time (target: < 2s)
- JSON parse time
- Statistics calculation time
- Chart render time
- Memory usage

**Tools**:
- Chrome DevTools → Performance tab
- Network tab for load times
- Console.time() for function timing

## Deployment Considerations

**Static Hosting Options**:
- GitHub Pages (frontend/ directory)
- Netlify (drop deployment)
- Vercel (static site)
- Simple HTTP server

**Required Files**:
- `/frontend/index.html`
- `/frontend/src/*.js`
- `/frontend/styles/main.css`
- `/books.json` (at root or accessible path)

**No Server Required**: Pure static files

**Future Backend Integration**: If books.json moves to API, update `data-loader.js` only (isolated change).
