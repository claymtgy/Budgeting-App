<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api/client'
import { centsFromDigitInput, formatCents, formatDigitInputAsDollars, todayISO } from '@/utils/money'

const expenses = ref([])
const envelopes = ref([])
const savedPlaces = ref([])
const error = ref('')
const success = ref('')
const loading = ref(true)
const submitting = ref(false)
const showHistory = ref(false)
const showPlaceSuggestions = ref(false)

const envelopeId = ref('')
const amountCents = ref(0)
const amountDisplay = ref('')
const description = ref('')
const place = ref('')
const expenseDate = ref(todayISO())

const envelopeMap = computed(() =>
  Object.fromEntries(envelopes.value.map((e) => [e.id, e.name]))
)

const recentExpenses = computed(() => expenses.value.slice(0, 10))

const placeSuggestions = computed(() => {
  const query = place.value.trim().toLowerCase()
  if (!query) return savedPlaces.value.slice(0, 8)
  return savedPlaces.value.filter((p) => p.toLowerCase().includes(query)).slice(0, 8)
})

async function loadPlaces() {
  try {
    const { data } = await api.get('/api/expenses/places')
    savedPlaces.value = data
  } catch {
    // Non-critical; typeahead just won't have history yet.
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
    savedPlaces.value = [...savedPlaces.value, trimmed].sort((a, b) =>
      a.localeCompare(b)
    )
  }
}

function onAmountInput(event) {
  amountCents.value = centsFromDigitInput(event.target.value)
  amountDisplay.value = formatDigitInputAsDollars(event.target.value)
  event.target.value = amountDisplay.value
}

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
      amount_cents: amountCents.value,
      description: description.value,
      place: place.value.trim(),
      expense_date: expenseDate.value
    })
    rememberPlace(place.value)
    amountCents.value = 0
    amountDisplay.value = ''
    description.value = ''
    place.value = ''
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

onMounted(() => {
  load()
  loadPlaces()
})
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

        <div class="form-group place-group">
          <label for="place">Place</label>
          <input
            id="place"
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
          <label for="date">Date</label>
          <input id="date" v-model="expenseDate" type="date" required />
        </div>

        <button
          class="btn btn-block"
          type="submit"
          :disabled="!envelopes.length || submitting || amountCents <= 0"
        >
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
              <template v-if="expense.place"> · {{ expense.place }}</template>
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
