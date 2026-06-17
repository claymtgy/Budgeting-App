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

/** Strip non-digits and interpret as cents (5748 → 5748 cents = $57.48). */
export function centsFromDigitInput(value) {
  const digits = String(value ?? '').replace(/\D/g, '')
  if (!digits) return 0
  return parseInt(digits, 10)
}

/** Format digit input as a dollar string with decimal (5748 → "57.48"). */
export function formatDigitInputAsDollars(value) {
  const cents = centsFromDigitInput(value)
  if (!cents) return ''
  return (cents / 100).toFixed(2)
}

export function todayISO() {
  return new Date().toISOString().slice(0, 10)
}

export function startOfMonthISO() {
  const date = new Date()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  return `${date.getFullYear()}-${month}-01`
}
