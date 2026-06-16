# VPS deployment guide

Deploy the app as Docker containers with automatic HTTPS via Caddy and Let's Encrypt.

## Architecture

```
Internet :443/:80
    └── Caddy (TLS termination)
            ├── APP_DOMAIN  → frontend (nginx, static Vue PWA)
            └── API_DOMAIN  → backend (Go API)
                                    └── postgres (internal only)
```

## Repo files used in production

| File | Purpose |
|------|---------|
| `docker-compose.prod.yml` | Production services (no exposed db/backend ports) |
| `Caddyfile` | Reverse proxy + automatic SSL |
| `.env.production.example` | Template for production secrets and domains |
| `scripts/backup-db.sh` | Postgres backup helper |

---

## Before you touch the VPS

Complete these on your machine or in your DNS provider:

- [ ] Choose two subdomains, e.g. `budget.example.com` and `api.budget.example.com`
- [ ] Push the latest code to GitHub (or your git remote)
- [ ] Generate secrets:
  ```bash
  openssl rand -hex 32   # JWT_SECRET
  openssl rand -hex 24   # POSTGRES_PASSWORD
  ```

---

## VPS setup checklist

### 1. Provision the server

- [ ] Ubuntu 22.04 or 24.04 LTS recommended (1 GB RAM minimum, 2 GB preferred)
- [ ] Note the public IPv4 address

### 2. DNS

Create **A records** pointing at the VPS IP:

| Host | Type | Value |
|------|------|-------|
| `budget` | A | `YOUR_VPS_IP` |
| `api.budget` | A | `YOUR_VPS_IP` |

- [ ] Wait for DNS propagation (`dig budget.example.com`, `dig api.budget.example.com`)

### 3. SSH hardening (recommended)

- [ ] Create a non-root sudo user
- [ ] Disable password SSH / use keys only
- [ ] Enable UFW:
  ```bash
  sudo ufw allow OpenSSH
  sudo ufw allow 80/tcp
  sudo ufw allow 443/tcp
  sudo ufw enable
  ```

### 4. Install Docker

```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

- [ ] Log out and back in so the `docker` group applies
- [ ] Verify: `docker run hello-world`

### 5. Clone the app

```bash
git clone https://github.com/YOUR_USER/Budgeting-App.git
cd Budgeting-App
```

### 6. Create production `.env`

```bash
cp .env.production.example .env
nano .env
```

Set every `CHANGE_ME` value and replace example domains with yours:

| Variable | Example |
|----------|---------|
| `POSTGRES_PASSWORD` | random hex |
| `DATABASE_URL` | must use the same password |
| `JWT_SECRET` | random hex |
| `APP_DOMAIN` | `budget.example.com` |
| `API_DOMAIN` | `api.budget.example.com` |
| `VITE_API_URL` | `https://api.budget.example.com` |
| `CORS_ORIGINS` | `https://budget.example.com` |
| `ACME_EMAIL` | your real email |

- [ ] `VITE_API_URL` uses **https** and matches `API_DOMAIN`
- [ ] `CORS_ORIGINS` uses **https** and matches `APP_DOMAIN` (no trailing slash)
- [ ] `.env` is **not** committed to git

### 7. Build and start

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

- [ ] All containers running: `docker compose -f docker-compose.prod.yml ps`
- [ ] Caddy obtained certificates (check logs): `docker compose -f docker-compose.prod.yml logs caddy`

### 8. Smoke test

```bash
curl https://api.budget.example.com/health
# {"status":"ok"}

curl -I https://budget.example.com
# HTTP/2 200
```

In a browser:

- [ ] Open `https://budget.example.com`
- [ ] Register a new account
- [ ] Add an envelope, log an expense, check dashboard
- [ ] Install as PWA (optional)

---

## Updating after code changes

```bash
cd Budgeting-App
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

If `VITE_API_URL` or domains changed, the **frontend must be rebuilt** (the command above does that).

---

## Backups

Manual backup:

```bash
chmod +x scripts/backup-db.sh
./scripts/backup-db.sh
```

Backups are written to `./backups/`.

Optional cron (daily at 3 AM):

```bash
crontab -e
# add:
0 3 * * * /home/YOU/Budgeting-App/scripts/backup-db.sh >> /home/YOU/backup.log 2>&1
```

- [ ] Copy backups off the VPS (S3, another server, etc.)

---

## Troubleshooting

| Symptom | Likely cause |
|---------|----------------|
| Caddy fails to get cert | DNS not propagated, port 80 blocked, wrong domain in `.env` |
| API calls fail in browser | `CORS_ORIGINS` mismatch, or `VITE_API_URL` wrong at build time |
| PWA won't install | Site not served over HTTPS |
| 502 from Caddy | Backend not healthy: `docker compose -f docker-compose.prod.yml logs backend` |
| Database connection error | `DATABASE_URL` password doesn't match `POSTGRES_PASSWORD` |

View logs:

```bash
docker compose -f docker-compose.prod.yml logs -f caddy
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
```

---

## Security notes

- Postgres and the Go API are **not** published to the host in production compose
- Only Caddy exposes ports 80 and 443
- Rotate `JWT_SECRET` only if you accept invalidating all active sessions
- Keep the VPS updated: `sudo apt update && sudo apt upgrade`
