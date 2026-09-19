// API client shim for the static site.
// Preserves the original fetch* function signatures but sources all data from
// books.json / goals.json via the client-side data layer (no backend).
import { getYears, getBooks, getStats, getGoal } from './data.js';

export async function fetchYears() {
    return getYears();
}

export async function fetchBooks(year, options = {}) {
    return getBooks(year, options);
}

export async function fetchStats(year) {
    return getStats(year);
}

export async function fetchGoal(year) {
    return getGoal(year);
}
