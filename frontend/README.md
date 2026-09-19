# Reading Tracker Frontend

Static, dependency-free JavaScript frontend for the reading tracker. It reads
`books.json` and `goals.json` directly in the browser — there is no backend.

## Prerequisites

- Node.js 18 or later (only for tests and the build script)

## Setup

```bash
npm install       # install dev dependencies (vitest)
npm test          # run unit tests once
npm run test:watch
npm run coverage
```

## Run locally

Build the site (copies frontend assets + `books.json` + `goals.json` into `dist/`)
and serve it with any static file server:

```bash
npm run build
npx serve dist          # or: python3 -m http.server --directory dist 8000
```

Then open the served URL. Everything is computed client-side from `books.json`.

## Project Structure

```
frontend/
├── index.html          # Dashboard page
├── build.js            # Assembles dist/ for local preview and Pages deploy
├── src/
│   ├── main.js         # Entry point and page initialization
│   ├── data.js         # Loads books.json/goals.json, computes years & stats
│   ├── api-client.js   # Thin shim delegating to data.js
│   ├── chart.js        # D3 monthly chart
│   └── ui.js           # DOM manipulation and UI updates
├── styles/
│   └── main.css        # CSS styling
└── tests/
    ├── data.test.js    # Data-layer unit tests
    └── ui.test.js      # Dashboard DOM tests
```

## Data model

All data comes from two git-tracked files at the repo root:

- `books.json` — one object per book (Goodreads-style keys: `Title`, `Author`,
  `Number of Pages`, `Date Read`, `Shelf`, `ISBN`, `ISBN13`, ...).
- `goals.json` — `{ "2025": 90, "2026": 90 }` mapping a year to a book target.

To add a book or change a goal, edit these files (via a PR). Merging to `main`
triggers the GitHub Pages deploy and the live site updates automatically.
