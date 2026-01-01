#!/usr/bin/env bash
set -e

# This script starts the local server and serves the built frontend.
# Defaults can be overridden via environment variables.

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$DIR" && pwd)"

export FRONTEND_DIR="${FRONTEND_DIR:-$ROOT/frontend/dist}"
export PORT="${PORT:-3600}"
export DATABASE_URL="${DATABASE_URL:-postgres://physio:physio@localhost:5433/physio?sslmode=disable}"

cd "$ROOT"

if [ ! -x "$ROOT/server" ]; then
  echo "server binary not found in $ROOT. Build/package first."
  exit 1
fi

echo "Starting server on :$PORT (DATABASE_URL=$DATABASE_URL)"
"$ROOT/server"
