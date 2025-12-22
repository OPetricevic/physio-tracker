#!/usr/bin/env bash
set -euo pipefail

# Simple helper to run backend API and frontend Vite dev server together.
# Assumes deps are installed: `go` modules fetched, `npm install` in frontend,
# and Postgres reachable via DATABASE_URL.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"

# Env defaults (override via env when calling the script)
export DATABASE_URL="${DATABASE_URL:-postgres://physio:physio@localhost:5433/physio?sslmode=disable}"
export PORT="${PORT:-3600}"
export VITE_API_URL="${VITE_API_URL:-/api}"

echo "Starting backend on :${PORT} (DATABASE_URL=${DATABASE_URL})"
(
  cd "$BACKEND_DIR"
  go run ./cmd/server
) &
BACK_PID=$!

echo "Starting frontend (VITE_API_URL=${VITE_API_URL})"
(
  cd "$FRONTEND_DIR"
  npm run dev -- --host 0.0.0.0
) &
FRONT_PID=$!

cleanup() {
  echo "Shutting down..."
  kill "${BACK_PID}" "${FRONT_PID}" 2>/dev/null || true
}
trap cleanup INT TERM EXIT

wait
