// Main application entry point
import { fetchYears, fetchStats, fetchBooks } from './api-client.js';
import { renderChart } from './chart.js';
import { populateYearSelector, updateStatistics, renderSummaryCard, showEmptyState, showError, showContent } from './ui.js';

let currentYear = new Date().getFullYear();

async function loadYearData(year) {
    currentYear = year;
    try {
        const stats = await fetchStats(year);
        
        renderSummaryCard(stats);
        
        if (stats.totalBooks === 0) {
            showEmptyState(year);
            hideBookList();
        } else {
            showContent();
            updateStatistics(stats);
            
            // Try to render chart, but don't fail if Chart.js isn't loaded
            try {
                const canvas = document.getElementById('monthly-chart');
                renderChart(canvas, stats.monthlyBreakdown);
            } catch (chartError) {
                console.warn('Chart rendering failed:', chartError);
                // Show a message instead of the chart if rendering fails
                const chartContainer = document.getElementById('chart-container');
                const canvas = document.getElementById('monthly-chart');
                if (chartContainer && canvas) {
                    canvas.style.display = 'none';
                    const message = document.createElement('p');
                    message.className = 'chart-error-message';
                    message.textContent = 'Chart visualization unavailable (Chart.js library failed to load)';
                    chartContainer.appendChild(message);
                }
            }
            
            // Load and display book list
            await loadBookList(year);
        }
        
    } catch (error) {
        console.error('Error loading year data:', error);
        showError('Unable to load reading data');
    }
}

async function loadBookList(year) {
    const bookList = document.getElementById('book-list');
    const bookCount = document.getElementById('book-count');
    const section = document.getElementById('book-list-section');
    
    if (!bookList || !section) return;
    
    section.style.display = 'block';
    bookList.innerHTML = '<p style="color: #7f8c8d; text-align: center; padding: 20px;">Loading...</p>';
    
    try {
        const data = await fetchBooks(year, { shelf: 'read' });
        const books = data.books || [];
        
        if (bookCount) {
            bookCount.textContent = books.length + ' book' + (books.length !== 1 ? 's' : '');
        }
        
        if (books.length === 0) {
            bookList.innerHTML = '<p style="color: #7f8c8d; text-align: center; padding: 20px;">No books read this year.</p>';
        } else {
            bookList.innerHTML = books.map(renderBookCard).join('');
        }
    } catch (error) {
        console.error('Error loading books:', error);
        bookList.innerHTML = '<p style="color: #e74c3c; text-align: center;">Failed to load books</p>';
    }
}

function hideBookList() {
    const section = document.getElementById('book-list-section');
    if (section) section.style.display = 'none';
}

function renderBookCard(book) {
    const dateStr = book.dateRead ? formatDate(book.dateRead) : '';
    const pagesStr = book.pages > 0 ? book.pages + ' pages' : '';
    
    let coverHtml = '<div class="book-cover-placeholder">📚</div>';
    if (book.coverUrl) {
        coverHtml = '<img src="' + escapeHtml(book.coverUrl) + '" alt="" class="book-cover" loading="lazy" onerror="this.style.display=\'none\'">';
    }
    
    let metaHtml = '';
    if (dateStr) metaHtml += '<span>📅 ' + dateStr + '</span>';
    if (pagesStr) metaHtml += '<span>📖 ' + pagesStr + '</span>';
    
    return '<div class="book-card">' +
        coverHtml +
        '<div class="book-info">' +
        '<div class="book-title">' + escapeHtml(book.title) + '</div>' +
        '<div class="book-author">' + escapeHtml(book.author) + '</div>' +
        (metaHtml ? '<div class="book-meta">' + metaHtml + '</div>' : '') +
        '</div></div>';
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const parts = dateStr.split('/');
    if (parts.length >= 2) {
        const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
        const month = parseInt(parts[1]) - 1;
        const day = parts[2] || '';
        return day ? months[month] + ' ' + parseInt(day) : months[month];
    }
    return dateStr;
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function init() {
    try {
        const yearsData = await fetchYears();
        
        if (!yearsData.years || yearsData.years.length === 0) {
            showError('No reading data available');
            return;
        }
        
        populateYearSelector(yearsData.years);
        
        const hasCurrentYear = yearsData.years.some(y => y.year === currentYear);
        const defaultYear = hasCurrentYear ? currentYear : yearsData.years[0].year;
        
        document.getElementById('year-selector').value = defaultYear;
        
        await loadYearData(defaultYear);
        
        document.getElementById('year-selector').addEventListener('change', (e) => {
            loadYearData(parseInt(e.target.value));
        });
        
    } catch (error) {
        console.error('Error initializing app:', error);
        showError('Unable to load reading data');
    }
}

document.addEventListener('DOMContentLoaded', init);
