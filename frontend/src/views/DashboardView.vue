<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api/client'
import { formatCents } from '@/utils/money'
import AllocationDonutChart from '@/components/charts/AllocationDonutChart.vue'
import EnvelopeBarChart from '@/components/charts/EnvelopeBarChart.vue'

const summary = ref(null)
const error = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await api.get('/api/summary')
    summary.value = data
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load dashboard'
  } finally {
    loading.value = false
  }
})

function availableClass(cents) {
  if (cents < 0) return 'over-budget'
  if (cents === 0) return 'empty'
  return ''
}
</script>

<template>
  <div class="dashboard">
    <h2 class="page-title">Dashboard</h2>
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="loading">Loading...</p>

    <template v-else-if="summary">
      <section class="envelope-balances">
        <h3 class="section-heading">Envelope balances</h3>
        <p v-if="!summary.envelopes.length" class="empty-hint">No envelopes yet. Add some under the Envelopes tab.</p>
        <ul v-else class="balance-list">
          <li
            v-for="envelope in summary.envelopes"
            :key="envelope.id"
            class="balance-item"
            :class="availableClass(envelope.available_cents)"
          >
            <div class="balance-info">
              <span class="balance-name">{{ envelope.name }}</span>
              <span class="balance-meta">
                {{ formatCents(envelope.spent_cents) }} spent of {{ formatCents(envelope.allocated_cents) }}
              </span>
            </div>
            <div class="balance-remaining">
              <span class="balance-label">Left</span>
              <span class="balance-amount">{{ formatCents(envelope.available_cents) }}</span>
            </div>
          </li>
        </ul>
      </section>

      <section v-if="summary.envelopes.length" class="chart-section card">
        <h3 class="section-heading">Allocated vs spent</h3>
        <EnvelopeBarChart :envelopes="summary.envelopes" />
      </section>

      <div class="summary-cards">
        <div class="summary-card">
          <h3>Total Income</h3>
          <p>{{ formatCents(summary.total_income_cents) }}</p>
        </div>
        <div class="summary-card">
          <h3>Total Allocated</h3>
          <p>{{ formatCents(summary.total_allocated_cents) }}</p>
        </div>
        <div class="summary-card">
          <h3>Unallocated</h3>
          <p>{{ formatCents(summary.unallocated_cents) }}</p>
        </div>
        <div class="summary-card">
          <h3>Total Spent</h3>
          <p>{{ formatCents(summary.total_spent_cents) }}</p>
        </div>
      </div>

      <section class="allocation-section card">
        <h3 class="section-heading">Allocation breakdown</h3>
        <AllocationDonutChart :envelopes="summary.envelopes" :unallocated-cents="summary.unallocated_cents" />
      </section>
    </template>
  </div>
</template>

<style scoped>
.section-heading {
  margin: 0 0 0.75rem;
  font-size: 1rem;
  color: #495057;
}

.envelope-balances {
  margin-bottom: 1.25rem;
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
  margin-bottom: 1.25rem;
}

.allocation-section {
  margin-top: 1.5rem;
}
</style>
