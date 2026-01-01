#!/usr/bin/env bash
set -euo pipefail

# Linux installer for Physio Tracker
# - Installs Postgres (if missing), creates DB/user
# - Copies bundled files to /opt/physio
# - Installs systemd service and desktop launcher
# Run from the release bundle folder (release/physio-bundle)

APP_DIR="/opt/physio"
SERVICE_NAME="physio"
DB_NAME="${DB_NAME:-physio}"
DB_USER="${DB_USER:-physio}"
DB_PASS="${DB_PASS:-physio}"
PORT="${PORT:-3600}"
DATABASE_URL="${DATABASE_URL:-postgres://$DB_USER:$DB_PASS@localhost:5432/$DB_NAME?sslmode=disable}"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
BUNDLE="$(cd "$ROOT" && pwd)"

echo "Using bundle at $BUNDLE"

need_sudo() {
  if [[ $EUID -ne 0 ]]; then
    sudo "$@"
  else
    "$@"
  fi
}

echo "Installing dependencies (PostgreSQL, curl)..."
need_sudo apt-get update -y
need_sudo apt-get install -y postgresql postgresql-contrib curl

echo "Ensuring Postgres is running..."
need_sudo systemctl enable postgresql
need_sudo systemctl start postgresql

echo "Creating DB user/db if missing..."
need_sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1 || \
  need_sudo -u postgres psql -c "CREATE ROLE $DB_USER WITH LOGIN PASSWORD '$DB_PASS';"
need_sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" | grep -q 1 || \
  need_sudo -u postgres psql -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"

echo "Copying app to $APP_DIR..."
need_sudo rm -rf "$APP_DIR"
need_sudo mkdir -p "$APP_DIR"
need_sudo cp -r "$BUNDLE"/* "$APP_DIR"/

echo "Running migrations..."
for f in "$APP_DIR"/migrations/*.sql; do
  need_sudo -u postgres psql "$DB_NAME" -f "$f"
done

echo "Creating systemd service..."
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
need_sudo bash -c "cat > $SERVICE_FILE" <<EOF
[Unit]
Description=Physio Tracker
After=network.target postgresql.service

[Service]
Type=simple
WorkingDirectory=$APP_DIR
Environment=DATABASE_URL=$DATABASE_URL
Environment=PORT=$PORT
Environment=FRONTEND_DIR=$APP_DIR/frontend/dist
ExecStart=$APP_DIR/server
Restart=always

[Install]
WantedBy=multi-user.target
EOF

need_sudo systemctl daemon-reload
need_sudo systemctl enable "$SERVICE_NAME"
need_sudo systemctl restart "$SERVICE_NAME"

echo "Creating desktop launcher..."
DESKTOP_FILE="$HOME/.local/share/applications/physio.desktop"
mkdir -p "$(dirname "$DESKTOP_FILE")"
cat > "$DESKTOP_FILE" <<EOF
[Desktop Entry]
Type=Application
Name=Physio Tracker
Exec=xdg-open http://localhost:$PORT
Icon=applications-office
Terminal=false
Categories=Office;
EOF

echo "Done. Service '$SERVICE_NAME' running on http://localhost:$PORT"
echo "Backups: use $APP_DIR/scripts/backup.sh"
