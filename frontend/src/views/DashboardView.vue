<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import api from '@/api/client'
import {
  currentMonthValue,
  formatCents,
  formatMonthLabel,
  monthOptions
} from '@/utils/money'
import AllocationDonutChart from '@/components/charts/AllocationDonutChart.vue'
import EnvelopeBarChart from '@/components/charts/EnvelopeBarChart.vue'
import IncomeVsSpendChart from '@/components/charts/IncomeVsSpendChart.vue'

const monthly = ref(null)
const budgetSummary = ref(null)
const trends = ref([])
const error = ref('')
const loading = ref(true)
const selectedMonth = ref(currentMonthValue())
const activeTab = ref('overview')

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'trends', label: 'Trends' },
  { id: 'budget', label: 'Budget' },
  { id: 'spending', label: 'Spending' }
]

const months = monthOptions(24)

const isCurrentMonth = computed(() => selectedMonth.value === currentMonthValue())
const monthLabel = computed(() => formatMonthLabel(selectedMonth.value))
const showMonthPicker = computed(() => activeTab.value !== 'trends')

const remainingLabel = computed(() =>
  isCurrentMonth.value ? 'Balance' : 'Left in month'
)

const monthRemainingCents = computed(() => {
  if (!monthly.value) return 0
  return monthly.value.envelopes.reduce((sum, envelope) => {
    if (isCurrentMonth.value) {
      return sum + envelope.available_cents
    }
    return sum + (envelope.allocated_cents - envelope.spent_cents)
  }, 0)
})

async function loadTrends() {
  try {
    const { data } = await api.get('/api/summary/trends', { params: { months: 12 } })
    trends.value = data
  } catch {
    trends.value = []
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const requests = [
      api.get('/api/summary/monthly', { params: { month: selectedMonth.value } })
    ]
    if (isCurrentMonth.value) {
      requests.push(api.get('/api/summary'))
    }

    const [monthlyRes, budgetRes] = await Promise.all(requests)
    monthly.value = monthlyRes.data
    budgetSummary.value = budgetRes?.data ?? null
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load dashboard'
  } finally {
    loading.value = false
  }
}

function availableClass(envelope) {
  const cents = isCurrentMonth.value
    ? envelope.available_cents
    : envelope.allocated_cents - envelope.spent_cents
  if (cents < 0) return 'over-budget'
  if (cents === 0) return 'empty'
  return ''
}

function remainingCents(envelope) {
  if (isCurrentMonth.value) return envelope.available_cents
  return envelope.allocated_cents - envelope.spent_cents
}

watch(selectedMonth, load)

onMounted(() => {
  loadTrends()
  load()
})
</script>

<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2 class="page-title">Dashboard</h2>
      <div v-if="showMonthPicker" class="month-picker">
        <label for="review-month">Month</label>
        <select id="review-month" v-model="selectedMonth">
          <option v-for="option in months" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
      </div>
    </div>

    <nav class="dashboard-tabs" aria-label="Dashboard sections">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        class="dashboard-tab"
        :class="{ active: activeTab === tab.id }"
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </nav>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="loading && activeTab !== 'trends'">Loading...</p>

    <!-- Overview -->
    <template v-else-if="activeTab === 'overview' && monthly">
      <p class="review-context">
        <template v-if="isCurrentMonth">Current month overview</template>
        <template v-else>Reviewing {{ monthLabel }}</template>
      </p>

      <div class="summary-cards">
        <div class="summary-card">
          <h3>Month budget</h3>
          <p>{{ formatCents(monthly.total_allocated_cents) }}</p>
        </div>
        <div class="summary-card">
          <h3>Month spent</h3>
          <p>{{ formatCents(monthly.total_spent_cents) }}</p>
        </div>
        <div class="summary-card">
          <h3>{{ isCurrentMonth ? 'Month remaining' : 'Under budget' }}</h3>
          <p>{{ formatCents(monthRemainingCents) }}</p>
        </div>
        <div v-if="isCurrentMonth && budgetSummary" class="summary-card">
          <h3>Unallocated</h3>
          <p>{{ formatCents(budgetSummary.unallocated_cents) }}</p>
        </div>
      </div>

      <section class="envelope-balances">
        <h3 class="section-heading">{{ monthLabel }}: envelope balances</h3>
        <p v-if="!monthly.envelopes.length" class="empty-hint">
          No envelopes yet. Add some under the Envelopes tab.
        </p>
        <ul v-else class="balance-list">
          <li
            v-for="envelope in monthly.envelopes"
            :key="envelope.id"
            class="balance-item"
            :class="availableClass(envelope)"
          >
            <div class="balance-info">
              <span class="balance-name">{{ envelope.name }}</span>
              <span class="balance-meta">
                {{ formatCents(envelope.spent_cents) }} spent of
                {{ formatCents(envelope.allocated_cents) }} budget
              </span>
            </div>
            <div class="balance-remaining">
              <span class="balance-label">{{ remainingLabel }}</span>
              <span class="balance-amount">{{ formatCents(remainingCents(envelope)) }}</span>
            </div>
          </li>
        </ul>
      </section>
    </template>

    <!-- Trends -->
    <template v-else-if="activeTab === 'trends'">
      <section class="chart-section card">
        <h3 class="section-heading">Income vs spending</h3>
        <p class="chart-hint">Last 12 months. Includes recurring Incomes sources and extras from History.</p>
        <IncomeVsSpendChart :trends="trends" />
      </section>
    </template>

    <!-- Budget -->
    <template v-else-if="activeTab === 'budget' && monthly">
      <p class="review-context">{{ monthLabel }}</p>
      <section v-if="monthly.envelopes.length" class="chart-section card">
        <h3 class="section-heading">Budget vs spent</h3>
        <EnvelopeBarChart :envelopes="monthly.envelopes" />
      </section>
      <p v-else class="empty-hint">No envelopes yet. Add some under the Envelopes tab.</p>
    </template>

    <!-- Spending -->
    <template v-else-if="activeTab === 'spending' && monthly">
      <p class="review-context">{{ monthLabel }}</p>
      <section v-if="monthly.envelopes.length" class="chart-section card">
        <h3 class="section-heading">Spending breakdown</h3>
        <AllocationDonutChart
          :envelopes="monthly.envelopes"
          value-key="spent_cents"
          :unallocated-cents="Math.max(0, monthly.total_allocated_cents - monthly.total_spent_cents)"
        />
      </section>
      <p v-else class="empty-hint">No envelopes yet. Add some under the Envelopes tab.</p>
    </template>
  </div>
</template>

<style scoped>
.dashboard-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.page-title {
  margin: 0;
}

.dashboard-tabs {
  display: flex;
  gap: 0.35rem;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  margin-bottom: 1rem;
  padding-bottom: 0.15rem;
}

.dashboard-tab {
  flex: 1 0 auto;
  border: none;
  background: #fff;
  border-radius: 999px;
  padding: 0.55rem 0.9rem;
  min-height: 2.5rem;
  font-size: 0.85rem;
  font-weight: 600;
  color: #495057;
  cursor: pointer;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  white-space: nowrap;
}

.dashboard-tab.active {
  background: #3b5bdb;
  color: #fff;
}

.month-picker {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  min-width: 9rem;
  flex-shrink: 0;
}

.month-picker label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #495057;
}

.month-picker select {
  padding: 0.5rem 0.65rem;
  border: 1px solid #ced4da;
  border-radius: 8px;
  font-size: 0.9rem;
  min-height: 2.5rem;
  background: #fff;
}

.review-context {
  margin: 0 0 1rem;
  font-size: 0.9rem;
  color: #868e96;
}

.section-heading {
  margin: 0 0 0.75rem;
  font-size: 1rem;
  color: #495057;
}

.envelope-balances {
  margin-top: 1.25rem;
}

.balance-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.balance-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  background: #fff;
  border-radius: 12px;
  padding: 0.9rem 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.balance-info {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  min-width: 0;
}

.balance-name {
  font-weight: 600;
  font-size: 1rem;
}

.balance-meta {
  font-size: 0.8rem;
  color: #868e96;
}

.balance-remaining {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  flex-shrink: 0;
}

.balance-label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #868e96;
}

.balance-amount {
  font-size: 1.2rem;
  font-weight: 700;
  color: #2b8a3e;
}

.balance-item.empty .balance-amount {
  color: #868e96;
}

.balance-item.over-budget .balance-amount {
  color: #c92a2a;
}

.empty-hint {
  color: #868e96;
  margin: 0;
}

.chart-section {
  margin-bottom: 0;
}

.chart-hint {
  margin: -0.35rem 0 0.75rem;
  font-size: 0.8rem;
  color: #868e96;
}
</style>
