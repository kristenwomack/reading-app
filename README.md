# 📚 Reading Tracker

[![CI](https://github.com/kristenwomack/reading-app/actions/workflows/ci.yml/badge.svg)](https://github.com/kristenwomack/reading-app/actions/workflows/ci.yml)

A personal reading tracker — a static dashboard that visualizes your reading
progress. It reads a single `books.json` file directly in the browser, so there
is no backend, no database, and nothing to keep running. It deploys to GitHub
Pages automatically on every push to `main`.

## Features

- 📊 **Dashboard** — reading stats, a monthly chart, and yearly goal progress
- 📖 **Book List** — browse books with covers from Open Library
- 🎯 **Reading Goals** — per-year targets via `goals.json`
- 🗂️ **Single source of truth** — all data lives in `books.json` (git-tracked)

## How it works

`books.json` (and `goals.json`) at the repo root are the only data sources. The
frontend fetches them and computes everything — available years, per-year book
lists, statistics, and the monthly breakdown — entirely in the browser.

Adding a book = editing `books.json` (via a pull request). When the PR merges to
`main`, the GitHub Pages workflow rebuilds and the live site updates.

## Quick Start

```bash
cd frontend
npm install
npm run build          # assembles dist/ (frontend + books.json + goals.json)
npx serve dist         # or: python3 -m http.server --directory dist 8000
```

Open the served URL in your browser.

## Adding a book

Append an entry to `books.json` (Goodreads-style keys) and open a PR:

```json
{
  "Title": "Yesteryear",
  "Author": "Caro Claire Burke",
  "ISBN": "059380421X",
  "ISBN13": "9780593804216",
  "Publisher": "Alfred A. Knopf",
  "Number of Pages": 400,
  "Year Published": 2026,
  "Date Read": "2026/07/26",
  "Date Added": "2026/07/26",
  "Shelf": "read"
}
```

Only books with `"Shelf": "read"` and a valid `Date Read` (`YYYY/MM/DD`) are
counted in the dashboard. Covers are derived from the ISBN via Open Library.

## Setting goals

Edit `goals.json` to map a year to a target number of books:

```json
{
  "2025": 90,
  "2026": 90
}
```

## Project Structure

```
reading-app/
├── books.json                # Book data (source of truth)
├── goals.json                # Per-year reading goals
├── frontend/                 # Static site
│   ├── index.html
│   ├── build.js              # Assembles dist/ for preview and deploy
│   ├── src/
│   │   ├── main.js           # Entry point
│   │   ├── data.js           # Loads books.json, computes years/stats
│   │   ├── api-client.js     # Shim delegating to data.js
│   │   ├── chart.js          # D3 monthly chart
│   │   └── ui.js             # DOM updates
│   ├── styles/main.css
│   └── tests/
└── .github/workflows/
    ├── ci.yml                # Tests + build on PRs
    └── pages.yml             # Deploys to GitHub Pages on push to main
```

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Vanilla JavaScript (ES modules), HTML, CSS |
| Charts | D3.js |
| Data | Static `books.json` / `goals.json` |
| Book covers | Open Library |
| Hosting | GitHub Pages |
| Tests | Vitest |

## Development

```bash
cd frontend
npm test          # run unit tests
npm run build     # build the static site into dist/
```

## Deployment

The site deploys to **GitHub Pages** via `.github/workflows/pages.yml` on every
push to `main`. Enable it once under **Settings → Pages → Build and deployment →
Source: GitHub Actions**.

## License

MIT License - see [LICENSE](LICENSE) for details.
