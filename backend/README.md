# Reading Tracker Backend

Go API backend for the reading tracker application.

## Prerequisites

- Go 1.21 or later
- Access to `books.json` in repository root

## Setup

```bash
# Initialize module (if not already done)
go mod init github.com/kristenwomack/reading-app/backend

# Run tests
go test ./...

# Run server
go run main.go
```

## API Endpoints

The server runs on `http://localhost:8080` and provides:

- `GET /api/years` - Returns available years with book counts
- `GET /api/books?year=YYYY` - Returns books for specified year
- `GET /api/stats?year=YYYY` - Returns statistics for specified year
- `GET /` - Serves frontend static files

## Project Structure

```
backend/
├── main.go                 # HTTP server entry point
├── internal/
│   ├── books/             # Book data loading and processing
│   │   ├── loader.go      # Load and parse books.json
│   │   ├── filter.go      # Filter by year and shelf
│   │   └── stats.go       # Calculate statistics
│   └── handlers/          # HTTP request handlers
│       ├── years.go       # GET /api/years
│       ├── books.go       # GET /api/books
│       └── stats.go       # GET /api/stats
└── testdata/
    └── books_test.json    # Test fixtures
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/books/
go test ./internal/handlers/
```

## Development

1. Tests are written first (TDD)
2. All tests must pass before committing
3. Use table-driven tests for comprehensive coverage
4. Follow standard Go conventions (gofmt, golint)

## API Contract

See `../shared/contracts/openapi.yaml` for complete API specification.
