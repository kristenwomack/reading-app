# reading-app Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-01-16

## Active Technologies
- **Backend**: Go (latest stable) with standard library net/http (001-reading-tracker-webpage)
- **Frontend**: JavaScript ES6+ (vanilla JS) with HTML5/CSS3 + Chart.js (001-reading-tracker-webpage)
- JavaScript ES6+ (vanilla JS), CSS3 + Chart.js (existing), no new dependencies (002-style-the-table)
- N/A (existing books.json data source, no changes) (002-style-the-table)

## Project Structure
```
backend/                  # Go API service
├── internal/
│   ├── books/           # Book data loading and filtering
│   └── handlers/        # HTTP request handlers
└── main.go              # Server entry point

frontend/                # Static JavaScript site
├── src/
│   ├── api-client.js    # API communication
│   ├── chart.js         # Chart.js integration
│   └── ui.js            # DOM manipulation
└── index.html

shared/
└── contracts/           # API contracts (OpenAPI)

books.json              # Data source (repository root)
```

## Commands
**Backend**:
- `cd backend && go test ./...` - Run all Go tests
- `cd backend && go run main.go` - Start API server on :8080

**Frontend**:
- `cd frontend && npm test` - Run Vitest tests
- Frontend served by Go backend at http://localhost:8080

## Code Style
- **Go**: Standard Go conventions (gofmt, golint)
- **JavaScript**: ES6+ with vanilla JS, no frameworks

## Architecture
- Monorepo with backend + frontend workspaces
- Go API reads books.json, exposes REST endpoints
- Frontend consumes API, renders UI with Chart.js
- No database, no authentication (local single-user app)

## API Endpoints
- `GET /api/years` - List available years
- `GET /api/books?year=2025` - Get books for year
- `GET /api/stats?year=2025` - Get statistics for year

## Recent Changes
- 002-style-the-table: Added JavaScript ES6+ (vanilla JS), CSS3 + Chart.js (existing), no new dependencies
- 001-reading-tracker-webpage: Added Go API backend + vanilla JS frontend with Chart.js

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
