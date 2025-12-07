#!/usr/bin/env bash
set -euo pipefail

# Creates app database/user if missing. Requires superuser access (default: postgres@localhost).
# Override as needed: DB_BOOTSTRAP_URL, DB_NAME, DB_APP_USER, DB_APP_PASS.

DB_BOOTSTRAP_URL="${DB_BOOTSTRAP_URL:-postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable}"
DB_NAME="${DB_NAME:-physio}"
DB_APP_USER="${DB_APP_USER:-physio_app}"
DB_APP_PASS="${DB_APP_PASS:-physio_app_pass}"

psql "$DB_BOOTSTRAP_URL" <<SQL
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '$DB_APP_USER') THEN
      CREATE ROLE $DB_APP_USER LOGIN PASSWORD '$DB_APP_PASS';
   END IF;
END
\$\$;

DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME') THEN
      CREATE DATABASE $DB_NAME OWNER $DB_APP_USER;
   END IF;
END
\$\$;
SQL

echo "Bootstrap complete. App user: $DB_APP_USER DB: $DB_NAME"
echo "Set DATABASE_URL=postgres://$DB_APP_USER:$DB_APP_PASS@localhost:5432/$DB_NAME?sslmode=disable"
