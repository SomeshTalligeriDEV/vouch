# vouch-api

Go + Fiber backend for **Vouch** — a demand-first builder reputation platform that
replaces GitHub stars and LeetCode ratings with verified proof of real shipped
products, real users, and real revenue.

## Architecture

Strict layered architecture. Dependencies point inward only:

```
handler  →  service  →  repository  →  MongoDB
   │            │
   │            └── domain (pure entities, interfaces, score engine — zero deps)
   └── middleware (auth, ratelimit, logger)

worker (Asynq)  →  service  →  repository
external (GitHub / Stripe gateways)  →  implement service ports
```

- **domain** — entities, errors, repository interfaces, and the pure score engine. No MongoDB, HTTP, or framework imports.
- **repository/mongo** — MongoDB implementations only. No business logic.
- **service** — business logic. Calls repositories and domain methods; never touches HTTP.
- **handler** — parses requests, calls services, returns standardized responses.
- **worker** — Asynq consumers for async score recalculation and Stripe sync.

Every function takes `context.Context` first. Every error is wrapped
(`fmt.Errorf("layer.Method: %w", err)`). All responses go through `pkg/response`.

## Score engine

```
Total = (UserScore + RevenueScore + ImpactScore + VelocityScore) × StripeMultiplier

UserScore     = min(verified_users × 10, 30000)
RevenueScore  = min(mrr × 2, 20000)
ImpactScore   = min(avg_rating × review_count × 5, 15000)
VelocityScore = min(90_day_growth × 0.1, 5000)
StripeMultiplier = 1.0 if Stripe verified, else 0.6
```

Tiers: Bronze (0–999) · Silver (1k–4.9k) · Gold (5k–14.9k) · Platinum (15k–49.9k) · 24 Karat (50k+)

## Run locally

```bash
cp .env.example .env          # fill in secrets
docker compose -f docker/docker-compose.yml up --build
# or, against local Mongo/Redis:
go run ./cmd/api
go run ./cmd/worker
```

## Test

```bash
go vet ./...
go test ./...
```

## Seed demo data

```bash
go run ./cmd/seed     # builders, projects, reviews, problems, and computed scores
```

## Storage, email & API collection

- **Uploads** — `POST /uploads/presign` returns a presigned Cloudflare R2 (S3-compatible)
  PUT URL plus the resulting public URL; the browser uploads directly. See
  `internal/external/r2.go`.
- **Email** — problem claims enqueue an async `email:problem_claimed` task; the worker
  sends via Resend (`internal/external/resend.go`). A missing `RESEND_API_KEY` is a
  no-op, so local dev works without email configured.
- **`api.http`** — a ready-to-run REST Client collection covering every endpoint.

## API

Base path `/api/v1`. See `internal/handler/router.go` for the full route table.
Mutation endpoints require a Bearer JWT and are rate-limited via Redis.

| Method | Path | Auth |
|---|---|---|
| POST | `/auth/github` | – |
| POST | `/auth/refresh` | – |
| GET | `/users/:username` | – |
| PATCH | `/users/me` | ✓ |
| POST | `/users/me/stripe` | ✓ |
| GET | `/projects` | – |
| POST | `/projects` | ✓ |
| GET/PATCH/DELETE | `/projects/:id` | mutations ✓ |
| GET | `/scores` (leaderboard) · `/scores/:username` | – |
| POST | `/scores/recalculate` | ✓ |
| GET | `/problems` · `/problems/:id` | – |
| POST | `/problems` · `/problems/:id/claim` · `/problems/:id/upvote` | ✓ |
| POST | `/reviews` | ✓ |
| GET | `/reviews/project/:id` | – |

## Deploy

Pushes to `main` touching `vouch-api/**` run vet + tests + build, then
`flyctl deploy` (see `.github/workflows/deploy.yml`). Requires `FLY_API_TOKEN` secret.
