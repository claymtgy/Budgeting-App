<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import * as d3 from 'd3'

const props = defineProps({
  envelopes: { type: Array, default: () => [] },
  unallocatedCents: { type: Number, default: 0 }
})

const chartRef = ref(null)
let svg

const colors = d3.schemeTableau10

function render() {
  if (!chartRef.value) return

  d3.select(chartRef.value).selectAll('*').remove()

  const data = [
    ...props.envelopes.map((e) => ({
      name: e.name,
      value: e.allocated_cents
    })),
    ...(props.unallocatedCents > 0
      ? [{ name: 'Unallocated', value: props.unallocatedCents }]
      : [])
  ].filter((d) => d.value > 0)

  if (!data.length) {
    chartRef.value.innerHTML = '<p>No allocation data yet.</p>'
    return
  }

  const width = 320
  const height = 320
  const radius = Math.min(width, height) / 2

  svg = d3
    .select(chartRef.value)
    .append('svg')
    .attr('viewBox', `0 0 ${width} ${height}`)
    .append('g')
    .attr('transform', `translate(${width / 2},${height / 2})`)

  const pie = d3.pie().value((d) => d.value)
  const arc = d3.arc().innerRadius(radius * 0.55).outerRadius(radius * 0.9)

  const arcs = svg.selectAll('arc').data(pie(data)).enter().append('g')

  arcs
    .append('path')
    .attr('d', arc)
    .attr('fill', (_, i) => colors[i % colors.length])

  arcs
    .append('text')
    .attr('transform', (d) => `translate(${arc.centroid(d)})`)
    .attr('text-anchor', 'middle')
    .attr('font-size', '11px')
    .text((d) => d.data.name)
}

onMounted(render)
onUnmounted(() => {
  if (chartRef.value) d3.select(chartRef.value).selectAll('*').remove()
})

watch(() => [props.envelopes, props.unallocatedCents], render, { deep: true })
</script>

<template>
  <div ref="chartRef" class="chart-container"></div>
</template>
