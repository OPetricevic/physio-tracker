# Backend local setup and bootstrap

Use these steps to bring up the backend on a fresh machine (no DB yet) and to understand what each script does.

## Prerequisites
- Go toolchain installed.
- PostgreSQL server running locally and a superuser connection available (default assumptions below).
- `psql` client installed.

## Scripts overview
- `backend/scripts/bootstrap_postgres.sh`  
  Creates an application role and database if they don’t exist. Defaults:
  - superuser DSN: `DB_BOOTSTRAP_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`
  - DB name: `DB_NAME=physio`
  - app user/pass: `DB_APP_USER=physio_app`, `DB_APP_PASS=physio_app_pass`
  After running, set `DATABASE_URL=postgres://physio_app:physio_app_pass@localhost:5432/physio?sslmode=disable`.

- `backend/scripts/migrate.sh`  
  Applies SQL migrations (`0001_init.sql`, `0002_seed_admin.sql`) to the database pointed to by `DATABASE_URL`.

## Makefile shortcuts
- `make backend-bootstrap` — runs the bootstrap script with defaults (override env vars as needed).
- `make backend-migrate` — runs migrations using `DB_URL` (defaults to `postgres://postgres:postgres@localhost:5432/physio?sslmode=disable`).
- `make backend-run` — starts the server with `DB_URL` and `PORT=3600` (override via env).

## Typical local flow
1) Bootstrap DB/user (once):
   ```bash
   make backend-bootstrap
   # or: cd backend && DB_BOOTSTRAP_URL=... DB_NAME=... DB_APP_USER=... DB_APP_PASS=... ./scripts/bootstrap_postgres.sh
   ```
2) Migrate:
   ```bash
   DB_URL="postgres://physio_app:physio_app_pass@localhost:5432/physio?sslmode=disable" make backend-migrate
   # or: cd backend && DATABASE_URL=... ./scripts/migrate.sh
   ```
3) Run server:
   ```bash
   DB_URL="postgres://physio_app:physio_app_pass@localhost:5432/physio?sslmode=disable" make backend-run
   # or: cd backend && DATABASE_URL=... PORT=3600 go run ./cmd/server
   ```

## Quick smoke test (after server is running)
```bash
curl -X POST http://localhost:3600/patients/create \
  -H "Content-Type: application/json" \
  -d '{"doctorUuid":"<doctor_uuid>","firstName":"Mia","lastName":"Horvat"}'
```
Use the seeded doctor from migrations for `doctorUuid`.
