# vouch-web

Next.js (App Router) frontend for **Vouch**. Typed against the `vouch-api` Go backend.

## Stack
- Next.js 15 + React 18, TypeScript (strict)
- Tailwind CSS
- TanStack React Query (server state) + Zustand (auth only)
- GitHub OAuth login (no passwords)

## Structure
```
app/            App Router routes (landing, auth, dashboard, discover, problems, builder profile)
components/     UI + feature components (project, problem, score)
lib/api.ts      Typed API client — ALL fetch calls live here
lib/auth.ts     Token persistence
hooks/          React Query hooks
store/          Zustand auth store
types/          TypeScript types mirroring the Go domain structs
```

## Develop
```bash
cp .env.local.example .env.local   # set NEXT_PUBLIC_API_URL + GitHub client id
npm install
npm run dev                         # http://localhost:3000
```

## Verify
```bash
npm run typecheck
npm run build
```

## Env
| Var | Purpose |
|---|---|
| `NEXT_PUBLIC_API_URL` | Base URL of vouch-api (e.g. `http://localhost:8080/api/v1`) |
| `NEXT_PUBLIC_GITHUB_CLIENT_ID` | GitHub OAuth app client id |
| `NEXT_PUBLIC_GITHUB_REDIRECT_URL` | OAuth redirect (the `/login` page handles the code exchange) |

The API client transparently refreshes the access token once on a 401 using the
stored refresh token, then replays the request.
