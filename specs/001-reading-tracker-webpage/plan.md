# Implementation Plan: Reading Tracker Webpage

**Branch**: `001-reading-tracker-webpage` | **Date**: 2025-01-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-reading-tracker-webpage/spec.md`

## Summary

Build a reading progress tracker with a Go API backend serving book data and a static JavaScript frontend for visualization. The Go API reads from books.json and provides REST endpoints for the frontend to fetch filtered/aggregated data. Frontend displays monthly breakdown bar charts and comprehensive statistics for the selected year. Clean separation: Go handles data processing, JavaScript handles UI/visualization. No database, no authentication - simple file-based data access.

## Technical Context

**Backend**:
- **Language/Version**: Go (latest stable, 1.21+)
- **Framework**: Standard library net/http (no framework)
- **Storage**: Read books.json from repository root (in-memory caching optional)
- **API Style**: REST with JSON responses

**Frontend**:
- **Language/Version**: JavaScript ES6+ (vanilla JS) with HTML5/CSS3
- **Chart Library**: Chart.js (minimal, for bar chart visualization)
- **API Client**: Fetch API

**Testing**: 
- Backend: Go standard testing package
- Frontend: Vitest or browser-based testing

**Target Platform**: 
- Backend: Local development server (port 8080)
- Frontend: Static files served by Go or separate web server

**Project Type**: Web application (backend + frontend monorepo)

**Performance Goals**: 
- API response time: <100ms p95 for data endpoints
- Frontend load time: <2 seconds
- Handle up to 1000 book entries

**Constraints**: 
- No database (read from books.json file)
- No authentication (local/single-user app)
- Desktop browsers only (MVP)

**Scale/Scope**: Single-user local application, ~1000 books max

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
### I. Monorepo Structure ✅ PASS
- Backend workspace: `/backend` (Go service with net/http)
- Frontend workspace: `/frontend` (Static HTML/CSS/JS)
- Shared: `/books.json` at repository root (data source)
- Each workspace independently buildable:
  - Backend: `go build` with go.mod
  - Frontend: Static files (no build step)
- Clear interfaces: REST API contract between backend/frontend
- Each includes README with setup/testing instructions

### II. Minimal Dependencies ✅ PASS

**Backend (Go)**:
- Standard library net/http (no Gin, Echo, or other frameworks)
- Standard encoding/json for JSON handling
- Standard os/io for file reading
- No external dependencies beyond standard library

**Frontend (JavaScript)**:
- Vanilla JavaScript (no React, Vue, Angular)
- Chart.js only (for bar chart visualization)
- No build tools, no bundlers
- ES6 modules directly in browser

**Rationale**: Both backend and frontend use minimal dependencies aligned with constitutional principle II.

### III. Test-First Development (TDD) ✅ PASS
- Backend tests written first:
  - Book loading/parsing tests
  - Filtering logic tests (by year, by shelf)
  - Statistics calculation tests
  - HTTP handler tests
- Frontend tests written first:
  - API client tests
  - Chart rendering tests
  - UI state management tests
- User approves tests before implementation
- Red-Green-Refactor cycle enforced

### IV. API Contract Testing ✅ PASS
- REST API contracts defined in `/shared/contracts/`
- Contract includes:
  - GET /api/years - List available years
  - GET /api/books?year=2025 - Get books for year
  - GET /api/stats?year=2025 - Get statistics for year
- Backend validates request parameters
- Frontend validates response schemas
- Contract tests on both sides (backend serves, frontend consumes)
- OpenAPI/JSON schema documentation

### V. Simplicity & YAGNI ✅ PASS
- No database (reading from JSON file is sufficient)
- No authentication (single-user local app)
- Standard library HTTP server (no framework overhead)
- Simple REST endpoints (3 endpoints total)
- No caching layer initially (add only if performance requires)
- No WebSocket/real-time updates (not needed)
- Straightforward request-response pattern

**Overall Assessment**: ✅ ALL CONSTITUTIONAL PRINCIPLES SATISFIED

The addition of a Go backend actually improves separation of concerns:
- Backend: Data processing and business logic
- Frontend: Pure presentation and visualization
- Clean API contract between layers

## Project Structure

### Documentation (this feature)
```
specs/001-reading-tracker-webpage/
├── spec.md               # Feature specification (complete)
├── plan.md               # This file (/plan command output)
├── research.md           # Phase 0 output (to be updated)
├── data-model.md         # Phase 1 output (to be updated)
├── contracts/            # Phase 1 output (API contracts)
│   ├── openapi.yaml      # OpenAPI specification
│   └── examples/         # Request/response examples
├── quickstart.md         # Phase 1 output (test scenarios)
└── tasks.md              # Phase 2 output (/tasks command)
```

### Source Code (repository root)

**Structure Decision**: Web application (backend + frontend)

```
backend/                           # NEW - Go API workspace
├── go.mod                         # Go module definition
├── go.sum                         # Dependency checksums
├── README.md                      # Setup and testing instructions
├── main.go                        # HTTP server entry point
├── internal/
│   ├── books/
│   │   ├── loader.go              # Load books.json
│   │   ├── loader_test.go         # TDD: Load tests
│   │   ├── filter.go              # Filter by year/shelf
│   │   ├── filter_test.go         # TDD: Filter tests
│   │   └── stats.go               # Calculate statistics
│   │       └── stats_test.go      # TDD: Stats tests
│   └── handlers/
│       ├── years.go               # GET /api/years handler
│       ├── years_test.go          # TDD: Handler tests
│       ├── books.go               # GET /api/books handler
│       ├── books_test.go          # TDD: Handler tests
│       ├── stats.go               # GET /api/stats handler
│       └── stats_test.go          # TDD: Handler tests
└── testdata/
    └── books_test.json            # Test fixtures

frontend/                          # NEW - Static site workspace
├── package.json                   # npm for Vitest only
├── README.md                      # Setup instructions
├── index.html                     # Main reading tracker page
├── src/
│   ├── main.js                    # Entry point, initialization
│   ├── api-client.js              # Fetch from Go API
│   ├── api-client.test.js         # TDD: API client tests
│   ├── chart.js                   # Chart.js integration
│   ├── chart.test.js              # TDD: Chart tests
│   └── ui.js                      # DOM manipulation, year selector
│       └── ui.test.js             # TDD: UI tests
├── styles/
│   └── main.css                   # Clean, readable styling
└── tests/
    └── integration.test.js        # Full API→UI flow tests

shared/                            # NEW - Shared contracts
└── contracts/
    ├── openapi.yaml               # API specification
    └── README.md                  # Contract documentation

books.json                         # Existing data source (repository root)
```

**Rationale**: 
- Clean separation: Go backend for data logic, JS frontend for presentation
- Follows monorepo principle with 3 independent workspaces
- Backend reads books.json, processes data, exposes REST API
- Frontend is pure client-side code consuming API
- Shared contracts directory for API documentation
- No database needed (books.json is the data source)
- Both workspaces independently testable

## Phase 0: Outline & Research

✅ **COMPLETE** - See `research.md` (updated for Go API architecture)

**Key Decisions Made**:
1. Backend Framework: Go standard library net/http (no framework)
2. API Design: Simple REST with query parameters (3 endpoints)
3. Data Loading: Load books.json on startup, cache in memory
4. CORS: Same-origin (Go serves frontend static files)
5. Frontend Chart: Chart.js (minimal, well-maintained)
6. Date Parsing: Go backend (better separation of concerns)
7. Testing: Go `testing` package + Vitest for frontend
8. Static Serving: Go http.FileServer (single server)
9. Error Handling: JSON errors with HTTP status codes
10. Development: Integrated workflow (Go serves everything)

**Architecture**: Go API backend + static JavaScript frontend, single server deployment

**No Unknowns Remain**: All technical decisions finalized

## Phase 1: Design & Contracts

✅ **COMPLETE**

### 1. Data Model → `data-model.md`

**Note**: Existing data-model.md remains valid but perspective shifts:
- Backend handles: Book loading, filtering, statistics calculation
- Frontend handles: UI state, chart data formatting
- Clear separation of concerns via API contract

**Backend Models** (Go):
- Book struct (maps to books.json)
- YearStatistics (computed)
- MonthlyCount (aggregated)

**Frontend Models** (JavaScript):
- Receives structured data from API
- Minimal local state for UI

### 2. API Contracts → `contracts/openapi.yaml`

✅ **COMPLETE** - OpenAPI 3.0 specification created

**Endpoints Defined**:
1. `GET /api/years` → Returns available years with book counts
2. `GET /api/books?year={year}` → Returns books for specified year
3. `GET /api/stats?year={year}` → Returns comprehensive statistics

**Schemas Defined**:
- YearsResponse, YearOption
- BooksResponse, BookSummary
- StatsResponse, MonthlyCount
- ErrorResponse

**Examples**: Included for typical and empty state responses

**Contract Testing Strategy**:
- Backend tests validate response schemas
- Frontend tests validate request/response handling
- OpenAPI spec is source of truth

### 3. Test Scenarios → `quickstart.md`

**Note**: Existing quickstart.md scenarios remain valid, execution shifts:
- Backend: Test Go handlers, business logic, data processing
- Frontend: Test API client, UI rendering, chart integration
- Integration: Test full stack (Go API → Frontend)

**10 Scenarios** still applicable with backend/frontend split

### 4. Agent Context Update

✅ **COMPLETE** - Updated `.github/copilot-instructions.md`

**Changes Applied**:
- Added Go backend with standard library
- Updated project structure (backend/ + frontend/ + shared/)
- Added API endpoint documentation
- Updated commands for Go + JS development

**Output**: 
- data-model.md ✓ (existing, valid)
- contracts/openapi.yaml ✓ (new)
- quickstart.md ✓ (existing, valid)
- .github/copilot-instructions.md ✓ (updated)

## Constitution Check (Re-evaluation)

### Post-Design Validation

✅ **ALL PRINCIPLES STILL SATISFIED** - Actually improved with Go backend!

**I. Monorepo Structure**: ✓ Enhanced
- Backend workspace: /backend with go.mod
- Frontend workspace: /frontend with package.json
- Shared workspace: /shared/contracts
- Each independently buildable and testable
- Clear API contract between workspaces

**II. Minimal Dependencies**: ✓ Maintained
- Backend: Go standard library only (net/http, encoding/json)
- Frontend: Chart.js only
- No build dependencies (Go compiles, JS uses ES6 modules)
- No frameworks on either side

**III. Test-First Development (TDD)**: ✓ Enhanced
- Backend: Table-driven tests for all logic
- Frontend: Vitest tests for UI and API client
- Contract tests on both sides
- Full integration test scenarios
- 10 comprehensive test scenarios in quickstart.md

**IV. API Contract Testing**: ✓ Fully Implemented
- OpenAPI 3.0 specification defined
- 3 REST endpoints documented
- Request/response schemas validated
- Examples for all scenarios
- Backend and frontend test against contract

**V. Simplicity & YAGNI**: ✓ Maintained
- No database (file-based)
- No authentication (local app)
- Simple REST API (3 endpoints)
- Standard library HTTP (no framework)
- Straightforward request-response pattern

**Final Assessment**: ✅ PASS - Go backend addition actually strengthens architecture by providing clear separation of concerns while maintaining all constitutional principles

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:

### Setup Tasks (Phase 3.1)
**Backend Setup**:
1. Create /backend directory structure
2. Initialize go.mod
3. Create main.go skeleton with HTTP server setup
4. Set up internal/ package structure (books/, handlers/)

**Frontend Setup**:
5. Create /frontend directory structure
6. Initialize package.json (Vitest only)
7. Create index.html skeleton
8. Set up src/ structure

**Shared Setup**:
9. Copy contracts/openapi.yaml to /shared/contracts/

### Test Tasks (Phase 3.2 - TDD)
**Backend Tests** (can run in parallel [P]):
10. [P] Test book loading from JSON (internal/books/loader_test.go)
11. [P] Test date parsing logic (internal/books/loader_test.go)
12. [P] Test filtering by year and shelf (internal/books/filter_test.go)
13. [P] Test statistics calculations (internal/books/stats_test.go)
14. [P] Test GET /api/years handler (internal/handlers/years_test.go)
15. [P] Test GET /api/books handler (internal/handlers/books_test.go)
16. [P] Test GET /api/stats handler (internal/handlers/stats_test.go)

**Frontend Tests** (can run in parallel [P]):
17. [P] Test API client (src/api-client.test.js)
18. [P] Test chart rendering logic (src/chart.test.js)
19. [P] Test UI state management (src/ui.test.js)

### Core Implementation Tasks (Phase 3.3)
**Backend Implementation** (sequential, TDD order):
20. Implement internal/books/loader.go (load books.json, parse dates)
21. Implement internal/books/filter.go (filter by year/shelf)
22. Implement internal/books/stats.go (calculate statistics)
23. Implement internal/handlers/years.go (GET /api/years)
24. Implement internal/handlers/books.go (GET /api/books)
25. Implement internal/handlers/stats.go (GET /api/stats)
26. Implement main.go (wire up routes, serve static files)

**Frontend Implementation** (after backend API available):
27. Implement src/api-client.js (fetch from Go API)
28. Implement src/chart.js (Chart.js integration)
29. Implement src/ui.js (DOM manipulation, year selector)
30. Implement src/main.js (orchestration, page initialization)

### Integration Tasks (Phase 3.4)
31. Create frontend/styles/main.css (clean, readable styling)
32. Integrate Chart.js (CDN link in index.html)
33. Wire up frontend/index.html with all modules
34. Integration test: Start Go server, test full stack
35. Manual testing in all browsers (Chrome, Firefox, Safari, Edge)

### Polish Tasks (Phase 3.5 - Parallel where possible)
36. [P] Empty state styling and messages (frontend)
37. [P] Error state handling and display (frontend)
38. [P] Performance testing with 1000 books (measure load times)
39. [P] backend/README.md with setup and testing instructions
40. [P] frontend/README.md with setup instructions
41. Final integration test run (all 10 scenarios from quickstart.md)

**Ordering Strategy**:
- TDD strictly enforced: All tests before implementation
- Backend first: API must be running before frontend can consume it
- Dependencies respected: loader → filter → stats → handlers → main
- Frontend after backend: API client needs working endpoints
- Parallel tasks marked [P]: Independent test files, documentation, styling

**File Impact Analysis**:
- Backend tests: Different packages = can run [P]
- Frontend tests: Different modules = can run [P]
- Backend implementation: Mostly sequential (dependency chain)
- Frontend implementation: After backend API is functional
- Polish tasks: Highly parallel (independent concerns)

**Estimated Output**: ~40 numbered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

**No violations to document** - All constitutional principles satisfied.

---

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) → research.md (updated for Go API)
- [x] Phase 1: Design complete (/plan command) → data-model.md, contracts/openapi.yaml, quickstart.md, copilot-instructions.md
- [x] Phase 2: Task planning complete (/plan command - approach described for Go + JS)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS (Go backend enhances architecture)
- [x] Post-Design Constitution Check: PASS (all principles strengthened)
- [x] All NEEDS CLARIFICATION resolved (via /clarify)
- [x] Complexity deviations documented (none - Go adds clarity, not complexity)

**Artifacts Generated**:
- [x] specs/001-reading-tracker-webpage/plan.md (this file - updated)
- [x] specs/001-reading-tracker-webpage/research.md (updated for Go API)
- [x] specs/001-reading-tracker-webpage/data-model.md (valid for both backend/frontend)
- [x] specs/001-reading-tracker-webpage/contracts/openapi.yaml (NEW)
- [x] specs/001-reading-tracker-webpage/quickstart.md (valid for full stack)
- [x] .github/copilot-instructions.md (updated with Go + JS)

**Next Command**: `/tasks` - Generate ~40 tasks for Go backend + JS frontend implementation

---

*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
