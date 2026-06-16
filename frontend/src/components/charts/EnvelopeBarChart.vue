<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import * as d3 from 'd3'

const props = defineProps({
  envelopes: { type: Array, default: () => [] }
})

const chartRef = ref(null)

function render() {
  if (!chartRef.value) return
  d3.select(chartRef.value).selectAll('*').remove()

  if (!props.envelopes.length) {
    chartRef.value.innerHTML = '<p>No envelope data yet.</p>'
    return
  }

  const data = props.envelopes
  const margin = { top: 20, right: 20, bottom: 40, left: 50 }
  const width = 460
  const height = 320
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
    .domain(data.map((d) => d.name))
    .range([0, innerWidth])
    .paddingInner(0.2)

  const x1 = d3.scaleBand().domain(['allocated', 'spent']).range([0, x0.bandwidth()]).padding(0.05)

  const maxValue = d3.max(data, (d) => Math.max(d.allocated_cents, d.spent_cents)) || 0
  const y = d3.scaleLinear().domain([0, maxValue]).nice().range([innerHeight, 0])

  svg
    .append('g')
    .attr('transform', `translate(0,${innerHeight})`)
    .call(d3.axisBottom(x0))
    .selectAll('text')
    .attr('transform', 'rotate(-20)')
    .style('text-anchor', 'end')

  svg.append('g').call(d3.axisLeft(y).ticks(5).tickFormat((d) => `$${d / 100}`))

  const groups = svg
    .selectAll('g.bar-group')
    .data(data)
    .enter()
    .append('g')
    .attr('transform', (d) => `translate(${x0(d.name)},0)`)

  groups
    .selectAll('rect')
    .data((d) => [
      { key: 'allocated', value: d.allocated_cents },
      { key: 'spent', value: d.spent_cents }
    ])
    .enter()
    .append('rect')
    .attr('x', (d) => x1(d.key))
    .attr('y', (d) => y(d.value))
    .attr('width', x1.bandwidth())
    .attr('height', (d) => innerHeight - y(d.value))
    .attr('fill', (d) => (d.key === 'allocated' ? '#3b5bdb' : '#e03131'))
}

onMounted(render)
onUnmounted(() => {
  if (chartRef.value) d3.select(chartRef.value).selectAll('*').remove()
})

watch(() => props.envelopes, render, { deep: true })
</script>

<template>
  <div ref="chartRef" class="chart-container"></div>
</template>
