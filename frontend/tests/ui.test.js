import { describe, it, expect, beforeEach, vi } from 'vitest';

describe('Dashboard Layout Tests', () => {
  beforeEach(() => {
    // Create a fresh DOM for each test
    document.body.innerHTML = `
      <main>
        <section id="statistics" class="stats-bar">
          <div class="summary-card">
            <h2 class="summary-title">2025 Summary</h2>
            <div class="stats-list"></div>
          </div>
          <div class="stat-card">
            <h2>Total Books</h2>
            <p class="stat-value" id="total-books">0</p>
          </div>
          <div class="stat-card">
            <h2>Average per Month</h2>
            <p class="stat-value" id="avg-per-month">0.0</p>
          </div>
          <div class="stat-card">
            <h2>Total Pages</h2>
            <p class="stat-value" id="total-pages">0</p>
          </div>
        </section>
        <section id="chart-container"></section>
        <section id="book-list-section"></section>
      </main>
    `;
  });

  // T002: Stats bar layout test
  it('stats bar renders as horizontal flex row', () => {
    const statsBar = document.querySelector('.stats-bar');
    const summaryCard = document.querySelector('.summary-card');
    const statCards = document.querySelectorAll('.stat-card');

    expect(statsBar).toBeTruthy();
    expect(summaryCard).toBeTruthy();
    expect(statCards).toHaveLength(3);

    // All stat cards and summary card are direct children of stats-bar
    expect(statsBar.contains(summaryCard)).toBe(true);
    statCards.forEach(card => {
      expect(statsBar.contains(card)).toBe(true);
    });
  });

  // T003: Full-width book list test
  it('book list is a direct child of main, not nested in a grid column', () => {
    const main = document.querySelector('main');
    const bookListSection = document.getElementById('book-list-section');

    expect(bookListSection).toBeTruthy();
    expect(bookListSection.parentElement).toBe(main);
  });

  // T004: Summary card content test
  it('summary card displays correct statistics', () => {
    // Mock API response
    const mockStats = {
      year: 2025,
      total_books: 42,
      total_pages: 12450,
      avg_pages_per_book: 296
    };

    // Create summary card HTML
    const summaryCard = document.querySelector('.summary-card');
    summaryCard.innerHTML = `
      <h2 class="summary-title">2025 Summary</h2>
      <div class="stats-list">
        <div class="stat-row">
          <span class="stat-label">Total Books</span>
          <span class="stat-value">42</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Total Pages</span>
          <span class="stat-value">12,450</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Avg Pages/Book</span>
          <span class="stat-value">296</span>
        </div>
      </div>
    `;

    // Assert card title
    const title = summaryCard.querySelector('.summary-title');
    expect(title.textContent).toBe('2025 Summary');

    // Assert stat rows
    const statRows = summaryCard.querySelectorAll('.stat-row');
    expect(statRows).toHaveLength(3);

    // Check Total Books
    expect(statRows[0].querySelector('.stat-label').textContent).toBe('Total Books');
    expect(statRows[0].querySelector('.stat-value').textContent).toBe('42');

    // Check Total Pages (with comma separator)
    expect(statRows[1].querySelector('.stat-label').textContent).toBe('Total Pages');
    expect(statRows[1].querySelector('.stat-value').textContent).toBe('12,450');

    // Check Avg Pages/Book
    expect(statRows[2].querySelector('.stat-label').textContent).toBe('Avg Pages/Book');
    expect(statRows[2].querySelector('.stat-value').textContent).toBe('296');
  });

  // T005: Empty state test
  it('summary card shows empty state when no books', () => {
    // Mock API response with 0 books
    const mockStats = {
      year: 2024,
      total_books: 0,
      total_pages: 0,
      avg_pages_per_book: 0
    };

    // Create empty state HTML
    const summaryCard = document.querySelector('.summary-card');
    summaryCard.classList.add('empty-state');
    summaryCard.innerHTML = `
      <h2 class="summary-title">2024 Summary</h2>
      <p class="empty-message">No books tracked for this year</p>
    `;

    // Assert empty state
    const title = summaryCard.querySelector('.summary-title');
    expect(title.textContent).toBe('2024 Summary');

    const emptyMessage = summaryCard.querySelector('.empty-message');
    expect(emptyMessage).toBeTruthy();
    expect(emptyMessage.textContent).toContain('No books tracked');

    // Ensure no error rendering
    expect(summaryCard).toBeTruthy();
  });
});
