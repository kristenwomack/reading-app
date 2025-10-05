# 📚 Reading Tracker

A personal reading progress tracker built with Go and vanilla JavaScript.

## Features

- 📊 **Interactive Dashboard** - Visualize your reading progress with Chart.js
- 📅 **Year Selector** - Browse reading data by year
- 📈 **Statistics** - Track total books, pages, and monthly breakdowns
- 🎨 **Clean UI** - Minimal, responsive design

## Architecture

### Monorepo Structure
```
reading-app/
├── backend/              # Go API server
│   ├── internal/
│   │   ├── books/       # Book loading, filtering, and stats
│   │   └── handlers/    # HTTP request handlers
│   └── main.go          # Server entry point with CORS
├── frontend/            # Vanilla JavaScript frontend
│   ├── src/
│   │   ├── api-client.js    # API communication
│   │   ├── chart.js         # Chart.js integration
│   │   ├── ui.js            # DOM manipulation
│   │   └── main.js          # Application entry point
│   ├── styles/
│   │   └── main.css         # Styling
│   └── index.html
├── shared/
│   └── contracts/       # OpenAPI specifications
└── books.json           # Reading data (repository root)
```

### Tech Stack

**Backend:**
- Go 1.21+ with standard library
- No external dependencies (stdlib only)
- RESTful API design

**Frontend:**
- Vanilla JavaScript (ES6+)
- Chart.js for visualizations
- No frameworks - pure HTML/CSS/JS

**Data:**
- JSON file storage (no database)
- No authentication (local single-user app)

### API Endpoints

```
GET /api/years              # List available years
GET /api/books?year=2025    # Get books for specific year
GET /api/stats?year=2025    # Get statistics for specific year
```

### Key Design Decisions

1. **CORS Support** - Middleware allows browser access from any origin
2. **Port 3000** - Runs on http://localhost:3000
3. **Minimal Dependencies** - Go stdlib + Chart.js only
4. **Static File Serving** - Backend serves frontend files
5. **TDD Approach** - Complete test coverage for all modules

## Quick Start

### Prerequisites
- Go 1.21 or later
- Modern web browser

### Running Locally

1. **Navigate to backend:**
   ```bash
   cd backend
   ```

2. **Run the server:**
   ```bash
   ./reading-tracker
   ```
   
   Or build first:
   ```bash
   go build -o reading-tracker main.go
   ./reading-tracker
   ```

3. **Open in browser:**
   ```
   http://localhost:3000
   ```

### Running Tests

```bash
cd backend
go test ./...
```

## Development

### Backend Development
```bash
cd backend
go test ./...              # Run tests
go build -o reading-tracker main.go   # Build binary
```

### Frontend Development
Frontend files are served by the Go backend. Simply refresh the browser to see changes.

## Project Structure Details

### Backend Modules

- **`internal/books/loader.go`** - Loads books from JSON file
- **`internal/books/filter.go`** - Filters books by year and other criteria
- **`internal/books/stats.go`** - Calculates reading statistics and monthly breakdowns
- **`internal/handlers/handlers.go`** - HTTP handlers for API endpoints

### Frontend Modules

- **`src/api-client.js`** - Fetches data from backend API
- **`src/chart.js`** - Creates Chart.js visualizations
- **`src/ui.js`** - Updates DOM with data
- **`src/main.js`** - Initializes app and coordinates modules

## Constitution

This project follows a constitution-based development approach:
- ✅ Monorepo structure
- ✅ Minimal dependencies
- ✅ Test-driven development
- ✅ Go + vanilla JavaScript
- ✅ No databases, no authentication

See `.specify/memory/constitution.md` for full guidelines.

## License

Personal project - all rights reserved.
