# 📚 Reading Tracker

[![CI](https://github.com/kristenwomack/reading-app/actions/workflows/ci.yml/badge.svg)](https://github.com/kristenwomack/reading-app/actions/workflows/ci.yml)

A personal reading tracker with a modern dashboard and mobile-friendly admin for logging books on the go.

## Features

- 📊 **Dashboard** - Visualize reading progress with charts and statistics
- 📖 **Book List** - Browse books with covers, linked to Open Library
- 🎯 **Reading Goals** - Set yearly book targets with progress tracking
- ✏️ **Admin Panel** - Add books from any device with password protection
- 📷 **ISBN Scanner** - Scan barcodes to auto-fill book details
- 🔍 **Open Library Integration** - Fetch book info and covers automatically

## Screenshots

The dashboard shows your reading stats, a monthly chart, goal progress, and a list of books read.

## Quick Start

### Prerequisites
- Go 1.21 or later

### Run the App

```bash
cd backend
READING_APP_PASSWORD=yourpassword go run main.go
```

Open http://localhost:3000 in your browser.

### Add Books

1. Go to http://localhost:3000/admin
2. Enter your password
3. Add books manually or scan an ISBN barcode

## Project Structure

```
reading-app/
├── backend/                  # Go API server
│   ├── internal/
│   │   ├── auth/            # JWT authentication
│   │   ├── books/           # Book loading, filtering, stats
│   │   ├── handlers/        # HTTP request handlers
│   │   └── store/           # SQLite database layer
│   └── main.go              # Server entry point
├── frontend/                 # Vanilla JavaScript frontend
│   ├── src/
│   │   ├── admin.js         # Admin page functionality
│   │   ├── api-client.js    # API communication
│   │   ├── chart.js         # Chart.js integration
│   │   ├── main.js          # Dashboard logic
│   │   └── ui.js            # DOM manipulation
│   ├── styles/
│   │   ├── main.css         # Dashboard styles
│   │   └── admin.css        # Admin page styles
│   ├── index.html           # Dashboard
│   └── admin.html           # Admin panel
├── books.json                # Initial book data (imported on first run)
└── books.db                  # SQLite database (created automatically)
```

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go with standard library + SQLite |
| Frontend | Vanilla JavaScript, HTML, CSS |
| Charts | Chart.js |
| Database | SQLite (via modernc.org/sqlite) |
| Auth | JWT tokens with bcrypt |
| Book Data | Open Library API |

## API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/years` | List available years |
| GET | `/api/books?year=2025` | Get books for year |
| GET | `/api/books?year=2025&shelf=read` | Filter by shelf |
| GET | `/api/stats?year=2025` | Get statistics |
| GET | `/api/goals/:year` | Get reading goal |

### Protected (requires auth)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/books` | Add a book |
| PUT | `/api/books/:id` | Update a book |
| DELETE | `/api/books/:id` | Delete a book |
| POST | `/api/goals` | Set reading goal |
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/logout` | Logout |

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `READING_APP_PASSWORD` | Password for admin access | (required) |
| `DATABASE_PATH` | Path to SQLite database file | `../books.db` |
| `PORT` | Server port (set automatically by Railway) | `3000` |

## Development

### Run Backend
```bash
cd backend
go run main.go
```

### Run Tests
```bash
cd backend && go test ./...
cd frontend && npm test
```

### Build for Production
```bash
cd backend
go build -o reading-tracker main.go
READING_APP_PASSWORD=secret ./reading-tracker
```

## Deployment

### Railway (Recommended)

Deploy the full app to [Railway](https://railway.com) with automatic HTTPS:

1. **Create a Railway account** at [railway.com](https://railway.com)
2. **Create a new project** → "Deploy from GitHub Repo" → select `reading-app`
3. **Add a persistent volume** in the service settings:
   - Mount path: `/data`
   - This keeps your SQLite database safe across deploys
4. **Set environment variables** in the Railway dashboard:
   - `READING_APP_PASSWORD` — password for admin access (required)
   - `DATABASE_PATH` — set to `/data/books.db`
5. **Deploy** — Railway auto-builds from the Dockerfile and assigns a `.up.railway.app` URL

Railway auto-deploys on every push to the connected branch.

### Self-hosted

For self-hosted deployment, use a reverse proxy like Caddy for automatic HTTPS:

```
# Caddyfile
books.yourdomain.com {
    reverse_proxy localhost:3000
}
```

## Data Migration

On first run, the app automatically imports `books.json` into SQLite. After that, all data is stored in `books.db`.

## License

MIT License - see [LICENSE](LICENSE) for details.
