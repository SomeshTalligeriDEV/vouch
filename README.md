# Vouch

> We replace GitHub stars and LeetCode ratings with verified proof —
> real products shipped, real users paying, real revenue earned.

A demand-first builder reputation platform.

## Three core flows
1. **Supply** — builders post existing/dead GitHub projects → real users discover and buy them.
2. **Demand** — users post real problems with a budget → builders claim and ship → the poster becomes the first paying user.
3. **Reputation** — every shipped product builds a verified Builder Score. Companies hire by score, not DSA.

## Monorepo layout
```
vouch-api/   Go + Fiber backend  — strict layered architecture, MongoDB, Redis + Asynq, JWT + GitHub OAuth, Stripe (read-only)
vouch-web/   Next.js frontend    — App Router, React Query, Zustand, Tailwind, typed API client
```

Each repo is self-contained with its own README, Dockerfile/build, and CI.

## Score engine
```
Total = (UserScore + RevenueScore + ImpactScore + VelocityScore) × StripeMultiplier
```
Tiers: Bronze → Silver → Gold → Platinum → 24 Karat.
Source of truth: `vouch-api/internal/domain/score.go` (pure, unit-tested).

## Quick start
```bash
# Backend
cd vouch-api && cp .env.example .env && docker compose -f docker/docker-compose.yml up --build

# Frontend
cd vouch-web && cp .env.local.example .env.local && npm install && npm run dev
```
