# Reading Tracker Frontend

Static JavaScript frontend for the reading tracker application.

## Prerequisites

- Node.js 18 or later
- npm 9 or later
- Backend API running on http://localhost:8080

## Setup

```bash
# Install dependencies
npm install

# Run tests
npm test

# Run tests with UI
npm run test:ui

# Generate coverage report
npm run coverage
```

## How to Access

The frontend is served by the Go backend at http://localhost:8080

1. Start the backend server: `cd ../backend && go run main.go`
2. Open browser to http://localhost:8080
3. The frontend will automatically load and fetch data from the API

## Project Structure

```
frontend/
├── index.html           # Main HTML page
├── src/
│   ├── main.js         # Entry point and page initialization
│   ├── api-client.js   # API communication module
│   ├── chart.js        # Chart.js integration
│   └── ui.js           # DOM manipulation and UI updates
├── styles/
│   └── main.css        # CSS styling
└── tests/
    ├── api-client.test.js
    ├── chart.test.js
    ├── ui.test.js
    └── integration.test.js
```

## Features

- **Year Selector**: Choose different years to view reading progress
- **Statistics Display**: Total books, average per month, total pages
- **Monthly Chart**: Bar chart showing books read per month
- **Empty State**: Friendly message when no books for selected year
- **Error Handling**: Simple error display if data fails to load

## Browser Support

- Chrome (last 2 versions)
- Firefox (last 2 versions)
- Safari (last 2 versions)
- Edge (last 2 versions)

## Development

1. Tests are written first (TDD)
2. Use ES6 modules (type="module")
3. Vanilla JavaScript (no framework)
4. Chart.js for visualization only

## API Contract

The frontend expects these API endpoints:

- `GET /api/years` → `{years: [{year: 2025, count: 12}, ...]}`
- `GET /api/books?year=2025` → `{books: [{title, author, dateRead, pages, month}, ...]}`
- `GET /api/stats?year=2025` → `{year, totalBooks, totalPages, averagePerMonth, monthlyBreakdown}`

See `../shared/contracts/openapi.yaml` for complete specification.
