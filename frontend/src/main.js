// Main application entry point
import { fetchYears, fetchStats } from './api-client.js';
import { renderChart } from './chart.js';
import { populateYearSelector, updateStatistics, renderSummaryCard, showEmptyState, showError, showContent } from './ui.js';

async function loadYearData(year) {
    try {
        const stats = await fetchStats(year);
        
        // T012: Always render summary card with stats
        renderSummaryCard(stats);
        
        if (stats.totalBooks === 0) {
            showEmptyState(year);
        } else {
            showContent();
            updateStatistics(stats);
            
            const canvas = document.getElementById('monthly-chart');
            renderChart(canvas, stats.monthlyBreakdown);
        }
    } catch (error) {
        console.error('Error loading year data:', error);
        showError('Unable to load reading data');
    }
}

async function init() {
    try {
        // Fetch available years
        const yearsData = await fetchYears();
        
        if (!yearsData.years || yearsData.years.length === 0) {
            showError('No reading data available');
            return;
        }
        
        // Populate year selector
        populateYearSelector(yearsData.years);
        
        // Set default year (most recent or current year)
        const currentYear = new Date().getFullYear();
        const hasCurrentYear = yearsData.years.some(y => y.year === currentYear);
        const defaultYear = hasCurrentYear ? currentYear : yearsData.years[0].year;
        
        document.getElementById('year-selector').value = defaultYear;
        
        // Load data for default year
        await loadYearData(defaultYear);
        
        // Set up year selector change event
        document.getElementById('year-selector').addEventListener('change', (e) => {
            loadYearData(parseInt(e.target.value));
        });
        
    } catch (error) {
        console.error('Error initializing app:', error);
        showError('Unable to load reading data');
    }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', init);
