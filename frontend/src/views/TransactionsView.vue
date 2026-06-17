<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import api from '@/api/client'
import {
  centsFromDigitInput,
  formatCents,
  formatDigitInputAsDollars,
  startOfMonthISO,
  todayISO
} from '@/utils/money'

const envelopes = ref([])
const savedPlaces = ref([])
const error = ref('')
const success = ref('')
const loading = ref(false)
const submitting = ref(false)
const showPlaceSuggestions = ref(false)

const entryType = ref('expense')
const fromDate = ref(startOfMonthISO())
const toDate = ref(todayISO())

const amountCents = ref(0)
const amountDisplay = ref('')
const description = ref('')
const place = ref('')
const envelopeId = ref('')
const transactionDate = ref(todayISO())

const expenses = ref([])
const incomeReceipts = ref([])

const envelopeMap = computed(() =>
  Object.fromEntries(envelopes.value.map((e) => [e.id, e.name]))
)

const placeSuggestions = computed(() => {
  const query = place.value.trim().toLowerCase()
  if (!query) return savedPlaces.value.slice(0, 8)
  return savedPlaces.value.filter((p) => p.toLowerCase().includes(query)).slice(0, 8)
})

const transactions = computed(() => {
  const items = [
    ...expenses.value.map((tx) => ({
      key: `expense-${tx.id}`,
      type: 'expense',
      id: tx.id,
      amount_cents: tx.amount_cents,
      date: tx.expense_date,
      description: tx.description,
      place: tx.place,
      envelope_id: tx.envelope_id,
      voided: tx.voided
    })),
    ...incomeReceipts.value.map((tx) => ({
      key: `income-${tx.id}`,
      type: 'income',
      id: tx.id,
      amount_cents: tx.amount_cents,
      date: tx.income_date,
      description: tx.description,
      place: '',
      envelope_id: null,
      voided: tx.voided
    }))
  ]

  return items.sort((a, b) => {
    if (a.date !== b.date) return a.date < b.date ? 1 : -1
    return a.type === b.type ? 0 : a.type === 'expense' ? -1 : 1
  })
})

const totalCents = computed(() =>
  transactions.value.reduce((sum, tx) => {
    if (tx.voided) return sum
    return tx.type === 'income' ? sum + tx.amount_cents : sum - tx.amount_cents
  }, 0)
)

function onAmountInput(event) {
  amountCents.value = centsFromDigitInput(event.target.value)
  amountDisplay.value = formatDigitInputAsDollars(event.target.value)
  event.target.value = amountDisplay.value
}

function resetForm() {
  amountCents.value = 0
  amountDisplay.value = ''
  description.value = ''
  place.value = ''
  transactionDate.value = todayISO()
}

async function loadPlaces() {
  try {
    const { data } = await api.get('/api/expenses/places')
    savedPlaces.value = data
  } catch {
    // Non-critical.
  }
}

function onPlaceFocus() {
  if (placeSuggestions.value.length) showPlaceSuggestions.value = true
}

function onPlaceBlur() {
  setTimeout(() => {
    showPlaceSuggestions.value = false
  }, 150)
}

function selectPlace(suggestion) {
  place.value = suggestion
  showPlaceSuggestions.value = false
}

function rememberPlace(value) {
  const trimmed = value.trim()
  if (!trimmed) return
  if (!savedPlaces.value.some((p) => p.toLowerCase() === trimmed.toLowerCase())) {
    savedPlaces.value = [...savedPlaces.value, trimmed].sort((a, b) => a.localeCompare(b))
  }
}

async function loadEnvelopes() {
  const { data } = await api.get('/api/envelopes')
  envelopes.value = data
  if (!envelopeId.value && envelopes.value.length) {
    envelopeId.value = envelopes.value[0].id
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [expenseRes, incomeRes] = await Promise.all([
      api.get('/api/expenses', { params: { from: fromDate.value, to: toDate.value } }),
      api.get('/api/income-receipts', { params: { from: fromDate.value, to: toDate.value } })
    ])
    expenses.value = expenseRes.data
    incomeReceipts.value = incomeRes.data
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load transactions'
  } finally {
    loading.value = false
  }
}

async function submit() {
  error.value = ''
  success.value = ''
  submitting.value = true
  try {
    if (entryType.value === 'expense') {
      await api.post('/api/expenses', {
        envelope_id: envelopeId.value,
        amount_cents: amountCents.value,
        description: description.value,
        place: place.value.trim(),
        expense_date: transactionDate.value
      })
      rememberPlace(place.value)
      success.value = 'Expense added'
    } else {
      await api.post('/api/income-receipts', {
        amount_cents: amountCents.value,
        description: description.value.trim(),
        income_date: transactionDate.value
      })
      success.value = 'Income added'
    }
    resetForm()
    await load()
    setTimeout(() => {
      success.value = ''
    }, 2500)
  } catch (e) {
    error.value =
      e.response?.data?.error ||
      (entryType.value === 'expense' ? 'Could not create expense' : 'Could not create income')
  } finally {
    submitting.value = false
  }
}

async function voidTransaction(tx) {
  const label = tx.type === 'income' ? 'income' : 'expense'
  if (!confirm(`Void this ${label}?`)) return
  try {
    if (tx.type === 'income') {
      await api.put(`/api/income-receipts/${tx.id}/void`)
    } else {
      await api.put(`/api/expenses/${tx.id}/void`)
    }
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || `Could not void ${label}`
  }
}

function formatSignedAmount(tx) {
  const formatted = formatCents(tx.amount_cents)
  return tx.type === 'income' ? `+${formatted}` : `-${formatted}`
}

watch(entryType, () => {
  error.value = ''
})

onMounted(async () => {
  try {
    await Promise.all([loadEnvelopes(), loadPlaces()])
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load data'
  }
})
</script>

<template>
  <div class="transactions">
    <h2 class="page-title">Transactions</h2>
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" class="success">{{ success }}</p>

    <section class="card add-card">
      <div class="type-toggle" role="tablist" aria-label="Transaction type">
        <button
          type="button"
          class="type-btn"
          :class="{ active: entryType === 'expense' }"
          @click="entryType = 'expense'"
        >
          Expense
        </button>
        <button
          type="button"
          class="type-btn"
          :class="{ active: entryType === 'income' }"
          @click="entryType = 'income'"
        >
          Income
        </button>
      </div>

      <form @submit.prevent="submit">
        <div class="form-group">
          <label for="tx-amount">Amount</label>
          <div class="amount-input-wrap" :class="{ income: entryType === 'income' }">
            <span class="currency">$</span>
            <input
              id="tx-amount"
              :value="amountDisplay"
              class="amount-input"
              type="text"
              inputmode="numeric"
              autocomplete="off"
              placeholder="0.00"
              required
              @input="onAmountInput"
            />
          </div>
        </div>

        <div v-if="entryType === 'expense'" class="form-group">
          <label for="tx-envelope">Envelope</label>
          <select id="tx-envelope" v-model="envelopeId" required :disabled="!envelopes.length">
            <option v-if="!envelopes.length" value="">No envelopes yet</option>
            <option v-for="envelope in envelopes" :key="envelope.id" :value="envelope.id">
              {{ envelope.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="tx-description">{{ entryType === 'income' ? 'Source' : 'Description' }}</label>
          <input
            id="tx-description"
            v-model="description"
            type="text"
            :placeholder="entryType === 'income' ? 'Paycheck, refund, etc.' : 'What was this for?'"
            autocomplete="off"
          />
        </div>

        <div v-if="entryType === 'expense'" class="form-group place-group">
          <label for="tx-place">Place</label>
          <input
            id="tx-place"
            v-model="place"
            type="text"
            placeholder="Where was this charge?"
            autocomplete="off"
            @focus="onPlaceFocus"
            @blur="onPlaceBlur"
            @input="showPlaceSuggestions = placeSuggestions.length > 0"
          />
          <ul
            v-if="showPlaceSuggestions && placeSuggestions.length"
            class="place-suggestions"
            role="listbox"
          >
            <li
              v-for="suggestion in placeSuggestions"
              :key="suggestion"
              role="option"
              @mousedown.prevent="selectPlace(suggestion)"
            >
              {{ suggestion }}
            </li>
          </ul>
        </div>

        <div class="form-group">
          <label for="tx-date">Date</label>
          <input id="tx-date" v-model="transactionDate" type="date" required />
        </div>

        <button
          class="btn btn-block"
          type="submit"
          :disabled="
            submitting ||
            amountCents <= 0 ||
            (entryType === 'expense' && !envelopes.length)
          "
        >
          {{ submitting ? 'Saving...' : entryType === 'income' ? 'Add income' : 'Add expense' }}
        </button>
      </form>

      <p v-if="entryType === 'expense' && !envelopes.length" class="hint">
        Create an envelope first under the Envelopes tab.
      </p>
    </section>

    <div class="card filters-card">
      <div class="date-filters">
        <div class="form-group">
          <label for="from-date">From</label>
          <input id="from-date" v-model="fromDate" type="date" required @change="load" />
        </div>
        <div class="form-group">
          <label for="to-date">To</label>
          <input id="to-date" v-model="toDate" type="date" required @change="load" />
        </div>
      </div>
      <p class="summary-line">
        <span>{{ transactions.length }} transaction{{ transactions.length === 1 ? '' : 's' }}</span>
        <span class="summary-total" :class="{ positive: totalCents > 0, negative: totalCents < 0 }">
          Net: {{ formatCents(Math.abs(totalCents)) }}{{ totalCents < 0 ? ' spent' : totalCents > 0 ? ' gained' : '' }}
        </span>
      </p>
    </div>

    <div class="card">
      <p v-if="loading">Loading...</p>
      <p v-else-if="!transactions.length" class="empty">No transactions in this date range.</p>

      <div v-else class="data-list mobile-only">
        <article
          v-for="tx in transactions"
          :key="tx.key"
          class="data-list-item"
          :class="{ voided: tx.voided, income: tx.type === 'income' }"
        >
          <div class="row">
            <div>
              <p class="title" :class="{ 'amount-income': tx.type === 'income' }">
                {{ formatSignedAmount(tx) }}
              </p>
              <p class="subtitle">
                <span class="type-badge" :class="tx.type">{{ tx.type }}</span>
                {{ tx.date }}
                <template v-if="tx.type === 'expense'">
                  · {{ envelopeMap[tx.envelope_id] || 'Unknown' }}
                  <template v-if="tx.place"> · {{ tx.place }}</template>
                </template>
              </p>
              <p v-if="tx.description" class="description">{{ tx.description }}</p>
            </div>
            <div class="actions">
              <span v-if="tx.voided" class="voided-label">Voided</span>
              <button
                v-else
                class="btn btn-sm btn-danger"
                type="button"
                @click="voidTransaction(tx)"
              >
                Void
              </button>
            </div>
          </div>
        </article>
      </div>

      <div v-if="!loading && transactions.length" class="table-wrap desktop-only">
        <table>
          <thead>
            <tr>
              <th>Type</th>
              <th>Date</th>
              <th>Amount</th>
              <th>Envelope</th>
              <th>Place</th>
              <th>Description</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="tx in transactions"
              :key="tx.key"
              :class="{ voided: tx.voided, income: tx.type === 'income' }"
            >
              <td>
                <span class="type-badge" :class="tx.type">{{ tx.type }}</span>
              </td>
              <td>{{ tx.date }}</td>
              <td :class="{ 'amount-income': tx.type === 'income' }">{{ formatSignedAmount(tx) }}</td>
              <td>{{ tx.type === 'expense' ? envelopeMap[tx.envelope_id] || 'Unknown' : '—' }}</td>
              <td>{{ tx.place || '—' }}</td>
              <td>{{ tx.description || '—' }}</td>
              <td>{{ tx.voided ? 'Voided' : 'Active' }}</td>
              <td>
                <button
                  v-if="!tx.voided"
                  class="btn btn-sm btn-danger"
                  type="button"
                  @click="voidTransaction(tx)"
                >
                  Void
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.transactions {
  max-width: 720px;
  margin: 0 auto;
}

.add-card {
  margin-bottom: 1rem;
}

.type-toggle {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
  margin-bottom: 1rem;
  padding: 0.25rem;
  background: #f1f3f5;
  border-radius: 10px;
}

.type-btn {
  border: none;
  background: transparent;
  border-radius: 8px;
  padding: 0.65rem;
  min-height: 2.75rem;
  font-weight: 600;
  color: #495057;
  cursor: pointer;
}

.type-btn.active {
  background: #fff;
  color: #3b5bdb;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.type-btn.active:last-child {
  color: #2b8a3e;
}

.amount-input-wrap {
  display: flex;
  align-items: center;
  border: 2px solid #3b5bdb;
  border-radius: 10px;
  padding: 0 0.75rem;
  background: #f8f9ff;
}

.amount-input-wrap.income {
  border-color: #2b8a3e;
  background: #f4fce3;
}

.amount-input-wrap.income .currency {
  color: #2b8a3e;
}

.currency {
  font-size: 1.25rem;
  font-weight: 600;
  color: #3b5bdb;
}

.amount-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 1.5rem;
  font-weight: 600;
  padding: 0.65rem 0.5rem;
  min-height: 3rem;
}

.amount-input:focus {
  outline: none;
}

.place-group {
  position: relative;
}

.place-suggestions {
  position: absolute;
  z-index: 10;
  top: calc(100% + 0.25rem);
  left: 0;
  right: 0;
  margin: 0;
  padding: 0.25rem 0;
  list-style: none;
  background: #fff;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  max-height: 12rem;
  overflow-y: auto;
}

.place-suggestions li {
  padding: 0.6rem 0.85rem;
  cursor: pointer;
}

.place-suggestions li:hover {
  background: #f1f3f5;
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

.filters-card {
  margin-bottom: 1rem;
}

.date-filters {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.date-filters .form-group {
  margin-bottom: 0;
}

.summary-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  margin: 1rem 0 0;
  font-size: 0.9rem;
  color: #495057;
}

.summary-total {
  font-weight: 600;
  color: #1a1a2e;
}

.summary-total.positive {
  color: #2b8a3e;
}

.summary-total.negative {
  color: #c92a2a;
}

.empty {
  margin: 0;
  color: #868e96;
}

.description {
  margin: 0.25rem 0 0;
  font-size: 0.9rem;
}

.voided-label {
  font-size: 0.8rem;
  color: #868e96;
}

.type-badge {
  display: inline-block;
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  padding: 0.1rem 0.35rem;
  border-radius: 4px;
  margin-right: 0.25rem;
}

.type-badge.expense {
  background: #fff5f5;
  color: #c92a2a;
}

.type-badge.income {
  background: #ebfbee;
  color: #2b8a3e;
}

.amount-income {
  color: #2b8a3e;
}

.data-list-item.income .title {
  color: #2b8a3e;
}
</style>
