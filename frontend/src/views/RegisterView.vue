<script setup>
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const email = ref('')
const password = ref('')
const joinCode = ref('')
const error = ref('')
const loading = ref(false)
const auth = useAuthStore()
const router = useRouter()

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(email.value, password.value, joinCode.value)
    router.push({ name: 'home' })
  } catch (e) {
    error.value = e.response?.data?.error || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="card auth-card">
      <h2>Create account</h2>
      <p class="subtitle">Password must be at least 8 characters</p>
      <p v-if="error" class="error">{{ error }}</p>
      <form @submit.prevent="submit">
        <div class="form-group">
          <label for="email">Email</label>
          <input id="email" v-model="email" type="email" required autocomplete="email" />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input id="password" v-model="password" type="password" required minlength="8" autocomplete="new-password" />
        </div>
        <div class="form-group">
          <label for="join-code">Family join code</label>
          <input
            id="join-code"
            v-model="joinCode"
            type="text"
            placeholder="Leave blank to create a new family budget"
            autocomplete="off"
            maxlength="8"
            class="join-code-input"
          />
          <p class="field-hint">
            Have a code from a family member? Enter it to share their budget. Leave empty to start a new one — you'll get a code to share after signing up.
          </p>
        </div>
        <button class="btn btn-block" type="submit" :disabled="loading">
          {{ loading ? 'Creating...' : 'Register' }}
        </button>
      </form>
      <p class="footer-link">
        Already have an account?
        <RouterLink to="/login">Log in</RouterLink>
      </p>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 1rem;
}

.auth-card {
  width: 100%;
  max-width: 420px;
}

.subtitle {
  color: #495057;
  margin-top: -0.5rem;
}

.join-code-input {
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 600;
}

.field-hint {
  margin: 0;
  font-size: 0.8rem;
  color: #868e96;
  line-height: 1.4;
}

.footer-link {
  margin-top: 1rem;
}
</style>
