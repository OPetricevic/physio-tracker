#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 path/to/backup.sql"
  exit 1
fi

FILE="$1"
if [ ! -f "$FILE" ]; then
  echo "File not found: $FILE"
  exit 1
fi

DB_URL="${DATABASE_URL:-postgres://physio:physio@localhost:5432/physio?sslmode=disable}"
echo "Restoring $FILE into $DB_URL"
psql "$DB_URL" < "$FILE"
echo "Restore complete"
