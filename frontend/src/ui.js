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
