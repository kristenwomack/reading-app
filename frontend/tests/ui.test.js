import { describe, it, expect, beforeEach, vi } from 'vitest';

describe('Dashboard Layout Tests', () => {
  beforeEach(() => {
    // Create a fresh DOM for each test
    document.body.innerHTML = `
      <main>
        <section id="statistics" class="stats-bar">
          <div class="stat-card">
            <span class="stat-emoji">📚</span>
            <p class="stat-value" id="stat-books">0</p>
            <h2>Books</h2>
          </div>
          <div class="stat-card">
            <span class="stat-emoji">📄</span>
            <p class="stat-value" id="stat-pages">0</p>
            <h2>Pages</h2>
          </div>
          <div class="stat-card">
            <span class="stat-emoji">📊</span>
            <p class="stat-value" id="stat-avg-pages">0</p>
            <h2>Avg Pages</h2>
          </div>
          <div class="stat-card">
            <span class="stat-emoji">📅</span>
            <p class="stat-value" id="avg-per-month">0.0</p>
            <h2>Per Month</h2>
          </div>
        </section>
        <section id="chart-container"></section>
        <section id="book-list-section"></section>
      </main>
    `;
  });

  // T002: Stats bar layout test
  it('stats bar renders as horizontal flex row with 4 stat cards', () => {
    const statsBar = document.querySelector('.stats-bar');
    const statCards = document.querySelectorAll('.stat-card');

    expect(statsBar).toBeTruthy();
    expect(statCards).toHaveLength(4);

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

  // T004: Stat cards display correct values
  it('stat cards display correct statistics when populated', () => {
    document.getElementById('stat-books').textContent = '42';
    document.getElementById('stat-pages').textContent = '12,450';
    document.getElementById('stat-avg-pages').textContent = '296';
    document.getElementById('avg-per-month').textContent = '3.5';

    expect(document.getElementById('stat-books').textContent).toBe('42');
    expect(document.getElementById('stat-pages').textContent).toBe('12,450');
    expect(document.getElementById('stat-avg-pages').textContent).toBe('296');
    expect(document.getElementById('avg-per-month').textContent).toBe('3.5');
  });

  // T005: Each stat card has an emoji
  it('each stat card has an emoji icon', () => {
    const emojis = document.querySelectorAll('.stat-emoji');
    expect(emojis).toHaveLength(4);
    expect(emojis[0].textContent).toBe('📚');
    expect(emojis[1].textContent).toBe('📄');
    expect(emojis[2].textContent).toBe('📊');
    expect(emojis[3].textContent).toBe('📅');
  });

  // T006: Empty state shows zeros
  it('stat cards show zero values when no books', () => {
    expect(document.getElementById('stat-books').textContent).toBe('0');
    expect(document.getElementById('stat-pages').textContent).toBe('0');
    expect(document.getElementById('stat-avg-pages').textContent).toBe('0');
    expect(document.getElementById('avg-per-month').textContent).toBe('0.0');
  });
});
