// Chart rendering with Chart.js
let chartInstance = null;

export function renderChart(canvasElement, monthlyData) {
    // Destroy existing chart if any
    if (chartInstance) {
        chartInstance.destroy();
    }

    const labels = monthlyData.map(m => m.MonthName);
    const data = monthlyData.map(m => m.Count);
    
    // Calculate cumulative progress toward goal of 90 books
    const cumulativeData = [];
    let runningTotal = 0;
    for (let i = 0; i < data.length; i++) {
        runningTotal += data[i];
        cumulativeData.push(runningTotal);
    }

    chartInstance = new Chart(canvasElement, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Books Read per Month',
                    data: data,
                    borderColor: 'rgb(75, 192, 192)',
                    backgroundColor: 'rgba(75, 192, 192, 0.2)',
                    borderWidth: 2,
                    tension: 0.1,
                    fill: true,
                    yAxisID: 'y'
                },
                {
                    label: 'Cumulative Progress (Goal: 90)',
                    data: cumulativeData,
                    borderColor: 'rgb(255, 99, 132)',
                    backgroundColor: 'rgba(255, 99, 132, 0.1)',
                    borderWidth: 2,
                    tension: 0.1,
                    fill: false,
                    yAxisID: 'y1'
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 1.618,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            scales: {
                y: {
                    type: 'linear',
                    display: true,
                    position: 'left',
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Books per Month'
                    },
                    ticks: {
                        stepSize: 1
                    }
                },
                y1: {
                    type: 'linear',
                    display: true,
                    position: 'right',
                    beginAtZero: true,
                    max: 90,
                    title: {
                        display: true,
                        text: 'Total Books (Goal: 90)'
                    },
                    grid: {
                        drawOnChartArea: false,
                    },
                    ticks: {
                        stepSize: 10
                    }
                }
            }
        }
    });
}
