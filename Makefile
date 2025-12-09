FRONTEND_DIR := frontend
BACKEND_DIR := backend
DB_URL ?= postgres://postgres:postgres@localhost:5432/physio?sslmode=disable
PROTO_DIR := $(BACKEND_DIR)/protos
PROTO_OUT := $(BACKEND_DIR)/golang
PROTOC := protoc
PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GORM := $(shell go env GOPATH)/bin/protoc-gen-gorm

.PHONY: frontend-install frontend-dev frontend-build run backend-run backend-migrate backend-bootstrap backend-proto

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

# Backend: generate Go code from protos
backend-proto:
	@command -v $(PROTOC) >/dev/null || { echo "protoc not found; install protoc"; exit 1; }
	@test -x $(PROTOC_GEN_GO) || { echo "protoc-gen-go missing; run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
	@test -x $(PROTOC_GEN_GORM) || { echo "protoc-gen-gorm missing; run: go install github.com/infobloxopen/protoc-gen-gorm@latest"; exit 1; }
	$(PROTOC) -I $(PROTO_DIR) --go_out=$(PROTO_OUT) --go_opt=paths=source_relative --gorm_out=$(PROTO_OUT) --gorm_opt=paths=source_relative $(PROTO_DIR)/*.proto
