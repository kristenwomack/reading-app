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

// T010: Render summary card with reading statistics
export function renderSummaryCard(stats) {
    const summaryCard = document.querySelector('.summary-card');
    if (!summaryCard) return;

    // Update title
    const summaryTitle = summaryCard.querySelector('.summary-title');
    if (summaryTitle) {
        summaryTitle.textContent = `${stats.year} Summary`;
    }

    // Handle empty state (0 books)
    if (stats.totalBooks === 0) {
        summaryCard.classList.add('empty-state');
        summaryCard.querySelector('.stats-list').innerHTML = `
            <p class="empty-message">No books tracked for this year</p>
        `;
        return;
    }

    // Remove empty state class if it exists
    summaryCard.classList.remove('empty-state');

    // Render stat rows
    const statsList = summaryCard.querySelector('.stats-list');
    if (statsList) {
        statsList.innerHTML = `
            <div class="stat-row">
                <span class="stat-label">Total Books</span>
                <span class="stat-value">${stats.totalBooks}</span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Total Pages</span>
                <span class="stat-value">${formatNumber(stats.totalPages)}</span>
            </div>
            <div class="stat-row">
                <span class="stat-label">Avg Pages/Book</span>
                <span class="stat-value">${stats.averagePagesPerBook || Math.floor(stats.totalPages / stats.totalBooks)}</span>
            </div>
        `;
    }
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
    document.getElementById('statistics').style.display = 'grid';
    document.getElementById('chart-container').style.display = 'block';
    document.getElementById('empty-state').classList.add('hidden');
    document.getElementById('error-state').classList.add('hidden');
}
