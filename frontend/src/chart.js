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
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Books Read',
                data: data,
                backgroundColor: 'rgba(54, 162, 235, 0.5)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
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
