<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api/client'
import { dollarsToCents, formatCents, todayISO } from '@/utils/money'

const expenses = ref([])
const envelopes = ref([])
const error = ref('')
const success = ref('')
const loading = ref(true)
const submitting = ref(false)
const showHistory = ref(false)

const envelopeId = ref('')
const amount = ref('')
const description = ref('')
const expenseDate = ref(todayISO())

const envelopeMap = computed(() =>
  Object.fromEntries(envelopes.value.map((e) => [e.id, e.name]))
)

const recentExpenses = computed(() => expenses.value.slice(0, 10))

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [expenseRes, envelopeRes] = await Promise.all([
      api.get('/api/expenses'),
      api.get('/api/envelopes')
    ])
    expenses.value = expenseRes.data
    envelopes.value = envelopeRes.data
    if (!envelopeId.value && envelopes.value.length) {
      envelopeId.value = envelopes.value[0].id
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load data'
  } finally {
    loading.value = false
  }
}

async function submit() {
  error.value = ''
  success.value = ''
  submitting.value = true
  try {
    await api.post('/api/expenses', {
      envelope_id: envelopeId.value,
      amount_cents: dollarsToCents(amount.value),
      description: description.value,
      expense_date: expenseDate.value
    })
    amount.value = ''
    description.value = ''
    expenseDate.value = todayISO()
    success.value = 'Expense added'
    await load()
    setTimeout(() => {
      success.value = ''
    }, 2500)
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not create expense'
  } finally {
    submitting.value = false
  }
}

async function voidExpense(id) {
  if (!confirm('Void this expense?')) return
  try {
    await api.put(`/api/expenses/${id}/void`)
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not void expense'
  }
}

onMounted(load)
</script>

<template>
  <div class="home">
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" class="success">{{ success }}</p>

    <section class="card add-expense-card">
      <h2 class="page-heading">Add expense</h2>

      <form @submit.prevent="submit">
        <div class="form-group">
          <label for="amount">Amount</label>
          <div class="amount-input-wrap">
            <span class="currency">$</span>
            <input
              id="amount"
              v-model="amount"
              class="amount-input"
              type="number"
              inputmode="decimal"
              step="0.01"
              min="0.01"
              placeholder="0.00"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="envelope">Envelope</label>
          <select id="envelope" v-model="envelopeId" required :disabled="!envelopes.length">
            <option v-if="!envelopes.length" value="">No envelopes yet</option>
            <option v-for="envelope in envelopes" :key="envelope.id" :value="envelope.id">
              {{ envelope.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="description">Description</label>
          <input
            id="description"
            v-model="description"
            type="text"
            placeholder="What was this for?"
            autocomplete="off"
          />
        </div>

        <div class="form-group">
          <label for="date">Date</label>
          <input id="date" v-model="expenseDate" type="date" required />
        </div>

        <button class="btn btn-block" type="submit" :disabled="!envelopes.length || submitting">
          {{ submitting ? 'Saving...' : 'Add expense' }}
        </button>
      </form>

      <p v-if="!loading && !envelopes.length" class="hint">
        Create an envelope first under the Envelopes tab.
      </p>
    </section>

    <section v-if="!loading && expenses.length" class="history-section">
      <button class="history-toggle" type="button" @click="showHistory = !showHistory">
        <span>Recent expenses ({{ expenses.length }})</span>
        <span class="chevron">{{ showHistory ? '▲' : '▼' }}</span>
      </button>

      <div v-show="showHistory" class="expense-list">
        <article
          v-for="expense in recentExpenses"
          :key="expense.id"
          class="expense-card"
          :class="{ voided: expense.voided }"
        >
          <div class="expense-main">
            <p class="expense-amount">{{ formatCents(expense.amount_cents) }}</p>
            <p class="expense-meta">
              {{ envelopeMap[expense.envelope_id] || 'Unknown' }}
              · {{ expense.expense_date }}
            </p>
            <p v-if="expense.description" class="expense-desc">{{ expense.description }}</p>
          </div>
          <button
            v-if="!expense.voided"
            class="btn btn-sm btn-danger"
            type="button"
            @click="voidExpense(expense.id)"
          >
            Void
          </button>
          <span v-else class="voided-label">Voided</span>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home {
  max-width: 480px;
  margin: 0 auto;
}

.page-heading {
  margin: 0 0 1.25rem;
  font-size: 1.35rem;
}

.add-expense-card {
  margin-bottom: 1rem;
}

.amount-input-wrap {
  display: flex;
  align-items: center;
  border: 2px solid #3b5bdb;
  border-radius: 10px;
  padding: 0 0.75rem;
  background: #f8f9ff;
}

.currency {
  font-size: 1.5rem;
  font-weight: 600;
  color: #3b5bdb;
}

.amount-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 2rem;
  font-weight: 600;
  padding: 0.75rem 0.5rem;
  min-height: 3.5rem;
}

.amount-input:focus {
  outline: none;
}

.hint {
  margin: 1rem 0 0;
  color: #868e96;
  font-size: 0.9rem;
}

.success {
  color: #2b8a3e;
  background: #ebfbee;
  padding: 0.65rem 0.85rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.history-section {
  margin-top: 0.5rem;
}

.history-toggle {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.85rem 1rem;
  border: none;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  font-weight: 600;
  cursor: pointer;
  min-height: 3rem;
}

.chevron {
  color: #868e96;
  font-size: 0.75rem;
}

.expense-list {
  margin-top: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.expense-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  background: #fff;
  border-radius: 10px;
  padding: 0.85rem 1rem;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.expense-card.voided {
  opacity: 0.55;
}

.expense-main {
  flex: 1;
  min-width: 0;
}

.expense-amount {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.expense-meta {
  margin: 0.15rem 0 0;
  font-size: 0.85rem;
  color: #868e96;
}

.expense-desc {
  margin: 0.25rem 0 0;
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.voided-label {
  font-size: 0.8rem;
  color: #868e96;
}
</style>
