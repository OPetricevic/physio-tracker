FRONTEND_DIR := frontend
BACKEND_DIR := backend
# default app creds (created by bootstrap script)
DB_URL ?= postgres://physio:physio@localhost:5433/physio?sslmode=disable
PROTO_DIR := $(BACKEND_DIR)/protos
PROTO_OUT := $(BACKEND_DIR)/golang
PROTOC := protoc
PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GORM := $(shell go env GOPATH)/bin/protoc-gen-gorm
PROTOC_GEN_VALIDATE := $(shell go env GOPATH)/bin/protoc-gen-validate
# Fixed include path for PGV (pinned to v0.6.13)
VALIDATE_INC := $(shell go env GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.13

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
	@test -x $(PROTOC_GEN_VALIDATE) || { echo "protoc-gen-validate missing; run: go install github.com/envoyproxy/protoc-gen-validate/cmd/protoc-gen-validate@latest"; exit 1; }
	@test -d "$(VALIDATE_INC)" || { echo "validate proto include not found; ensure protoc-gen-validate v0.6.13 is downloaded"; exit 1; }
	$(PROTOC) -I $(PROTO_DIR) -I $(VALIDATE_INC) --go_out=$(PROTO_OUT) --gorm_out=$(PROTO_OUT) --validate_out="lang=go,paths=source_relative:$(PROTO_OUT)" $(PROTO_DIR)/*.proto

# Run backend and frontend together (expects DB ready and npm install done)
dev:
	./scripts/dev.sh

# Alias: generate protos (Go + gorm)
proto: backend-proto
