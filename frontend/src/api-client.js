// API client for fetching data from backend
const API_BASE = '/api';

export async function fetchYears() {
    const response = await fetch(`${API_BASE}/years`);
    if (!response.ok) {
        throw new Error('Failed to fetch years');
    }
    return response.json();
}

export async function fetchBooks(year) {
    const response = await fetch(`${API_BASE}/books?year=${year}`);
    if (!response.ok) {
        throw new Error('Failed to fetch books');
    }
    return response.json();
}

export async function fetchStats(year) {
    const response = await fetch(`${API_BASE}/stats?year=${year}`);
    if (!response.ok) {
        throw new Error('Failed to fetch stats');
    }
    return response.json();
}
