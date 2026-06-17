<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const copied = ref(false)

const navItems = [
  { to: '/', name: 'home', label: 'Add', icon: '+' },
  { to: '/transactions', name: 'transactions', label: 'History', icon: '≡' },
  { to: '/dashboard', name: 'dashboard', label: 'Dashboard', icon: '◉' },
  { to: '/incomes', name: 'incomes', label: 'Incomes', icon: '$' },
  { to: '/envelopes', name: 'envelopes', label: 'Envelopes', icon: '✉' }
]

onMounted(() => {
  if (auth.isAuthenticated && !auth.joinCode) {
    auth.fetchMe()
  }
})

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}

function isNavActive(item) {
  return route.name === item.name
}

async function copyJoinCode() {
  if (!auth.joinCode) return
  try {
    await navigator.clipboard.writeText(auth.joinCode)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    // clipboard unavailable
  }
}
</script>

<template>
  <div class="app-shell">
    <header class="header">
      <div class="header-inner">
        <h1 class="app-title">Envelope Budget</h1>
        <div class="header-actions">
          <button
            v-if="auth.joinCode"
            class="join-code-btn"
            type="button"
            :title="copied ? 'Copied!' : 'Copy family join code'"
            @click="copyJoinCode"
          >
            <span class="join-code-label">Family code</span>
            <span class="join-code-value">{{ auth.joinCode }}</span>
          </button>
          <button class="btn-icon" type="button" aria-label="Log out" @click="logout">⎋</button>
        </div>
      </div>
    </header>

    <main class="main-content">
      <RouterView />
    </main>

    <nav class="bottom-nav" aria-label="Main navigation">
      <RouterLink
        v-for="item in navItems"
        :key="item.name"
        :to="item.to"
        class="bottom-nav-item"
        :class="{ active: isNavActive(item) }"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span class="nav-label">{{ item.label }}</span>
      </RouterLink>
    </nav>
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
}

.header {
  background: #1a1a2e;
  color: #fff;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0.85rem 1rem;
  padding-top: max(0.85rem, env(safe-area-inset-top));
}

.app-title {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  flex-shrink: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.join-code-btn {
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  border-radius: 8px;
  padding: 0.35rem 0.6rem;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  min-width: 0;
}

.join-code-label {
  font-size: 0.6rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  opacity: 0.75;
}

.join-code-value {
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.1em;
}

.btn-icon {
  border: none;
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 8px;
  font-size: 1.1rem;
  cursor: pointer;
  flex-shrink: 0;
}

.main-content {
  flex: 1;
  padding: 1rem;
  padding-bottom: calc(5.5rem + env(safe-area-inset-bottom));
  max-width: 1100px;
  width: 100%;
  margin: 0 auto;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  background: #fff;
  border-top: 1px solid #dee2e6;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 20;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.06);
}

.bottom-nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.15rem;
  padding: 0.55rem 0.25rem;
  min-height: 3.5rem;
  color: #868e96;
  text-decoration: none;
  font-size: 0.7rem;
  -webkit-tap-highlight-color: transparent;
}

.bottom-nav-item.active {
  color: #3b5bdb;
  font-weight: 600;
}

.nav-icon {
  font-size: 1.25rem;
  line-height: 1;
}

.nav-label {
  line-height: 1.2;
}

@media (min-width: 768px) {
  .main-content {
    padding: 1.5rem;
    padding-bottom: calc(5.5rem + env(safe-area-inset-bottom));
  }

  .bottom-nav-item {
    font-size: 0.75rem;
    min-height: 4rem;
  }
}
</style>
