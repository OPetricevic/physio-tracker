# Backend TODOs and Ops Plan

Working list for backend and deployment tasks. Adjust as requirements arrive (e.g., final domain, branding assets).

## Core backend
- Scaffold Go service (`cmd/server/main.go`): load config (.env), open Postgres, run migrations, start HTTP server.
- DB: Postgres with UUID primary keys. Tables: `patients`, `anamneses`; later `users` for auth.
- Migrations: add tool (`golang-migrate` CLI or embedded). Provide `make migrate-up/down`.
- Routes (initial): health, patients CRUD, anamneses CRUD, PDF generation (later).

## Config & secrets
- Use `.env.local` (not committed) for DB URL, JWT secret (later), and Git token if needed.
- Fine-grained GitHub token (read/pull) stored locally; never in repo. Configure `git` to use credential helper or env var for automation scripts.

## Run scripts (future)
- Single prod entry (e.g., `make prod` or `./run.ps1`) that:
  1) Builds frontend (npm run build) and serves static assets (backend or lightweight file server).
  2) Ensures DB exists; runs migrations idempotently; creates DB if missing.
  3) Starts backend on a fixed port (e.g., 3600) and serves the app.
- Windows helper (PowerShell) plan:
  - `git pull origin main` to update.
  - Ensure hosts entry for chosen domain (placeholder: `app.localclinic.local` -> `127.0.0.1`). Make it idempotent.
  - If an existing instance is running, stop it cleanly (by port or named process), then start new.
  - Run prod script, hide window, open `http://app.localclinic.local` in browser.

## Domain placeholder
- Decide later: e.g., `app.localclinic.local` or similar; avoid hardcoding until confirmed.

## PDF (later when assets arrive)
- Store branding assets locally (`backend/assets/branding/logo.png`, `doctor.jpg`).
- Template helper for headers/footers; stream PDF and store copy under `storage/pdfs/{patientUuid}/{anamnesisUuid}.pdf` (gitignored).

## Auth (later)
- Add users table, bcrypt hashes, JWT (or session) middleware; keep toggle `AUTH_ENABLED`.

## Backups
- Add `pg_dump` scripts (manual/cron) and API trigger; zip dump + manifest for USB transfer.
