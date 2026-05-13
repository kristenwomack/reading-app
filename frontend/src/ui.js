// UI manipulation functions

export function populateYearSelector(years) {
    const selector = document.getElementById('year-selector');
    selector.innerHTML = '';
    
    // Sort years descending
    years.sort((a, b) => b.year - a.year);
    
    years.forEach(yearInfo => {
        const option = document.createElement('option');
        option.value = yearInfo.year;
        option.textContent = `${yearInfo.year} (${yearInfo.count} books)`;
        selector.appendChild(option);
    });
}

export function updateStatistics(stats) {
    document.getElementById('total-books').textContent = stats.totalBooks;
    document.getElementById('avg-per-month').textContent = stats.averagePerMonth.toFixed(1);
    document.getElementById('total-pages').textContent = stats.totalPages;
    document.getElementById('current-year').textContent = stats.year;
}

// T011: Number formatting function
export function formatNumber(num) {
    return num.toLocaleString('en-US');
}

// T010: Render individual stat cards with reading statistics
export function renderStatCards(stats) {
    const booksEl = document.getElementById('stat-books');
    const pagesEl = document.getElementById('stat-pages');
    const avgPagesEl = document.getElementById('stat-avg-pages');

    if (stats.totalBooks === 0) {
        if (booksEl) booksEl.textContent = '0';
        if (pagesEl) pagesEl.textContent = '0';
        if (avgPagesEl) avgPagesEl.textContent = '0';
        return;
    }

    if (booksEl) booksEl.textContent = stats.totalBooks;
    if (pagesEl) pagesEl.textContent = formatNumber(stats.totalPages);
    if (avgPagesEl) avgPagesEl.textContent = stats.averagePagesPerBook || Math.floor(stats.totalPages / stats.totalBooks);
}


export function showEmptyState(year) {
    document.getElementById('statistics').style.display = 'none';
    document.getElementById('chart-container').style.display = 'none';
    document.getElementById('empty-state').classList.remove('hidden');
    document.getElementById('empty-year').textContent = year;
    document.getElementById('error-state').classList.add('hidden');
}

export function showError(message) {
    document.getElementById('statistics').style.display = 'none';
    document.getElementById('chart-container').style.display = 'none';
    document.getElementById('empty-state').classList.add('hidden');
    document.getElementById('error-state').classList.remove('hidden');
    document.getElementById('error-message').textContent = message;
}

export function showContent() {
    document.getElementById('statistics').style.display = 'flex';
    document.getElementById('chart-container').style.display = 'block';
    document.getElementById('empty-state').classList.add('hidden');
    document.getElementById('error-state').classList.add('hidden');
}
