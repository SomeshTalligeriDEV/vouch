# Contributing to Vouch

## Prerequisites

- Go 1.22+
- Node.js 20+
- Docker + Docker Compose (for local deps)

## Local setup

```bash
# Start MongoDB + Redis
cd vouch-api && docker-compose up -d mongo redis

# Backend
cp vouch-api/.env.example vouch-api/.env   # fill in values
cd vouch-api && go run ./cmd/api

# Frontend
cp vouch-web/.env.example vouch-web/.env.local
cd vouch-web && npm install && npm run dev
```

## Development workflow

1. **Branch** off `main`: `git checkout -b feat/your-feature`
2. **Write tests first** where practical (service layer especially).
3. **Run checks** before pushing:
   ```bash
   # Backend
   cd vouch-api
   go build ./...
   go test ./... -race
   golangci-lint run

   # Frontend
   cd vouch-web
   npx tsc --noEmit
   npm run lint
   ```
4. **Open a PR** — fill in the PR template.

## Code conventions

- Follow the layered architecture: `domain → repository → service → handler`.
- Never import `handler` from `service`, or `service` from `repository`.
- All new service methods need at least one test in `*_test.go`.
- No `fmt.Println` in production code — use `zerolog`.
- Frontend: React Query for server state, Zustand only for auth.

## Commit style

```
feat(scope): short present-tense description
fix(scope): what was wrong
chore(scope): maintenance
docs(scope): documentation
test(scope): tests only
```

## Questions?

Open an issue or discussion on GitHub.
