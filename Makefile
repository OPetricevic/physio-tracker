FRONTEND_DIR := frontend
BACKEND_DIR := backend
DB_URL ?= postgres://postgres:postgres@localhost:5432/physio?sslmode=disable

.PHONY: frontend-install frontend-dev frontend-build run backend-run backend-migrate backend-bootstrap

frontend-install:
	npm --prefix $(FRONTEND_DIR) install

frontend-dev:
	npm --prefix $(FRONTEND_DIR) run dev

frontend-build:
	npm --prefix $(FRONTEND_DIR) run build

# Convenience alias for local development
run: frontend-dev

# Backend: run the server (expects DB already created/migrated)
backend-run:
	cd $(BACKEND_DIR) && DATABASE_URL=$(DB_URL) PORT=3600 go run ./cmd/server

# Backend: apply SQL migrations in order using psql (requires DB to exist)
backend-migrate:
	cd $(BACKEND_DIR) && DATABASE_URL=$(DB_URL) ./scripts/migrate.sh

# Backend: create DB/user (if missing) using the default superuser connection
# WARNING: adjust DB_BOOTSTRAP_URL/DB_NAME/DB_APP_USER/DB_APP_PASS to your environment.
backend-bootstrap:
	cd $(BACKEND_DIR) && ./scripts/bootstrap_postgres.sh
