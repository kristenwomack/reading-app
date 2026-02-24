// Chart rendering with D3.js
const ANNUAL_GOAL = 90;

let chartSvg = null;

export function renderChart(container, monthlyData) {
    if (!container || !monthlyData) return;

    // Clear previous chart
    container.innerHTML = '';

    const labels = monthlyData.map(m => m.MonthName);
    const data = monthlyData.map(m => m.Count);

    // Calculate cumulative progress
    const cumulativeData = [];
    let runningTotal = 0;
    for (let i = 0; i < data.length; i++) {
        runningTotal += data[i];
        cumulativeData.push(runningTotal);
    }

    const margin = { top: 20, right: 60, bottom: 40, left: 50 };
    const containerWidth = container.clientWidth || 600;
    const width = containerWidth - margin.left - margin.right;
    const height = Math.round(containerWidth / 1.618) - margin.top - margin.bottom;

    const svg = d3.select(container)
        .append('svg')
        .attr('viewBox', `0 0 ${containerWidth} ${height + margin.top + margin.bottom}`)
        .attr('preserveAspectRatio', 'xMidYMid meet')
        .style('width', '100%')
        .style('height', 'auto');

    const g = svg.append('g')
        .attr('transform', `translate(${margin.left},${margin.top})`);

    // Scales
    const x = d3.scaleBand()
        .domain(labels)
        .range([0, width])
        .padding(0.1);

    const yLeft = d3.scaleLinear()
        .domain([0, Math.max(...data, 1)])
        .nice()
        .range([height, 0]);

    const yRight = d3.scaleLinear()
        .domain([0, ANNUAL_GOAL])
        .range([height, 0]);

    // Axes
    g.append('g')
        .attr('transform', `translate(0,${height})`)
        .call(d3.axisBottom(x))
        .selectAll('text')
        .style('font-size', '11px');

    g.append('g')
        .call(d3.axisLeft(yLeft).ticks(Math.max(...data, 1)).tickFormat(d3.format('d')))
        .append('text')
        .attr('fill', 'rgb(75, 192, 192)')
        .attr('transform', 'rotate(-90)')
        .attr('y', -40)
        .attr('x', -height / 2)
        .attr('text-anchor', 'middle')
        .style('font-size', '12px')
        .text('Books per Month');

    g.append('g')
        .attr('transform', `translate(${width},0)`)
        .call(d3.axisRight(yRight).ticks(9).tickFormat(d3.format('d')))
        .append('text')
        .attr('fill', 'rgb(255, 99, 132)')
        .attr('transform', 'rotate(-90)')
        .attr('y', 50)
        .attr('x', -height / 2)
        .attr('text-anchor', 'middle')
        .style('font-size', '12px')
        .text(`Total Books (Goal: ${ANNUAL_GOAL})`);

    // Line generator using band midpoints
    const xMid = (label) => x(label) + x.bandwidth() / 2;

    // Monthly area + line (left axis)
    const areaGen = d3.area()
        .x((_, i) => xMid(labels[i]))
        .y0(height)
        .y1((d) => yLeft(d))
        .curve(d3.curveMonotoneX);

    const lineGen = d3.line()
        .x((_, i) => xMid(labels[i]))
        .y((d) => yLeft(d))
        .curve(d3.curveMonotoneX);

    g.append('path')
        .datum(data)
        .attr('fill', 'rgba(75, 192, 192, 0.2)')
        .attr('d', areaGen);

    g.append('path')
        .datum(data)
        .attr('fill', 'none')
        .attr('stroke', 'rgb(75, 192, 192)')
        .attr('stroke-width', 2)
        .attr('d', lineGen);

    // Cumulative line (right axis)
    const cumLineGen = d3.line()
        .x((_, i) => xMid(labels[i]))
        .y((d) => yRight(d))
        .curve(d3.curveMonotoneX);

    g.append('path')
        .datum(cumulativeData)
        .attr('fill', 'none')
        .attr('stroke', 'rgb(255, 99, 132)')
        .attr('stroke-width', 2)
        .attr('d', cumLineGen);

    // Data point dots
    g.selectAll('.dot-monthly')
        .data(data)
        .enter().append('circle')
        .attr('cx', (_, i) => xMid(labels[i]))
        .attr('cy', d => yLeft(d))
        .attr('r', 3)
        .attr('fill', 'rgb(75, 192, 192)');

    g.selectAll('.dot-cumulative')
        .data(cumulativeData)
        .enter().append('circle')
        .attr('cx', (_, i) => xMid(labels[i]))
        .attr('cy', d => yRight(d))
        .attr('r', 3)
        .attr('fill', 'rgb(255, 99, 132)');

    // Legend
    const legend = svg.append('g')
        .attr('transform', `translate(${margin.left + 10}, ${margin.top})`);

    legend.append('rect').attr('width', 12).attr('height', 12).attr('fill', 'rgb(75, 192, 192)');
    legend.append('text').attr('x', 16).attr('y', 10).text('Books Read per Month').style('font-size', '11px');

    legend.append('rect').attr('x', 160).attr('width', 12).attr('height', 12).attr('fill', 'rgb(255, 99, 132)');
    legend.append('text').attr('x', 176).attr('y', 10).text(`Cumulative Progress (Goal: ${ANNUAL_GOAL})`).style('font-size', '11px');

    chartSvg = svg;
}
