# Simple Reddit Web

The frontend has been migrated from Angular to React with Vite.

## Stack

- React
- TypeScript
- Vite
- React Router
- Native `fetch` API client

## Local development

```bash
cd frontend/forum-app
npm install
npm start
```

The development server runs on `http://localhost:4200`.

## API configuration

By default, the frontend calls the backend at `http://localhost:8080`.

For deployed environments, set:

```bash
VITE_API_BASE_URL=https://<backend-domain>
```

## Build

```bash
npm run build
```

The production build outputs to `dist/`.

## Notes

The React migration intentionally prioritizes the highest-value read paths first:

- feed
- communities
- post detail
- comments

Auth/write flows are separated behind the backend's existing secured endpoints and should be completed in the next focused frontend PR.
