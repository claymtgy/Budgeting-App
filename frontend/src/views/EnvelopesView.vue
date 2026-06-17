<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api/client'
import { dollarsToCents, formatCents } from '@/utils/money'

const envelopes = ref([])
const error = ref('')
const loading = ref(true)

const name = ref('')
const amount = ref('')
const editingId = ref(null)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await api.get('/api/envelopes')
    envelopes.value = data
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load envelopes'
  } finally {
    loading.value = false
  }
}

function resetForm() {
  editingId.value = null
  name.value = ''
  amount.value = ''
}

function startEdit(envelope) {
  editingId.value = envelope.id
  name.value = envelope.name
  amount.value = (envelope.allocated_cents / 100).toFixed(2)
}

async function submit() {
  error.value = ''
  const payload = {
    name: name.value,
    allocated_cents: dollarsToCents(amount.value)
  }

  try {
    if (editingId.value) {
      await api.put(`/api/envelopes/${editingId.value}`, payload)
    } else {
      await api.post('/api/envelopes', payload)
    }
    resetForm()
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not save envelope'
  }
}

async function remove(id) {
  if (!confirm('Delete this envelope?')) return
  try {
    await api.delete(`/api/envelopes/${id}`)
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not delete envelope'
  }
}

onMounted(load)
</script>

<template>
  <div>
    <h2 class="page-title">Envelopes</h2>
    <p v-if="error" class="error">{{ error }}</p>

    <div class="grid-2">
      <div class="card">
        <h3>{{ editingId ? 'Edit envelope' : 'Add envelope' }}</h3>
        <form @submit.prevent="submit">
          <div class="form-group">
            <label>Name</label>
            <input v-model="name" required />
          </div>
          <div class="form-group">
            <label>Monthly allocation (USD)</label>
            <input v-model="amount" type="number" inputmode="decimal" step="0.01" min="0" required />
          </div>
          <button class="btn btn-block" type="submit">{{ editingId ? 'Update' : 'Add' }}</button>
          <button v-if="editingId" class="btn btn-secondary btn-block" type="button" @click="resetForm">Cancel</button>
        </form>
      </div>

      <div class="card">
        <h3>Your envelopes</h3>
        <p v-if="loading">Loading...</p>

        <div v-else class="data-list mobile-only">
          <article v-for="envelope in envelopes" :key="envelope.id" class="data-list-item">
            <div class="row">
              <div>
                <p class="title">{{ envelope.name }}</p>
                <p class="subtitle">
                  {{ formatCents(envelope.allocated_cents) }}/mo · Balance {{ formatCents(envelope.balance_cents) }}
                </p>
              </div>
              <div class="actions">
                <button class="btn btn-sm btn-secondary" type="button" @click="startEdit(envelope)">Edit</button>
                <button class="btn btn-sm btn-danger" type="button" @click="remove(envelope.id)">Del</button>
              </div>
            </div>
          </article>
        </div>

        <div v-if="!loading" class="table-wrap desktop-only">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Monthly</th>
                <th>Balance</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="envelope in envelopes" :key="envelope.id">
                <td>{{ envelope.name }}</td>
                <td>{{ formatCents(envelope.allocated_cents) }}</td>
                <td>{{ formatCents(envelope.balance_cents) }}</td>
                <td>
                  <button class="btn btn-sm btn-secondary" type="button" @click="startEdit(envelope)">Edit</button>
                  <button class="btn btn-sm btn-danger" type="button" @click="remove(envelope.id)">Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
