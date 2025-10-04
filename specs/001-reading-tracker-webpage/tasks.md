# Tasks: Reading Tracker Webpage

**Input**: Design documents from `/specs/001-reading-tracker-webpage/`
**Prerequisites**: plan.md ✓, research.md ✓, data-model.md ✓, contracts/openapi.yaml ✓, quickstart.md ✓

## Architecture Summary

- **Backend**: Go API with standard library net/http (3 REST endpoints)
- **Frontend**: Vanilla JavaScript with Chart.js
- **API Contract**: OpenAPI 3.0 specification in shared/contracts/
- **Data Source**: books.json (repository root)
- **Testing**: TDD enforced - Go testing + Vitest

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **Sequential**: Tasks without [P] must run in order
- **File paths**: Exact paths from repository root

## Path Conventions
```
backend/          # Go API workspace
frontend/         # JavaScript static site
shared/contracts/ # API specifications
books.json        # Data source (root)
```

---

## Phase 3.1: Setup

### Backend Setup
- [ ] **T001** Create backend directory structure: `backend/`, `backend/internal/books/`, `backend/internal/handlers/`, `backend/testdata/`
- [ ] **T002** Initialize Go module: `cd backend && go mod init github.com/kristenwomack/reading-app/backend`
- [ ] **T003** Create `backend/main.go` skeleton with package main, imports, and empty main() function
- [ ] **T004** Create `backend/README.md` with setup instructions (go version, how to run, how to test)

### Frontend Setup
- [ ] **T005** Create frontend directory structure: `frontend/`, `frontend/src/`, `frontend/styles/`, `frontend/tests/`
- [ ] **T006** Initialize `frontend/package.json` with Vitest dependency: `{"devDependencies": {"vitest": "^1.0.0"}}`
- [ ] **T007** Create `frontend/index.html` skeleton with DOCTYPE, head, body, and script module tag
- [ ] **T008** Create `frontend/README.md` with setup instructions (node version, npm install, npm test)

### Shared Setup
- [ ] **T009** Create shared directory: `shared/contracts/`
- [ ] **T010** Copy `specs/001-reading-tracker-webpage/contracts/openapi.yaml` to `shared/contracts/openapi.yaml`
- [ ] **T011** Create `shared/contracts/README.md` documenting the 3 API endpoints and their contracts

---

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation code**

### Backend Tests (Go)
All these tests can run in parallel [P] as they are in different packages:

- [ ] **T012** [P] Write test for loading books.json in `backend/internal/books/loader_test.go`:
  - TestLoadBooks: Verify books load from ../testdata/books_test.json
  - TestLoadBooksFileNotFound: Verify error when file missing
  - TestLoadBooksInvalidJSON: Verify error when JSON malformed
  
- [ ] **T013** [P] Write test for date parsing in `backend/internal/books/loader_test.go`:
  - TestParseDate: Test "2025/09/19" → year:2025, month:9, day:19
  - TestParseDateYearOnly: Test "2025" → year:2025, month:0, day:0
  - TestParseDateEmpty: Test "" → skip/error

- [ ] **T014** [P] Write test for filtering by year in `backend/internal/books/filter_test.go`:
  - TestFilterByYear: Given 100 books, filter by year 2025, expect only 2025 books
  - TestFilterByYearNoMatch: Filter by year with no books, expect empty array
  - TestFilterByShelf: Filter by shelf="read", expect only read books

- [ ] **T015** [P] Write test for statistics calculation in `backend/internal/books/stats_test.go`:
  - TestCalculateStatistics: Verify totalBooks, totalPages, averagePerMonth
  - TestCalculateMonthlyBreakdown: Verify 12 months with correct counts
  - TestCalculateStatisticsEmpty: Year with no books returns zeros

- [ ] **T016** [P] Write test for GET /api/years handler in `backend/internal/handlers/years_test.go`:
  - TestGetYears: HTTP GET returns 200 with JSON array of years
  - TestGetYearsFormat: Verify response matches YearsResponse schema

- [ ] **T017** [P] Write test for GET /api/books handler in `backend/internal/handlers/books_test.go`:
  - TestGetBooks: HTTP GET with ?year=2025 returns 200 with books array
  - TestGetBooksMissingYear: HTTP GET without year param returns 400
  - TestGetBooksInvalidYear: HTTP GET with invalid year returns 400

- [ ] **T018** [P] Write test for GET /api/stats handler in `backend/internal/handlers/stats_test.go`:
  - TestGetStats: HTTP GET with ?year=2025 returns 200 with statistics
  - TestGetStatsMissingYear: HTTP GET without year param returns 400
  - TestGetStatsFormat: Verify response matches StatsResponse schema

- [ ] **T019** Create test fixture `backend/testdata/books_test.json` with sample books for testing

### Frontend Tests (JavaScript)
All these tests can run in parallel [P] as they test different modules:

- [ ] **T020** [P] Write test for API client in `frontend/src/api-client.test.js`:
  - testFetchYears: Mock fetch, verify GET /api/years called correctly
  - testFetchBooks: Mock fetch with year param, verify response parsing
  - testFetchStats: Mock fetch with year param, verify response parsing
  - testFetchError: Mock fetch failure, verify error handling

- [ ] **T021** [P] Write test for chart rendering in `frontend/src/chart.test.js`:
  - testRenderChart: Given monthly data, verify Chart.js called with correct config
  - testRenderChartEmpty: Given zero data, verify chart shows empty state
  - testMonthLabels: Verify labels are ["Jan", "Feb", ..., "Dec"]

- [ ] **T022** [P] Write test for UI state in `frontend/src/ui.test.js`:
  - testPopulateYearSelector: Given years array, verify dropdown populated
  - testUpdateStatistics: Given stats object, verify DOM updated with numbers
  - testShowEmptyState: Verify empty state message displayed when no books
  - testShowError: Verify error message displayed on failure

---

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Backend Implementation (Sequential - dependency order)

- [ ] **T023** Implement `backend/internal/books/loader.go`:
  - LoadBooks() function to read and parse ../books.json
  - ParseDate() function to extract year/month/day from "YYYY/MM/DD"
  - Book struct matching books.json schema
  - Make T012 and T013 tests pass

- [ ] **T024** Implement `backend/internal/books/filter.go`:
  - FilterByYear() function to filter books array by year
  - FilterByShelf() function to filter by shelf status
  - Make T014 tests pass

- [ ] **T025** Implement `backend/internal/books/stats.go`:
  - CalculateStatistics() function returning totalBooks, totalPages, averagePerMonth
  - CalculateMonthlyBreakdown() function returning array of 12 MonthlyCount structs
  - Make T015 tests pass

- [ ] **T026** Implement `backend/internal/handlers/years.go`:
  - GetYears HTTP handler function
  - Extract unique years from books, count per year
  - Return JSON response matching YearsResponse schema
  - Make T016 tests pass

- [ ] **T027** Implement `backend/internal/handlers/books.go`:
  - GetBooks HTTP handler function
  - Parse year query parameter, validate
  - Call FilterByYear and FilterByShelf
  - Return JSON response matching BooksResponse schema
  - Make T017 tests pass

- [ ] **T028** Implement `backend/internal/handlers/stats.go`:
  - GetStats HTTP handler function
  - Parse year query parameter, validate
  - Call CalculateStatistics and CalculateMonthlyBreakdown
  - Return JSON response matching StatsResponse schema
  - Make T018 tests pass

- [ ] **T029** Implement `backend/main.go`:
  - Load books.json on startup into memory cache
  - Set up HTTP routes: /api/years, /api/books, /api/stats
  - Set up static file server for frontend: http.FileServer(http.Dir("../frontend"))
  - Start HTTP server on :8080
  - Add proper error handling and logging

### Frontend Implementation (After backend API is functional)

- [ ] **T030** Implement `frontend/src/api-client.js`:
  - fetchYears() async function calling GET /api/years
  - fetchBooks(year) async function calling GET /api/books?year={year}
  - fetchStats(year) async function calling GET /api/stats?year={year}
  - Error handling for fetch failures
  - Make T020 tests pass

- [ ] **T031** Implement `frontend/src/chart.js`:
  - renderChart(canvasElement, monthlyData) function
  - Configure Chart.js bar chart with monthly labels
  - Handle empty data case (no bars or all zero)
  - Make T021 tests pass

- [ ] **T032** Implement `frontend/src/ui.js`:
  - populateYearSelector(years) to fill dropdown with year options
  - updateStatistics(stats) to display numbers in DOM
  - showEmptyState(year) to display friendly "no books" message
  - showError(message) to display error state
  - Make T022 tests pass

- [ ] **T033** Implement `frontend/src/main.js`:
  - Initialize page on DOMContentLoaded
  - Call fetchYears(), populate year selector
  - Set default selected year (current year or latest with books)
  - Wire up year selector change event to reload data
  - Call fetchStats() and fetchBooks() for selected year
  - Call renderChart() with stats data
  - Handle errors and empty states

---

## Phase 3.4: Integration

- [ ] **T034** Create `frontend/styles/main.css`:
  - Layout: CSS Grid or Flexbox for main sections
  - Statistics cards: Clean, readable with good spacing
  - Chart container: Responsive sizing
  - Year selector: Dropdown styling
  - Empty state: Centered, friendly message
  - Error state: Visible but not alarming red

- [ ] **T035** Integrate Chart.js in `frontend/index.html`:
  - Add `<script src="https://cdn.jsdelivr.net/npm/chart.js@4"></script>` before app scripts
  - Add `<canvas id="monthlyChart"></canvas>` for chart rendering
  - Add module script tags for src/main.js

- [ ] **T036** Complete `frontend/index.html` with full page structure:
  - Header with title "Reading Tracker 2025"
  - Year selector dropdown
  - Statistics display area (total books, avg per month, total pages)
  - Chart container
  - Empty state container (hidden by default)
  - Error message container (hidden by default)

- [ ] **T037** Integration test: Start Go server with `go run backend/main.go`, open browser to http://localhost:8080:
  - Verify page loads
  - Verify year selector populates
  - Verify statistics display for current year
  - Verify chart renders with data
  - Verify switching years updates all data
  - Check browser console for errors

---

## Phase 3.5: Polish

These tasks can mostly run in parallel [P]:

- [ ] **T038** [P] Add empty state styling to `frontend/styles/main.css`:
  - Center-aligned message
  - Book icon or emoji
  - Muted colors
  - Make it friendly and encouraging

- [ ] **T039** [P] Add error state styling and handling:
  - Simple error message: "Unable to load reading data"
  - Clear visibility without being alarming
  - Log detailed error to console

- [ ] **T040** [P] Performance testing with 1000 books:
  - Measure page load time (should be <2 seconds)
  - Measure API response times (should be <100ms)
  - Use Chrome DevTools Performance tab
  - Document results in a comment

- [ ] **T041** [P] Update `backend/README.md` with complete instructions:
  - Prerequisites: Go 1.21+
  - How to run: `go run main.go`
  - How to test: `go test ./...`
  - API endpoints documentation
  - Architecture notes

- [ ] **T042** [P] Update `frontend/README.md` with complete instructions:
  - Prerequisites: Node.js 18+
  - How to install: `npm install`
  - How to test: `npm test`
  - How to access: http://localhost:8080 (served by Go backend)
  - Browser requirements: Chrome, Firefox, Safari, Edge (last 2 versions)

- [ ] **T043** Final integration test run - Execute all 10 scenarios from `specs/001-reading-tracker-webpage/quickstart.md`:
  - Scenario 1: Happy path with books in 2025
  - Scenario 2: Empty state (no books in year)
  - Scenario 3: Switching between years
  - Scenario 4: Data loading error
  - Scenario 5: Books with partial dates
  - Scenario 6: Books with zero pages
  - Scenario 7: Non-"read" shelf status filtering
  - Scenario 8: Large dataset (1000 books)
  - Scenario 9: Browser compatibility (test in all 4 browsers)
  - Scenario 10: Year selector population

---

## Dependencies

### Critical Dependencies (MUST follow this order):
1. **Setup (T001-T011)** → Must complete before any tests
2. **Tests (T012-T022)** → Must complete before implementation
3. **Backend Implementation (T023-T029)** → Sequential (dependency chain)
4. **Frontend Implementation (T030-T033)** → After backend API exists
5. **Integration (T034-T037)** → After both backend and frontend work
6. **Polish (T038-T043)** → After integration complete

### Parallel Execution Examples:

**Backend tests** (can all run simultaneously):
```bash
go test ./internal/books/... ./internal/handlers/... -v
# T012-T018 can execute in parallel
```

**Frontend tests** (can all run simultaneously):
```bash
npm test
# T020-T022 can execute in parallel
```

**Polish tasks** (independent work):
```bash
# T038, T039, T040, T041, T042 can all be done in parallel by different people
```

---

## Task Completion Summary

**Total Tasks**: 43
- Setup: 11 tasks (T001-T011)
- Tests (TDD): 11 tasks (T012-T022)
- Backend Implementation: 7 tasks (T023-T029)
- Frontend Implementation: 4 tasks (T030-T033)
- Integration: 4 tasks (T034-T037)
- Polish: 6 tasks (T038-T043)

**Estimated Execution**:
- With proper TDD: Tests fail first, then implement to make them pass
- Backend before frontend: API must exist for UI to consume
- Parallel where possible: Tests and polish tasks highly parallelizable

**Constitution Compliance**:
- ✅ TDD enforced (all tests before implementation)
- ✅ Minimal dependencies (Go stdlib, Chart.js only)
- ✅ Contract testing (OpenAPI spec, tests on both sides)
- ✅ Simple architecture (no frameworks, no database)

---

Ready to execute! Start with T001 and follow the order. Mark tasks with [x] as you complete them.
