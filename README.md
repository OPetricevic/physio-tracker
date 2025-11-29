# Doctor App Project Plan

Working notes to guide development of the offline-first doctor application (React frontend, Go backend, PostgreSQL). Update this file as the project evolves.

## Overview
- Goal: Offline desktop/local app for doctors to manage patients and anamneses, and generate PDFs of notes on demand.
- Offline-first: All assets local; no CDN. Later web deployment should require minimal tweaks (keep HTTP APIs clean).
- Users: Single doctor now; design for multi-doctor/login later (keep separation of concerns and add auth when needed).

## Stack Choices
- Backend: Go, REST style, `gorilla/mux` or `chi`, `pgx` driver. Migrations via `golang-migrate` (or similar). PDF generation via `gofpdf` (or `pdfcpu`).
- DB: PostgreSQL (preferred for future parity); SQLite acceptable for quick prototyping but Postgres recommended.
- Frontend: React + TypeScript (Vite). Data fetching via `@tanstack/react-query`, forms via `react-hook-form`. Styling: CSS modules or minimal local CSS (no external CDN).
- Packaging: Optional Docker Compose later for Postgres; keep local binaries/scripts working without containers.

## Project Layout
- `backend/`
  - `cmd/server/main.go`
  - `internal/`
    - `patients/` (handlers, service, repo)
    - `pdf/` (PDF generation utilities)
    - `db/` (queries, migrations helper)
  - `migrations/` (SQL up/down files)
  - `Makefile` (run, test, migrate, format)
- `frontend/`
  - `src/`
    - `components/` (forms, table, pdf buttons)
    - `pages/` (Patients page, Anamneses page)
    - `lib/api.ts` (API client)
    - `types/`
- `ops/`
  - `docker-compose.yml` (optional Postgres for later)
  - `backup/` (backup/restore scripts and dumps)

## Data Model (initial)
- `patients`: `id` (uuid), `first_name`, `last_name`, `phone`, `created_at`, `updated_at`.
- `anamneses`: `id`, `patient_id` (fk), `note` (text), `created_at`.
- Indexing: names/phone for search. Allow multiple anamneses per patient to keep history.

## API Sketch
- `POST /patients` (create patient)
- `GET /patients` (list/search)
- `GET /patients/{id}` (fetch detail)
- `POST /patients/{id}/anamneses` (add anamnesis)
- `GET /patients/{id}/anamneses` (list anamneses)
- `POST /patients/{id}/anamneses/{anamnesisId}/pdf` (generate/stream PDF of that note)

## PDF Generation
- Server-side Go using `gofpdf` (simple, offline). Return `application/pdf` stream.
- Optionally store generated PDFs on disk (e.g., `storage/pdfs/{patientId}/{anamnesisId}.pdf`) for quick re-download; regenerate if missing.

## Backup Strategy (offline-friendly)
- Preferred: Daily rolling dump plus manual "Backup now" trigger.
  - Postgres example: `pg_dump -Fc dbname > ops/backup/dumps/backup-YYYYMMDD-HHMM.dump`
  - Optionally zip dump + small JSON manifest: `backup-YYYYMMDD-HHMM.zip` for USB transfer.
- Endpoint and CLI script to trigger backup and download. Avoid heavy dump on every write; that is slow and brittle.
- Restore script (e.g., `ops/backup/restore.sh`) to load dump.

## Frontend Notes
- Pages: Patients list with search; form/drawer to add/edit; per-patient anamnesis list; per-entry PDF button.
- State/data: `react-query` for caching/fetch; keep API base URL configurable (env).
- UI: Local fonts/assets; no external calls. Keep layout simple and printable-friendly PDF view if needed.

## Backend Notes
- Config via env (DB URL, HTTP port, storage paths). Sensible defaults for local/dev.
- Logging, basic validation, and simple error responses (JSON).
- Migrations run on startup or via `make migrate`.

## Initial Setup Commands (to run manually)
- Backend scaffold:
  - `cd backend`
  - `go mod init doctorapp/backend`
  - `go get github.com/gorilla/mux github.com/jackc/pgx/v5 github.com/golang-migrate/migrate/v4 github.com/jung-kurt/gofpdf`
  - `mkdir -p cmd/server internal/{patients,pdf,db} migrations`
- Frontend scaffold:
  - `cd frontend`
  - `npm create vite@latest frontend -- --template react-ts`
  - `cd frontend`
  - `npm install @tanstack/react-query react-hook-form`

## Running locally (current)
- From repo root, use the Makefile helpers:
  - `make frontend-install` to install frontend deps.
  - `make frontend-dev` (or `make run`) to start the React dev server.
  - `make frontend-build` to build static assets.
- Later, add backend targets (run server, migrations, tests) and wire a launcher/shortcut that starts both backend and frontend locally.

## Frontend auth placeholder
- Basic login/registration screens exist for future auth. They store a mock user in `localStorage` and gate the workspace routes. No real backend yet.
- Flow: open app → `/login` or `/register` → upon submit you are routed to the patient workspace.

## Feature notes / changes
- Pacijenti/Anamneze view now paginates anamneses: 5 najnovijih po stranici, sortirano po datumu (DESC). Navigacija stranica je u panelu; backend treba vratiti sortirano po datumu kad se spoji.
- Top navigacija dodana: Pacijenti (trenutna funkcionalnost) i Raspored (placeholder za Google Calendar integraciju).

## Backend TODOs
- See `docs/BACKEND_TODO.md` for planned backend scaffolding, run scripts, deployment helper, PDF/assets plan, and backups.
- See `docs/ARCHITECTURE.md` for ports/adapters naming and folder conventions (inbound/outbound/service).

## Next Steps (candidate tasks)
1) Scaffold backend server, wiring routes and simple health check. Add sample migration for patients/anamneses.  
2) Scaffold frontend (Vite React TS) with Patients page and forms hooked to mock API client.  
3) Add PDF generation endpoint + storage.  
4) Add backup scripts (`backup-now.sh`, `restore.sh`) and optional API trigger.  
5) Add search and basic validations; prepare for auth later.
