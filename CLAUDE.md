# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Services

Four services, all run locally (only the databases are containerised):

| Service | Dir | Port | Stack |
|---|---|---|---|
| Backend | `SMS.Backend/` | 8080 | Go 1.23 |
| Warehouse | `SMS.Warehouse/` | 8081 | Go 1.23 |
| BackOffice | `SMS.BackOffice/` | 5173 | SvelteKit + Svelte 5 |
| Page Viewer (example) | `examples/page-viewer/` | 5174 | Vite + Svelte 5 |

## Running things

**Databases** (required first):
```bash
docker compose up -d          # starts postgres:5432, postgres-warehouse:5433, pgadmin:5050
```

**Backend:**
```bash
cd SMS.Backend
go run .
# env: DATABASE_URL, WAREHOUSE_URL (default http://localhost:8081), ALLOWED_ORIGINS, PORT
```

**Warehouse:**
```bash
cd SMS.Warehouse
go run .
# env: DATABASE_URL (default postgres://postgres:postgres@localhost:5433/warehouse), ALLOWED_ORIGINS, PORT
```

**BackOffice:**
```bash
cd SMS.BackOffice
npm run dev       # dev server
npm run check     # svelte-check type checking (use this, not tsc)
npm run build     # production build
```

**Page Viewer:**
```bash
cd examples/page-viewer
npm run dev
npm run build
```

## Go build & check

```bash
cd SMS.Backend && go build ./...
cd SMS.Warehouse && go build ./...
```

No tests exist yet. There are no linters configured beyond the compiler.

## Architecture

### Data model

The two services own separate Postgres databases:

**Backend DB** (`kms`, port 5432): `pages` → `boxes`
- A `box` belongs to a `page`, has a `type`, JSONB `content`, integer `position`, and an optional `datasource_id` (plain `TEXT`, not a foreign key — references a datasource in the warehouse DB).

**Warehouse DB** (`warehouse`, port 5433): `datasources` → `datapoints`
- A `datasource` has `type`, `description`, and `content` fields.
- A `datapoint` belongs to a datasource and carries a `tag` (one of `h1`, `h2`, `p`, `li`), `content` text, and `description`.

### Service responsibilities

**SMS.Backend** owns pages and boxes. It exposes one warehouse proxy route (`GET /datasources`) so the backoffice can populate the datasource picker without needing to know the warehouse URL. All other warehouse operations bypass the backend.

**SMS.Warehouse** owns datasources and datapoints. It is the source of truth for all content data.

**SMS.BackOffice** calls the backend for pages/boxes CRUD, and calls the warehouse **directly** (via `VITE_WAREHOUSE_URL`, default `http://localhost:8081`) for datasource and datapoint CRUD. The `api.ts` file has two base-URL helpers: `req()` → backend, `wreq()` → warehouse.

**examples/page-viewer** is a minimal consumer demonstrating real usage. It fetches a page by slug from the backend (`GET /slug/{slug}`), then fetches each referenced datasource directly from the warehouse, and renders the datapoints through tag-mapped Svelte components (`H1`, `H2`, `P`, `Li` in `src/components/`).

### Migrations

Both Go services embed their SQL migrations and apply them with `IF NOT EXISTS` guards at startup via `database.Migrate()`. Adding a migration means creating a new numbered `.sql` file in `internal/database/migrations/`, embedding it, and appending it to the slice in `Migrate()`.

### CORS

Both Go services use an `ALLOWED_ORIGINS` env var (comma-separated list). Defaults to `http://localhost:5173,http://localhost:5174`. The `cors()` middleware in each `main.go` reflects the request `Origin` header only when it matches the allowlist.

### Public slug route

`GET /slug/{slug}` on the backend returns a full page with boxes. This is intentionally separate from `GET /pages/{pageID}` to avoid Go's mux conflict with `GET /pages/{pageID}/boxes`. The slug route is the public-facing entrypoint; the ID route is for backoffice use.
