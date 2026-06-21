# Deployment

This project should be deployed as two services plus managed serverless data stores:

1. Angular frontend: Cloudflare Pages or Vercel static hosting
2. Go backend: Vercel Go Runtime, Railway, Render, Fly.io, or another Go-capable host
3. Primary database target: serverless Postgres, preferably Neon
4. Optional hot-path layer: Upstash Redis

> Current code still uses MongoDB. See `docs/database-decision.md` for the migration target and schema. Do not treat MongoDB Atlas as the desired long-term production database.

## Recommended production target

Use this path because it avoids always-on database infrastructure while keeping the data model relational:

- Frontend: Cloudflare Pages
- Backend: Render/Railway/Fly.io first; Vercel is possible with repo layout adjustment
- Primary DB: Neon Postgres
- Optional cache/rate-limit store: Upstash Redis

## Required environment variables

### Backend today, before DB migration

```bash
MONGOURI=mongodb://localhost:27017
SECRET_KEY=<long-random-secret>
ALLOWED_ORIGINS=https://<frontend-domain>
PORT=<set-by-host>
```

### Backend after Postgres migration

```bash
DATABASE_URL=postgres://<user>:<password>@<host>/<database>?sslmode=require
SECRET_KEY=<long-random-secret>
ALLOWED_ORIGINS=https://<frontend-domain>
PORT=<set-by-host>

# optional
UPSTASH_REDIS_REST_URL=https://...
UPSTASH_REDIS_REST_TOKEN=...
```

`PORT` is usually injected by the hosting provider. Set it manually only for local development.

### Frontend

Use an environment-specific API base URL that points to the deployed backend:

```bash
API_BASE_URL=https://<backend-domain>
```

The current frontend may need a small Angular environment/config cleanup before production deployment if API URLs are hardcoded.

## Backend local execution

```bash
cd backend
go mod download
go run main.go
```

Optional local `.env` today:

```bash
MONGOURI=mongodb://localhost:27017
SECRET_KEY=dev-only-secret
ALLOWED_ORIGINS=http://localhost:4200
PORT=8080
```

After the Postgres migration, replace `MONGOURI` with `DATABASE_URL`.

## Frontend local execution

```bash
cd frontend/forum-app
npm install
npm start
```

## Cloudflare Pages frontend settings

From the Cloudflare Pages Angular guide, use:

- Build command: `npm run build`
- Build output directory: Angular `dist/...` output directory
- Production branch: `main`

For this repo, set the project root to:

```text
frontend/forum-app
```

Confirm the final output directory after running `npm run build`; Angular 13 projects commonly output under `dist/forum-app`.

## Backend deployment notes

### Vercel

Vercel's Go runtime expects either:

- a root `go.mod` with `main.go`, `cmd/api/main.go`, or `cmd/server/main.go`; or
- Go files inside `/api` exporting `http.HandlerFunc`.

Because this repo keeps Go code under `backend/`, either deploy `backend/` as the Vercel project root or move/wrap the backend entrypoint for Vercel.

### Render/Railway/Fly.io

These are simpler for a long-running Go API because the existing `backend/main.go` already starts a server. Configure:

- root directory: `backend`
- build command: `go build -o app .`
- start command: `./app`
- environment variables today: `MONGOURI`, `SECRET_KEY`, `ALLOWED_ORIGINS`
- environment variables after migration: `DATABASE_URL`, `SECRET_KEY`, `ALLOWED_ORIGINS`

## Serverless Postgres setup

1. Create a Neon Postgres project.
2. Create a role/user with only the permissions needed by this app.
3. Copy the pooled connection string into `DATABASE_URL`.
4. Run migrations.
5. Smoke test signup, login, create post, list posts, create comment, vote.

## Upstash Redis setup, optional

Use Upstash Redis for data that can be rebuilt or expired:

- rate limits
- feed cache
- hot counters cache
- session denylist
- short-lived queues

Do not store canonical users/posts/comments only in Redis unless the project intentionally accepts weaker queryability and more application-level consistency work.

## Production checklist

- [ ] Merge backend pagination/security/votes PR.
- [ ] Complete Postgres migration from `docs/database-decision.md`.
- [ ] Create Neon Postgres project.
- [ ] Set backend environment variables.
- [ ] Deploy backend.
- [ ] Configure frontend API base URL.
- [ ] Deploy frontend.
- [ ] Set `ALLOWED_ORIGINS` to the exact frontend domain.
- [ ] Smoke test signup, login, create post, list posts, create comment, vote.
- [ ] Add uptime/error monitoring before sharing broadly.
