# Research: Reading Tracker Webpage

**Feature**: Reading Tracker Webpage  
**Date**: 2025-01-16 (Updated)  
**Phase**: 0 - Research & Decisions
**Architecture**: Go API Backend + Static JavaScript Frontend

## Research Items

### 1. Backend Framework Choice

**Decision**: Go standard library net/http (no framework)

**Rationale**:
- Aligns with Constitution Principle II (Minimal Dependencies)
- Standard library net/http is production-ready and performant
- No need for Gin, Echo, or other frameworks for 3 simple endpoints
- Zero external dependencies for HTTP server
- Simpler codebase, easier to understand and maintain

**Alternatives Considered**:
- **Gin/Echo**: Too heavyweight for 3 endpoints, adds unnecessary complexity
- **chi/gorilla/mux**: Minimal routers, but standard library suffices for our needs
- **Standard library**: ✅ Sufficient, simple, zero dependencies

**Implementation Approach**:
```go
http.HandleFunc("/api/years", handlers.GetYears)
http.HandleFunc("/api/books", handlers.GetBooks)
http.HandleFunc("/api/stats", handlers.GetStats)
http.HandleFunc("/", handlers.ServeStatic) // Serve frontend
```

### 2. API Design Pattern

**Decision**: Simple REST with query parameters

**Endpoints**:
- `GET /api/years` - Returns array of years with book counts
- `GET /api/books?year=2025` - Returns books for specified year
- `GET /api/stats?year=2025` - Returns statistics for specified year

**Rationale**:
- RESTful and intuitive
- Query parameters sufficient (no path parameters needed)
- JSON responses throughout
- No versioning needed (v1 built into path if needed later)
- Follows YAGNI principle

**Response Format**:
```json
// GET /api/years
{"years": [{"year": 2025, "count": 12}, {"year": 2024, "count": 45}]}

// GET /api/books?year=2025
{"books": [{"title": "...", "author": "...", "dateRead": "2025/09/19", ...}]}

// GET /api/stats?year=2025
{"year": 2025, "totalBooks": 12, "totalPages": 3240, "averagePerMonth": 1.0, 
 "monthlyBreakdown": [{"month": 1, "count": 2}, ...]}
```

### 3. Data Loading Strategy

**Decision**: Load books.json on startup, cache in memory

**Rationale**:
- File doesn't change during server runtime (static data)
- Loading once avoids repeated disk I/O
- 1000 books ~600KB JSON, easily fits in memory
- Restart server to reload data (acceptable for local use)
- Simple, performant, no cache invalidation complexity

**Implementation**:
```go
var booksCache []Book

func init() {
  data, _ := os.ReadFile("../books.json")
  json.Unmarshal(data, &booksCache)
}
```

**Alternative**: Re-read file on each request (slower, unnecessary I/O)

### 4. CORS Configuration

**Decision**: Allow localhost origins for development

**Rationale**:
- Frontend served from same origin (Go server serves static files) = No CORS needed
- Alternative: If frontend served separately (e.g., port 3000), enable CORS for localhost

**Implementation** (if needed):
```go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
```

**Preferred**: Serve frontend from same Go server (no CORS issues)

### 5. Frontend Chart Library

**Decision**: Chart.js (minimal, well-maintained)

**Rationale**:
- Lightweight (~60KB minified)
- Active maintenance ✓
- Simple API for bar charts
- Works with vanilla JavaScript
- No framework dependencies

**Implementation**: Include via CDN in HTML:
```html
<script src="https://cdn.jsdelivr.net/npm/chart.js@4"></script>
```

### 6. Date Parsing Location

**Decision**: Parse dates in Go backend, send structured data to frontend

**Rationale**:
- Backend is better suited for data processing
- Frontend receives clean, structured data
- No duplicate parsing logic
- Backend can validate date format
- Simpler frontend code (pure presentation)

**Backend Parsing** (Go):
```go
type ParsedDate struct {
  Year  int `json:"year"`
  Month int `json:"month"` // 1-12, 0 if not present
  Day   int `json:"day"`   // 1-31, 0 if not present
}

func parseDate(dateStr string) ParsedDate {
  parts := strings.Split(dateStr, "/")
  // Parse and return structured date
}
```

**Frontend Receives**: Already-parsed dates in JSON response

### 7. Testing Strategy

**Backend Testing**:
- **Framework**: Go standard `testing` package
- **Test Files**: `*_test.go` files alongside implementation
- **Table-Driven Tests**: Go idiom for comprehensive test coverage
- **Coverage Goal**: 100% for new code (TDD)

**Frontend Testing**:
- **Framework**: Vitest (fast, modern)
- **Focus**: API client, chart rendering, UI interactions
- **Mocking**: Mock API responses for unit tests
- **Integration**: Test against local Go server

**Contract Testing**:
- Backend tests validate it produces correct response schemas
- Frontend tests validate it handles expected response schemas
- OpenAPI spec documents the contract

### 8. Static File Serving

**Decision**: Go server serves frontend static files

**Rationale**:
- Single server for API + frontend (simpler deployment)
- No CORS issues (same origin)
- Go's `http.FileServer` is efficient
- Frontend files bundled with backend binary

**Implementation**:
```go
// Serve frontend from /frontend directory
fs := http.FileServer(http.Dir("../frontend"))
http.Handle("/", fs)

// API routes prefixed with /api
http.HandleFunc("/api/years", handlers.GetYears)
// ...
```

**Deployment**: Single binary + frontend files = complete application

### 9. Error Handling Pattern

**Backend**:
- Return appropriate HTTP status codes
- JSON error responses: `{"error": "message"}`
- Log errors server-side
- Don't expose internal details to client

**Frontend**:
- Check response.ok before parsing JSON
- Display user-friendly error messages
- Fallback: "Unable to load reading data"
- Log technical details to console

### 10. Development Workflow

**Backend Development**:
```bash
cd backend
go test ./...        # Run all tests
go run main.go       # Start server on :8080
```

**Frontend Development**:
```bash
cd frontend
npm test             # Run Vitest tests
# Frontend served by Go backend at http://localhost:8080
```

**Integrated Testing**:
- Start Go server
- Open browser to http://localhost:8080
- Frontend calls Go API automatically
- Full stack testing in one environment

## Summary of Decisions

| Area | Decision | Justification |
|------|----------|---------------|
| Backend Framework | Go standard library net/http | Zero dependencies, sufficient for 3 endpoints |
| API Design | Simple REST with query params | Intuitive, follows YAGNI |
| Data Loading | Load on startup, cache in memory | Fast, simple, 600KB fits in RAM |
| CORS | Same-origin (Go serves frontend) | No CORS issues, single server |
| Frontend Chart | Chart.js | Lightweight, well-maintained, simple API |
| Date Parsing | Go backend | Better separation, cleaner frontend |
| Testing | Go `testing` + Vitest | Standard tools, minimal deps |
| Static Serving | Go http.FileServer | Single server, no CORS, simple deployment |
| Error Handling | JSON errors + HTTP codes | Standard REST pattern |
| Development | Integrated (Go serves all) | Single server, full-stack in one process |

## Architecture Summary

```
┌─────────────┐
│   Browser   │
│  (Frontend) │
│  HTML/CSS/JS│
└──────┬──────┘
       │ HTTP GET /
       │ HTTP GET /api/*
       ▼
┌─────────────────┐
│   Go Server     │
│   Port 8080     │
├─────────────────┤
│ • Serve static  │
│   files         │
│ • API endpoints │
│ • Load & filter │
│   books.json    │
└────────┬────────┘
         │ Read on startup
         ▼
    ┌─────────────┐
    │ books.json  │
    │ (repo root) │
    └─────────────┘
```

## No Unknowns Remain

All technical decisions have been made. Architecture is clear: Go API backend serving both the REST API and static frontend files. Ready for Phase 1 (data model, API contracts, and test scenarios).
