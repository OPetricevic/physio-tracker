#!/usr/bin/env bash
set -euo pipefail

DB_URL="${DATABASE_URL:-postgres://physio:physio@localhost:5432/physio?sslmode=disable}"
BACKUP_DIR="${BACKUP_DIR:-/opt/physio/backups}"
mkdir -p "$BACKUP_DIR"
ts=$(date +"%Y%m%d_%H%M%S")
file="$BACKUP_DIR/backup_$ts.sql"
echo "Backing up to $file"
pg_dump "$DB_URL" > "$file"
echo "Backup complete"
