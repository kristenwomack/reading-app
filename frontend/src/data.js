// Client-side data layer for the static site.
// Loads books.json (and goals.json) directly in the browser and computes
// years, per-year book lists, and statistics — replacing the old Go backend.

const BOOKS_URL = 'books.json';
const GOALS_URL = 'goals.json';

const MONTH_NAMES = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

let booksPromise = null;
let goalsPromise = null;

// loadRawBooks fetches and caches the raw books.json array.
function loadRawBooks() {
    if (!booksPromise) {
        booksPromise = fetch(BOOKS_URL)
            .then((res) => {
                if (!res.ok) throw new Error(`Failed to load ${BOOKS_URL} (${res.status})`);
                return res.json();
            })
            .catch((err) => {
                booksPromise = null; // allow retry on failure
                throw err;
            });
    }
    return booksPromise;
}

// loadGoals fetches and caches goals.json. Missing file resolves to {}.
function loadGoals() {
    if (!goalsPromise) {
        goalsPromise = fetch(GOALS_URL)
            .then((res) => (res.ok ? res.json() : {}))
            .catch(() => ({}));
    }
    return goalsPromise;
}

// --- Field accessors (books.json uses Goodreads-style keys, mixed types) ---

function toStringValue(v) {
    if (v === null || v === undefined) return '';
    if (typeof v === 'number') return Number.isInteger(v) ? String(v) : String(v);
    return String(v);
}

export function getTitle(book) {
    return toStringValue(book['Title']).trim();
}

export function getPages(book) {
    const v = book['Number of Pages'];
    if (typeof v === 'number') return Math.trunc(v);
    if (typeof v === 'string') {
        const n = parseInt(v, 10);
        return Number.isNaN(n) ? 0 : n;
    }
    return 0;
}

export function getISBN(book) {
    const isbn13 = book['ISBN13'];
    if (typeof isbn13 === 'string' && isbn13 !== '') return isbn13;
    if (typeof isbn13 === 'number' && isbn13 > 0) return String(isbn13);
    const isbn = book['ISBN'];
    if (typeof isbn === 'string' && isbn !== '') return isbn;
    if (typeof isbn === 'number' && isbn > 0) return String(isbn);
    return '';
}

// coverURLByISBN mirrors the backend Open Library cover URL builder.
export function coverURLByISBN(isbn) {
    if (!isbn) return '';
    return `https://covers.openlibrary.org/b/isbn/${isbn}-M.jpg`;
}

function resolveCoverURL(book) {
    if (book['CoverURL']) return book['CoverURL'];
    return coverURLByISBN(getISBN(book));
}

// parseDate parses "YYYY/MM/DD" (month/day optional). Returns null if invalid.
// Mirrors backend books.ParseDate (year must be >= 1900).
export function parseDate(dateStr) {
    if (typeof dateStr !== 'string') return null;
    const trimmed = dateStr.trim();
    if (trimmed === '') return null;

    const parts = trimmed.split('/');
    const year = parseInt(parts[0], 10);
    if (Number.isNaN(year) || year < 1900) return null;

    const result = { year, month: 0, day: 0 };

    if (parts.length > 1 && parts[1] !== '') {
        const month = parseInt(parts[1], 10);
        if (Number.isNaN(month) || month < 1 || month > 12) return null;
        result.month = month;
    }
    if (parts.length > 2 && parts[2] !== '') {
        const day = parseInt(parts[2], 10);
        if (Number.isNaN(day) || day < 1 || day > 31) return null;
        result.day = day;
    }
    return result;
}

// --- Aggregations (mirror backend handlers/stats) ---

// computeYears returns [{ year, count }] for read books with a valid Date Read.
export function computeYears(books) {
    const counts = new Map();
    for (const book of books) {
        if (book['Shelf'] !== 'read') continue;
        const date = parseDate(book['Date Read']);
        if (!date) continue;
        counts.set(date.year, (counts.get(date.year) || 0) + 1);
    }
    return Array.from(counts, ([year, count]) => ({ year, count }));
}

function filterByYear(books, year) {
    return books.filter((book) => {
        const date = parseDate(book['Date Read']);
        return date && date.year === year;
    });
}

function filterByShelf(books, shelf) {
    return books.filter((book) => book['Shelf'] === shelf);
}

// computeBooksForYear returns view models for a year, optionally filtered by shelf.
export function computeBooksForYear(books, year, shelf) {
    let list = filterByYear(books, year);
    if (shelf) list = filterByShelf(list, shelf);
    return list.map((book) => {
        const date = parseDate(book['Date Read']);
        return {
            title: getTitle(book),
            author: book['Author'] || '',
            dateRead: book['Date Read'] || '',
            pages: getPages(book),
            month: date ? date.month : 0,
            shelf: book['Shelf'] || '',
            isbn: getISBN(book),
            coverUrl: resolveCoverURL(book),
        };
    });
}

// computeMonthlyBreakdown returns 12 entries with PascalCase keys (chart.js expects these).
export function computeMonthlyBreakdown(readBooks) {
    const breakdown = MONTH_NAMES.map((name, i) => ({ Month: i + 1, MonthName: name, Count: 0 }));
    for (const book of readBooks) {
        const date = parseDate(book['Date Read']);
        if (!date || date.month < 1 || date.month > 12) continue;
        breakdown[date.month - 1].Count++;
    }
    return breakdown;
}

// computeStats mirrors the backend /api/stats response for a year.
export function computeStats(books, year) {
    const readBooks = filterByShelf(filterByYear(books, year), 'read');
    let totalPages = 0;
    for (const book of readBooks) {
        const pages = getPages(book);
        if (pages > 0) totalPages += pages;
    }
    const totalBooks = readBooks.length;
    return {
        year,
        totalBooks,
        totalPages,
        averagePerMonth: totalBooks > 0 ? totalBooks / 12 : 0,
        monthlyBreakdown: computeMonthlyBreakdown(readBooks),
    };
}

// --- Public async API (consumed by api-client.js) ---

export async function getYears() {
    const books = await loadRawBooks();
    return { years: computeYears(books) };
}

export async function getBooks(year, options = {}) {
    const books = await loadRawBooks();
    return { books: computeBooksForYear(books, year, options.shelf) };
}

export async function getStats(year) {
    const books = await loadRawBooks();
    return computeStats(books, year);
}

export async function getGoal(year) {
    const goals = await loadGoals();
    const target = goals[String(year)];
    if (typeof target === 'number' && target > 0) {
        return { year, target };
    }
    return { year, target: null };
}
