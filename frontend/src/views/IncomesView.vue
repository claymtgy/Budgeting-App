<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api/client'
import { dollarsToCents, formatCents } from '@/utils/money'

const incomes = ref([])
const error = ref('')
const loading = ref(true)

const name = ref('')
const amount = ref('')
const period = ref('monthly')
const editingId = ref(null)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await api.get('/api/incomes')
    incomes.value = data
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load incomes'
  } finally {
    loading.value = false
  }
}

function resetForm() {
  editingId.value = null
  name.value = ''
  amount.value = ''
  period.value = 'monthly'
}

function startEdit(income) {
  editingId.value = income.id
  name.value = income.name
  amount.value = (income.amount_cents / 100).toFixed(2)
  period.value = income.period
}

async function submit() {
  error.value = ''
  const payload = {
    name: name.value,
    amount_cents: dollarsToCents(amount.value),
    period: period.value
  }

  try {
    if (editingId.value) {
      await api.put(`/api/incomes/${editingId.value}`, payload)
    } else {
      await api.post('/api/incomes', payload)
    }
    resetForm()
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not save income'
  }
}

async function remove(id) {
  if (!confirm('Delete this income?')) return
  try {
    await api.delete(`/api/incomes/${id}`)
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not delete income'
  }
}

onMounted(load)
</script>

<template>
  <div>
    <h2 class="page-title">Incomes</h2>
    <p v-if="error" class="error">{{ error }}</p>

    <div class="grid-2">
      <div class="card">
        <h3>{{ editingId ? 'Edit income' : 'Add income' }}</h3>
        <form @submit.prevent="submit">
          <div class="form-group">
            <label>Name</label>
            <input v-model="name" required />
          </div>
          <div class="form-group">
            <label>Amount (USD)</label>
            <input v-model="amount" type="number" inputmode="decimal" step="0.01" min="0" required />
          </div>
          <div class="form-group">
            <label>Period</label>
            <select v-model="period">
              <option value="monthly">Monthly</option>
              <option value="biweekly">Biweekly</option>
              <option value="weekly">Weekly</option>
            </select>
          </div>
          <button class="btn btn-block" type="submit">{{ editingId ? 'Update' : 'Add' }}</button>
          <button v-if="editingId" class="btn btn-secondary btn-block" type="button" @click="resetForm">Cancel</button>
        </form>
      </div>

      <div class="card">
        <h3>Your incomes</h3>
        <p v-if="loading">Loading...</p>

        <div v-else class="data-list mobile-only">
          <article v-for="income in incomes" :key="income.id" class="data-list-item">
            <div class="row">
              <div>
                <p class="title">{{ income.name }}</p>
                <p class="subtitle">{{ formatCents(income.amount_cents) }} · {{ income.period }}</p>
              </div>
              <div class="actions">
                <button class="btn btn-sm btn-secondary" type="button" @click="startEdit(income)">Edit</button>
                <button class="btn btn-sm btn-danger" type="button" @click="remove(income.id)">Del</button>
              </div>
            </div>
          </article>
        </div>

        <div v-if="!loading" class="table-wrap desktop-only">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Amount</th>
                <th>Period</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="income in incomes" :key="income.id">
                <td>{{ income.name }}</td>
                <td>{{ formatCents(income.amount_cents) }}</td>
                <td>{{ income.period }}</td>
                <td>
                  <button class="btn btn-sm btn-secondary" type="button" @click="startEdit(income)">Edit</button>
                  <button class="btn btn-sm btn-danger" type="button" @click="remove(income.id)">Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
