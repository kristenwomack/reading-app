# Research: Style Books Table Layout

**Feature**: 002-style-the-table  
**Date**: 2025-01-16  
**Status**: Complete

## Research Questions

### 1. Layout Approach for 25%/75% Split

**Decision**: Use CSS Grid with `grid-template-columns: 1fr 3fr`

**Rationale**:
- Grid provides explicit 2-column layout control
- Fractional units (1fr + 3fr = 25%/75%) are more flexible than fixed percentages for gap handling
- Better browser support than newer layout features
- Simpler than Flexbox for fixed-ratio columns
- Built-in gap support eliminates manual margin calculations

**Alternatives Considered**:
- **Flexbox with flex-basis: 25%/75%**: Works but requires manual gap handling and less semantic for grid-like layouts
- **Float-based layout**: Legacy approach, harder to maintain, poor responsive behavior
- **CSS calc() with percentages**: Overly complex for simple 2-column split
- **Table display**: Semantically incorrect and accessibility concerns

**Implementation Notes**:
- Container: `.dashboard-layout { display: grid; grid-template-columns: 1fr 3fr; gap: 1rem; }`
- Mobile breakpoint: `@media (max-width: 768px) { grid-template-columns: 1fr; }` (stacks columns)

---

### 2. Responsive Breakpoint Strategy

**Decision**: Single breakpoint at 768px (tablet/mobile boundary)

**Rationale**:
- 768px is industry standard mobile/tablet breakpoint (matches iPad and most tablets in portrait)
- Single breakpoint keeps CSS simple (YAGNI principle)
- Below 768px: Stack vertically (summary table full width, then charts)
- Above 768px: Side-by-side 25%/75% layout
- Matches typical "mobile-first" responsive patterns

**Alternatives Considered**:
- **Multiple breakpoints (320px, 768px, 1024px, 1440px)**: Overengineering for current requirements
- **Container queries**: Too new, browser support insufficient (as of 2025)
- **600px breakpoint**: Too narrow, excludes tablets in portrait
- **No responsive design**: Violates FR-007 (mobile requirement)

**Implementation Notes**:
- Default (mobile-first): Stacked layout
- `@media (min-width: 769px)`: Apply 25%/75% grid
- Test on actual devices: iPhone SE (375px), iPad (768px), desktop (1440px+)

---

### 3. Summary Statistics Display Format

**Decision**: Compact card-style table with labeled stat rows

**Rationale**:
- Existing `/api/stats?year=YYYY` endpoint already provides aggregated data:
  ```json
  {
    "year": 2025,
    "total_books": 42,
    "total_pages": 12450,
    "avg_pages_per_book": 296
  }
  ```
- Card-style table (not HTML `<table>`) is more flexible for responsive design
- Each stat as a labeled row: "Total Books: 42" (semantic and accessible)
- No pagination or scrolling needed (only 3-5 stats total)

**Alternatives Considered**:
- **HTML table element**: Semantic but overkill for simple key-value display
- **Individual cards per stat**: Wastes vertical space, violates "condensed" requirement
- **List with bullets**: Less scannable than aligned labels
- **Inline text paragraph**: Poor readability and structure

**Data Mapping**:
| API Field | Display Label | Format |
|-----------|--------------|--------|
| `total_books` | "Total Books" | Integer |
| `total_pages` | "Total Pages" | Comma-separated (e.g., "12,450") |
| `avg_pages_per_book` | "Avg Pages/Book" | Integer |

**HTML Structure**:
```html
<div class="summary-card">
  <h2>2025 Summary</h2>
  <div class="stat-row">
    <span class="stat-label">Total Books</span>
    <span class="stat-value">42</span>
  </div>
  <!-- Repeat for other stats -->
</div>
```

---

### 4. Testing Strategy for Layout

**Decision**: Vitest with JSDOM for dimension assertions + manual visual testing

**Rationale**:
- Vitest already in use (frontend/package.json)
- JSDOM can assert CSS properties and computed styles
- Test grid-template-columns values at different viewport widths
- Manual testing on real devices for visual confirmation (automated visual regression out of scope)

**Test Cases**:
1. **Desktop layout (>768px)**: Assert `.dashboard-layout` has `grid-template-columns: 1fr 3fr`
2. **Mobile layout (<768px)**: Assert stacked layout (1fr single column)
3. **Summary card rendering**: Assert presence of stat rows with correct labels
4. **Empty state**: Assert "0 books" displays when API returns no data
5. **Width calculation**: Assert summary card width is approximately 25% of container (via offsetWidth)

**Alternatives Considered**:
- **Playwright/Cypress for visual testing**: Overkill for simple layout validation
- **Percy/Chromatic visual regression**: Adds external dependency (violates Constitution II)
- **Manual testing only**: No automated regression protection
- **Unit tests only (no integration)**: Misses layout calculation errors

**Implementation Notes**:
- Mock window.matchMedia in tests for breakpoint simulation
- Use `getComputedStyle()` to verify applied CSS rules
- Test data from mocked API responses

---

### 5. Chart Visualization Style

**Decision**: Use Chart.js line chart with connected data points for monthly reading trends

**Rationale**:
- Line charts better show trends over time compared to bar charts
- Connected data points emphasize reading progression throughout the year
- Chart.js already included in project (no new dependency)
- Line chart style from Chart.js documentation: https://www.chartjs.org/docs/latest/samples/line/line.html
- More visually elegant and better use of 75% screen space
- Easier to see patterns and compare months at a glance

**Alternatives Considered**:
- **Bar chart**: Original implementation, good for discrete comparisons but less effective for showing trends
- **Area chart**: Visually heavy, can obscure data points
- **Scatter plot**: Doesn't emphasize continuity between months
- **Multi-line chart**: Overkill for single metric (books read)

**Implementation Notes**:
- Chart type: `type: 'line'` in Chart.js config
- Dataset configuration:
  ```javascript
  {
    label: 'Books Read',
    data: monthlyData.map(m => m.Count),
    borderColor: 'rgb(75, 192, 192)',
    backgroundColor: 'rgba(75, 192, 192, 0.2)',
    tension: 0.1  // slight curve for smooth line
  }
  ```
- X-axis: Month names (Jan, Feb, Mar, etc.)
- Y-axis: Book count (starts at 0)
- Responsive: `responsive: true, maintainAspectRatio: false`

---

## Technology Stack Confirmation

**Frontend**:
- **Layout**: CSS Grid (native browser support)
- **Responsive**: CSS Media Queries (@media)
- **Styling**: Plain CSS (no preprocessors)
- **Testing**: Vitest + JSDOM (existing)
- **No new dependencies**: ✅

**Backend**: No changes required

**Contracts**: No changes required (existing `/api/stats` sufficient)

---

## Constitutional Compliance Summary

- ✅ **Monorepo**: Changes isolated to `frontend/` workspace
- ✅ **Minimal Dependencies**: Zero new dependencies added
- ✅ **TDD**: Tests written first (layout assertions before CSS implementation)
- ✅ **API Contracts**: No contract changes (reuses existing `/api/stats` endpoint)
- ✅ **Simplicity**: Single breakpoint, CSS Grid, no frameworks

---

## Open Questions & Risks

**None remaining.** All technical unknowns resolved.

**Risks**:
- **Browser compatibility**: CSS Grid is supported in all modern browsers (IE11 excluded, acceptable for 2025)
- **Chart library layout interactions**: Existing Chart.js may need container size adjustment (handled in implementation via resize observer or chart.resize() method)

---

## Next Steps

Proceed to **Phase 1: Design & Contracts**
- Extract data model (summary statistics entity)
- Define UI component structure (summary card)
- Write tests for layout behavior
- Create quickstart validation steps
