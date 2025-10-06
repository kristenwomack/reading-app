// Chart rendering with Chart.js
let chartInstance = null;

export function renderChart(canvasElement, monthlyData) {
    // Destroy existing chart if any
    if (chartInstance) {
        chartInstance.destroy();
    }

    const labels = monthlyData.map(m => m.MonthName);
    const data = monthlyData.map(m => m.Count);

    chartInstance = new Chart(canvasElement, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Books Read',
                data: data,
                borderColor: 'rgb(75, 192, 192)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                borderWidth: 2,
                tension: 0.1,
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 1.618,
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1
                    }
                }
            }
        }
    });
}
