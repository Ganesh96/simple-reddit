# Deployment

This project should be deployed as two services plus one managed database:

1. Angular frontend: Cloudflare Pages or Vercel static hosting
2. Go backend: Vercel Go Runtime, Railway, Render, Fly.io, or another Go-capable host
3. MongoDB: MongoDB Atlas

## Recommended MVP deployment

Use this path first because it has the least infrastructure overhead:

- Frontend: Cloudflare Pages
- Backend: Vercel Go Runtime or Render/Railway
- Database: MongoDB Atlas

## Required environment variables

### Backend

```bash
MONGOURI=mongodb+srv://<user>:<password>@<cluster>/<database>?retryWrites=true&w=majority
SECRET_KEY=<long-random-secret>
ALLOWED_ORIGINS=https://<frontend-domain>
PORT=<set-by-host>
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

Optional local `.env`:

```bash
MONGOURI=mongodb://localhost:27017
SECRET_KEY=dev-only-secret
ALLOWED_ORIGINS=http://localhost:4200
PORT=8080
```

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
- environment variables: `MONGOURI`, `SECRET_KEY`, `ALLOWED_ORIGINS`

## MongoDB Atlas setup

1. Create an Atlas cluster.
2. Create a database user with only the permissions needed by this app.
3. Restrict network access to the backend host where possible.
4. Copy the Go connection string into `MONGOURI`.
5. Start the backend once so `EnsureIndexes()` creates required indexes.

## Production checklist

- [ ] Merge backend pagination/security/votes PR.
- [ ] Create MongoDB Atlas cluster.
- [ ] Set backend environment variables.
- [ ] Deploy backend.
- [ ] Configure frontend API base URL.
- [ ] Deploy frontend.
- [ ] Set `ALLOWED_ORIGINS` to the exact frontend domain.
- [ ] Smoke test signup, login, create post, list posts, create comment, vote.
- [ ] Add uptime/error monitoring before sharing broadly.
