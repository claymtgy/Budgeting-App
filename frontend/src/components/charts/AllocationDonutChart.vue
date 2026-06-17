<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import * as d3 from 'd3'
import { formatCents } from '@/utils/money'

const props = defineProps({
  envelopes: { type: Array, default: () => [] },
  unallocatedCents: { type: Number, default: 0 },
  valueKey: { type: String, default: 'allocated_cents' }
})

const chartRef = ref(null)
const colors = d3.schemeTableau10

const chartData = computed(() => {
  const extraLabel = props.valueKey === 'spent_cents' ? 'Unspent' : 'Unallocated'
  return [
    ...props.envelopes.map((e) => ({
      name: e.name,
      value: e[props.valueKey] ?? 0
    })),
    ...(props.unallocatedCents > 0
      ? [{ name: extraLabel, value: props.unallocatedCents }]
      : [])
  ].filter((d) => d.value > 0)
})

const legendItems = computed(() =>
  chartData.value.map((item, index) => ({
    ...item,
    color: colors[index % colors.length]
  }))
)

function render() {
  if (!chartRef.value) return

  d3.select(chartRef.value).selectAll('*').remove()

  if (!chartData.value.length) {
    chartRef.value.innerHTML = '<p class="empty-chart">No allocation data yet.</p>'
    return
  }

  const width = 280
  const height = 280
  const radius = Math.min(width, height) / 2

  const svg = d3
    .select(chartRef.value)
    .append('svg')
    .attr('viewBox', `0 0 ${width} ${height}`)
    .append('g')
    .attr('transform', `translate(${width / 2},${height / 2})`)

  const pie = d3.pie().value((d) => d.value)
  const arc = d3.arc().innerRadius(radius * 0.55).outerRadius(radius * 0.9)

  svg
    .selectAll('path')
    .data(pie(chartData.value))
    .enter()
    .append('path')
    .attr('d', arc)
    .attr('fill', (_, i) => colors[i % colors.length])
}

onMounted(render)
onUnmounted(() => {
  if (chartRef.value) d3.select(chartRef.value).selectAll('*').remove()
})

watch(() => [props.envelopes, props.unallocatedCents, props.valueKey], render, { deep: true })
</script>

<template>
  <div class="donut-chart">
    <div ref="chartRef" class="chart-container"></div>
    <ul v-if="legendItems.length" class="legend" aria-label="Chart legend">
      <li v-for="item in legendItems" :key="item.name" class="legend-item">
        <span class="legend-swatch" :style="{ backgroundColor: item.color }" aria-hidden="true"></span>
        <span class="legend-label">{{ item.name }}</span>
        <span class="legend-value">{{ formatCents(item.value) }}</span>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.donut-chart {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.chart-container {
  width: 100%;
  max-width: 280px;
  min-height: 280px;
}

.chart-container :deep(.empty-chart) {
  margin: 0;
  color: #868e96;
  text-align: center;
  padding: 2rem 0;
}

.legend {
  list-style: none;
  margin: 0;
  padding: 0;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  font-size: 0.9rem;
}

.legend-swatch {
  width: 0.85rem;
  height: 0.85rem;
  border-radius: 3px;
  flex-shrink: 0;
}

.legend-label {
  flex: 1;
  min-width: 0;
  color: #1a1a2e;
}

.legend-value {
  font-weight: 600;
  color: #495057;
  flex-shrink: 0;
}

@media (min-width: 480px) {
  .legend {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.5rem 1rem;
  }
}
</style>
