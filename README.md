# Physio Tracker (offline-first)

An offline-first app for physiotherapists/doctors to manage patients, anamneses, and generate PDFs. Runs locally (no internet required) with a Go backend, React frontend, and PostgreSQL. Installers/scripts provided for Windows and Linux.

## Features
- Patients CRUD, anamneses CRUD, PDF generation (Bosnian/Croatian diacritics supported).
- Include previous visits in PDFs; “only this visit” option.
- Doctor profile (logo, header, contact) stored locally.
- Backups via scripts (pg_dump/psql).

## Quick Start (developers)
Prereqs: Go, Node, PostgreSQL.
- Install deps: `make frontend-install`
- Run dev: `make dev` (backend + Vite dev proxy to `/api`)
- Build frontend: `make frontend-build`
- Run backend: `make backend-run` (set `DB_URL` if needed)
- Migrate: `make backend-migrate` (DB must exist)
- Package bundle: `make package` (creates `release/physio-bundle` with server, frontend, assets, migrations, start scripts)

## Installation for end users

### Windows (installer)
1) We ship an installer built from the bundle:
   - Build bundle: `make package`
   - Open `scripts/win/physio.iss` in Inno Setup and build the installer.
2) Run the installer:
   - Copies the app to `C:\Program Files\PhysioTracker`.
   - Installs/updates the “PhysioTracker” Windows service (runs backend, serves frontend).
   - Adds Start menu/Desktop shortcuts to `http://localhost:3600`.
3) Backups/restore:
   - `scripts/win/backup.ps1` (pg_dump)
   - `scripts/win/restore.ps1 -File path\to\backup.sql`

### Windows (portable/manual)
- Use the bundle: `release/physio-bundle`
- Ensure Postgres is running; set `DATABASE_URL` if different.
- Run `scripts/start_windows.ps1` (starts server on port 3600, serves frontend).

### Linux
1) Build bundle: `make package`
2) On target machine, run from the bundle folder:
   - `scripts/linux/install.sh`
   - Installs Postgres (if missing), creates DB/user (`physio`/`physio` by default), runs migrations.
   - Copies app to `/opt/physio`, installs a systemd service (`physio`), and a desktop launcher to open `http://localhost:3600`.
3) Backups/restore:
   - `scripts/linux/backup.sh`
   - `scripts/linux/restore.sh path/to/backup.sql`

## Defaults / Config
- HTTP port: `PORT` (default 3600)
- DB URL: `DATABASE_URL` (default `postgres://physio:physio@localhost:5433/physio?sslmode=disable` in Makefile; installer uses `:5432` on Linux)
- Frontend build served from `frontend/dist` (packaged into bundle)
- Uploads/logos stored under `uploads/` (served at `/static`)

## Repo layout (relevant)
- `backend/` Go server, migrations, fonts for PDFs
- `frontend/` React/Vite app
- `scripts/` start/installer/backup/restore (Windows & Linux)
- `release/physio-bundle/` (created by `make package`)

## Notes
- All data stays local. No internet calls.
- Fonts: DejaVu included for proper čćđšž.
- If you rebuild the app, rerun `make package` and rebuild the installer.
