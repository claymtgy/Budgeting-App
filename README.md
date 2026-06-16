# Envelope Budgeting App

A Dockerized envelope budgeting application with a Go API, Vue 3 frontend, PostgreSQL database, JWT authentication, and D3 dashboard charts.

## Stack

- **Backend:** Go, chi, pgx, golang-migrate, JWT
- **Frontend:** Vue 3, Vite, Pinia, Vue Router, Axios, D3, PWA (vite-plugin-pwa)
- **Database:** PostgreSQL 16

## Quick start

1. Copy environment variables:

```bash
cp .env.example .env
```

2. Start all services:

```bash
docker compose up --build
```

3. Open the app:

- Frontend: http://localhost
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

4. Register a new account, then add incomes, envelope allocations, and expenses from the UI.

### Family budgets

Each account belongs to a **household**. On signup:

- **Leave join code blank** to create a new family budget — you'll receive an 8-character code to share.
- **Enter a join code** from a family member to access the same budget.

The family join code appears in the app header (tap to copy). All household members see and edit the same incomes, envelopes, and expenses.

## Local development

### Backend only (Docker)

```bash
docker compose up db backend
```

### Frontend dev server

```bash
cd frontend
npm install
VITE_API_URL=http://localhost:8080 npm run dev
```

The Vite dev server runs at http://localhost:5173 and proxies API calls to the backend.

## API overview

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Create account |
| POST | `/api/auth/login` | Log in |
| GET | `/api/auth/me` | Current user |
| GET/POST | `/api/incomes` | List/create incomes |
| PUT/DELETE | `/api/incomes/{id}` | Update/delete income |
| GET/POST | `/api/envelopes` | List/create envelopes |
| PUT/DELETE | `/api/envelopes/{id}` | Update/delete envelope |
| GET/POST | `/api/expenses` | List/create expenses |
| PUT | `/api/expenses/{id}/void` | Void expense |
| GET | `/api/summary` | Dashboard totals |

All routes except register/login require `Authorization: Bearer <token>`.

## Envelope method

1. Add your **incomes** for the budget period.
2. Create **envelopes** (categories) and allocate amounts to each.
3. Log **expenses** against envelopes; void mistakes instead of deleting.
4. The **dashboard** shows unallocated funds, spending, and D3 charts.

Money is stored as integer cents in the database to avoid floating-point errors.

## Install as a PWA

The frontend is a Progressive Web App. After building (`docker compose up --build` or `npm run build`), you can install it like a native app.

### Desktop (Chrome / Edge)

1. Open http://localhost (or your deployed URL).
2. Click the install icon in the address bar, or use the browser menu → **Install Envelope Budget**.

### Android (Chrome)

1. Open the site in Chrome.
2. Tap the menu (⋮) → **Install app** or **Add to Home screen**.

### iPhone / iPad (Safari)

1. Open the site in **Safari**.
2. Tap Share → **Add to Home Screen**.

### Notes

- **HTTPS is required** for install prompts on real devices (localhost is exempt for development).
- If you install on a phone, set `VITE_API_URL` at build time to an API URL the phone can reach (not `http://localhost:8080` unless the backend is on the same device).
- The service worker caches the app shell for faster loads; API calls always go to the network.

## Production deployment

See **[DEPLOY.md](DEPLOY.md)** for the full VPS checklist.

**Already running nginx for other sites (e.g. casashoa.com)?** Use the nginx guide in DEPLOY.md — do **not** use `docker-compose.prod.yml` (Caddy). Use:

```bash
docker compose -p budgeting -f docker-compose.prod-nginx.yml up -d --build
```

**Empty VPS?** Use Caddy:

```bash
cp .env.production.example .env
docker compose -f docker-compose.prod.yml up -d --build
```

## Project structure

```
backend/     Go API and SQL migrations
frontend/    Vue 3 SPA with D3 charts
docker-compose.yml          Local development
docker-compose.prod.yml     Production (Caddy + SSL)
Caddyfile                   Reverse proxy config
DEPLOY.md                   VPS deployment guide
.env.production.example     Production env template
scripts/backup-db.sh        Database backup helper
```
