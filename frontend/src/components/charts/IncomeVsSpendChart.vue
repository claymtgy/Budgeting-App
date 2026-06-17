<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import * as d3 from 'd3'

const props = defineProps({
  trends: { type: Array, default: () => [] }
})

const chartRef = ref(null)

function monthLabel(monthValue) {
  const [year, month] = monthValue.split('-').map(Number)
  const date = new Date(year, month - 1, 1)
  return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' })
}

function render() {
  if (!chartRef.value) return
  d3.select(chartRef.value).selectAll('*').remove()

  if (!props.trends.length) {
    chartRef.value.innerHTML = '<p>No income or spending data yet.</p>'
    return
  }

  const data = props.trends
  const margin = { top: 20, right: 16, bottom: 52, left: 52 }
  const width = Math.max(360, data.length * 56)
  const height = 300
  const innerWidth = width - margin.left - margin.right
  const innerHeight = height - margin.top - margin.bottom

  const svg = d3
    .select(chartRef.value)
    .append('svg')
    .attr('viewBox', `0 0 ${width} ${height}`)
    .append('g')
    .attr('transform', `translate(${margin.left},${margin.top})`)

  const x0 = d3
    .scaleBand()
    .domain(data.map((d) => d.month))
    .range([0, innerWidth])
    .paddingInner(0.25)

  const x1 = d3.scaleBand().domain(['income', 'spent']).range([0, x0.bandwidth()]).padding(0.08)

  const maxValue = d3.max(data, (d) => Math.max(d.income_cents, d.spent_cents)) || 0
  const y = d3.scaleLinear().domain([0, maxValue]).nice().range([innerHeight, 0])

  svg
    .append('g')
    .attr('transform', `translate(0,${innerHeight})`)
    .call(d3.axisBottom(x0).tickFormat(monthLabel))
    .selectAll('text')
    .attr('transform', 'rotate(-35)')
    .style('text-anchor', 'end')

  svg.append('g').call(d3.axisLeft(y).ticks(5).tickFormat((d) => `$${d / 100}`))

  const legend = svg.append('g').attr('transform', `translate(0,${-8})`)
  ;[
    { label: 'Income', color: '#2b8a3e' },
    { label: 'Spent', color: '#e03131' }
  ].forEach((item, i) => {
    const g = legend.append('g').attr('transform', `translate(${i * 88}, 0)`)
    g.append('rect').attr('width', 12).attr('height', 12).attr('fill', item.color).attr('rx', 2)
    g.append('text').attr('x', 16).attr('y', 10).attr('font-size', '11px').text(item.label)
  })

  const groups = svg
    .selectAll('g.month-group')
    .data(data)
    .enter()
    .append('g')
    .attr('transform', (d) => `translate(${x0(d.month)},0)`)

  groups
    .selectAll('rect')
    .data((d) => [
      { key: 'income', value: d.income_cents },
      { key: 'spent', value: d.spent_cents }
    ])
    .enter()
    .append('rect')
    .attr('x', (d) => x1(d.key))
    .attr('y', (d) => y(d.value))
    .attr('width', x1.bandwidth())
    .attr('height', (d) => innerHeight - y(d.value))
    .attr('fill', (d) => (d.key === 'income' ? '#2b8a3e' : '#e03131'))
    .attr('rx', 2)
}

onMounted(render)
onUnmounted(() => {
  if (chartRef.value) d3.select(chartRef.value).selectAll('*').remove()
})

watch(() => props.trends, render, { deep: true })
</script>

<template>
  <div ref="chartRef" class="chart-container"></div>
</template>
