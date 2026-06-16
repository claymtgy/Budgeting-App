export function formatCents(cents) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format((cents || 0) / 100)
}

export function dollarsToCents(value) {
  const num = Number(value)
  if (Number.isNaN(num)) return 0
  return Math.round(num * 100)
}

export function todayISO() {
  return new Date().toISOString().slice(0, 10)
}
